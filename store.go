package main

import (
	"context"
	"io"
)

// ActionType represents the type of sync operation.
type ActionType int

const (
	ActionCreate ActionType = iota
	ActionUpdate
	ActionDelete
	ActionSkip
)

func (a ActionType) String() string {
	switch a {
	case ActionCreate:
		return "create"
	case ActionUpdate:
		return "update"
	case ActionDelete:
		return "delete"
	case ActionSkip:
		return "skip"
	default:
		return "unknown"
	}
}

// Config holds CLI flags and arguments.
type Config struct {
	Subcommand  string
	Profile     string
	DryRun      bool
	SkipApprove bool
	Debug       bool
	DeleteFile  string
	File        string
	ShowVersion bool
}

// Entry is a single key-value pair from YAML.
type Entry struct {
	Key   string
	Value string
}

// Action represents a planned or executed sync operation.
type Action struct {
	Key   string
	Type  ActionType
	Value string // needed for Put operations
	Error error
}

// Summary holds aggregate counts after sync execution.
type Summary struct {
	Created   int
	Updated   int
	Deleted   int
	Unchanged int
	Failed    int
}

// Store abstracts AWS parameter/secret storage.
type Store interface {
	GetAll(ctx context.Context) (map[string]string, error)
	Put(ctx context.Context, key, value string) error
	Delete(ctx context.Context, keys []string) error
}

// IOStreams holds I/O dependencies for testability.
type IOStreams struct {
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	IsTerminal func() bool
}
