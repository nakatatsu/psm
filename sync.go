package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
)

// plan computes actions by comparing YAML entries against existing AWS state.
func plan(entries []Entry, existing map[string]string) []Action {
	var actions []Action

	for _, e := range entries {
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

	return actions
}

// displayPlan renders the action list to stdout without executing.
func displayPlan(actions []Action, stdout io.Writer) {
	for _, a := range actions {
		switch a.Type {
		case ActionSkip:
			continue
		default:
			fmt.Fprintf(stdout, "%s: %s\n", a.Type, a.Key)
		}
	}
}

// printSummary outputs the summary line for an action list.
func printSummary(actions []Action, dryRun bool, stdout io.Writer) {
	var created, updated, deleted, unchanged int
	for _, a := range actions {
		switch a.Type {
		case ActionCreate:
			created++
		case ActionUpdate:
			updated++
		case ActionDelete:
			deleted++
		case ActionSkip:
			unchanged++
		}
	}
	suffix := ""
	if dryRun {
		suffix = " (dry-run)"
	}
	fmt.Fprintf(stdout, "%d created, %d updated, %d deleted, %d unchanged, 0 failed%s\n",
		created, updated, deleted, unchanged, suffix)
}

// promptApprove asks the user for confirmation. Returns true only for "y" or "Y".
func promptApprove(reader io.Reader, writer io.Writer) bool {
	fmt.Fprint(writer, "Proceed? [y/N] ")
	scanner := bufio.NewScanner(reader)
	if !scanner.Scan() {
		return false
	}
	input := strings.TrimSpace(scanner.Text())
	return input == "y" || input == "Y"
}

// execute runs planned actions against the store. Display is handled by displayPlan before calling this.
func execute(ctx context.Context, actions []Action, s Store, stdout, stderr io.Writer) Summary {
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
			continue
		case ActionCreate, ActionUpdate:
			// Put operations run concurrently
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
				fmt.Fprintf(stderr, "error: %s: %v\n", action.Key, err)
			} else if action.Type == ActionCreate {
				summary.Created++
			} else {
				summary.Updated++
			}
		}(a)
	}

	wg.Wait()

	// Handle deletes
	if len(deleteKeys) > 0 {
		err := s.Delete(ctx, deleteKeys)
		if err != nil {
			summary.Failed += len(deleteKeys)
			for _, k := range deleteKeys {
				fmt.Fprintf(stderr, "error: %s: %v\n", k, err)
			}
		} else {
			summary.Deleted += len(deleteKeys)
		}
	}

	fmt.Fprintf(stdout, "%d created, %d updated, %d deleted, %d unchanged, %d failed\n",
		summary.Created, summary.Updated, summary.Deleted, summary.Unchanged, summary.Failed)

	return summary
}
