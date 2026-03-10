package main

import (
	"context"
	"fmt"
	"io"
	"sort"
	"sync"
)

// plan computes actions by comparing YAML entries against existing AWS state.
func plan(entries []Entry, existing map[string]string, prune bool) []Action {
	var actions []Action
	yamlKeys := make(map[string]bool)

	for _, e := range entries {
		yamlKeys[e.Key] = true
		awsVal, exists := existing[e.Key]
		switch {
		case !exists:
			actions = append(actions, Action{Key: e.Key, Type: ActionCreate, Value: e.Value})
		case awsVal != e.Value:
			actions = append(actions, Action{Key: e.Key, Type: ActionUpdate, Value: e.Value})
		default:
			actions = append(actions, Action{Key: e.Key, Type: ActionSkip})
		}
	}

	if prune {
		var deleteKeys []string
		for k := range existing {
			if !yamlKeys[k] {
				deleteKeys = append(deleteKeys, k)
			}
		}
		sort.Strings(deleteKeys)
		for _, k := range deleteKeys {
			actions = append(actions, Action{Key: k, Type: ActionDelete})
		}
	}

	return actions
}

// execute runs planned actions against the store.
func execute(ctx context.Context, actions []Action, s Store, dryRun bool, stdout, stderr io.Writer) Summary {
	var summary Summary
	var mu sync.Mutex
	sem := make(chan struct{}, 10)
	var wg sync.WaitGroup

	// Collect delete keys for batch operation
	var deleteKeys []string

	for i := range actions {
		a := &actions[i]
		switch a.Type {
		case ActionSkip:
			summary.Unchanged++
			continue
		case ActionDelete:
			deleteKeys = append(deleteKeys, a.Key)
			if !dryRun {
				continue // handle batch below
			}
			_, _ = fmt.Fprintf(stdout, "delete: %s\n", a.Key)
			summary.Deleted++
			continue
		case ActionCreate, ActionUpdate:
			_, _ = fmt.Fprintf(stdout, "%s: %s\n", a.Type, a.Key)
			if dryRun {
				if a.Type == ActionCreate {
					summary.Created++
				} else {
					summary.Updated++
				}
				continue
			}
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(action *Action) {
			defer wg.Done()
			defer func() { <-sem }()
			err := s.Put(ctx, action.Key, action.Value)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				action.Error = err
				summary.Failed++
				_, _ = fmt.Fprintf(stderr, "error: %s: %v\n", action.Key, err)
			} else if action.Type == ActionCreate {
				summary.Created++
			} else {
				summary.Updated++
			}
		}(a)
	}

	wg.Wait()

	// Handle deletes
	if !dryRun && len(deleteKeys) > 0 {
		for _, k := range deleteKeys {
			_, _ = fmt.Fprintf(stdout, "delete: %s\n", k)
		}
		err := s.Delete(ctx, deleteKeys)
		if err != nil {
			summary.Failed += len(deleteKeys)
			for _, k := range deleteKeys {
				_, _ = fmt.Fprintf(stderr, "error: %s: %v\n", k, err)
			}
		} else {
			summary.Deleted += len(deleteKeys)
		}
	}

	_, _ = fmt.Fprintf(stdout, "%d created, %d updated, %d deleted, %d unchanged, %d failed\n",
		summary.Created, summary.Updated, summary.Deleted, summary.Unchanged, summary.Failed)

	return summary
}
