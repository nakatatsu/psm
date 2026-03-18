# Error Handling and Logging

This document defines the error handling and logging policy for the psm project. It extends the [psm Constitution](.specify/memory/constitution.md) with operational details.

## Error Handling

Per the constitution: "Error handling follows Go conventions: return `error`, wrap with `fmt.Errorf("context: %w", err)`. Do not use `panic`."

This section extends and operationalizes that principle.

### Return Conventions

Functions that can fail MUST return `error` as the last return value.

```go
// Good
func syncParameters(ctx context.Context, entries []Entry) (Summary, error) {
    // ...
}

// Bad — caller cannot distinguish success from failure
func syncParameters(ctx context.Context, entries []Entry) Summary {
    // ...
}
```

### Panic Prohibition

`panic` MUST NOT be used anywhere in the codebase.

When an unrecoverable error occurs, return the error to the caller and let `main()` handle the exit.

```go
// Good
func run(cfg Config) (int, error) {
    if cfg.Store == "" {
        return 1, fmt.Errorf("--store is required")
    }
    // ...
}

// Bad — panic crashes without cleanup
func run(cfg Config) {
    if cfg.Store == "" {
        panic("--store is required")
    }
}
```

### Error Wrapping

Errors MUST be wrapped with context using `fmt.Errorf` with the `%w` verb when the caller adds meaningful context. Wrap MUST describe what the current function was trying to do. Avoid redundant wrapping at pass-through layers where no new context is added.

```go
// Good — adds meaningful context
existing, err := store.GetAll(ctx)
if err != nil {
    return 1, fmt.Errorf("failed to get existing parameters: %w", err)
}

// Good — pass-through is OK when no useful context to add
func (s *SSMStore) GetAll(ctx context.Context) (map[string]string, error) {
    // ... multiple API calls with individual wrapping ...
    return result, nil
}

// Bad — no context, caller cannot tell where the error originated
if err != nil {
    return 1, err
}

// Bad — uses %v instead of %w, breaks errors.Is/errors.As chain
if err != nil {
    return 1, fmt.Errorf("failed: %v", err)
}

// Bad — redundant wrapping that adds noise without value
if err != nil {
    return fmt.Errorf("GetAll failed: %w", err) // "GetAll" is obvious from the stack
}
```

### Error Classification

Use `errors.Is` and `errors.As` for control flow decisions. Define sentinel errors or typed errors when callers need to distinguish error kinds.

```go
// Sentinel error — for simple "is it this kind of error?" checks
var ErrConflict = errors.New("deletion candidate exists in sync YAML")

// Usage
if errors.Is(err, ErrConflict) {
    // abort entire operation
}
```

```go
// Typed error — when the error carries structured data
type ValidationError struct {
    Key    string
    Reason string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", e.Key, e.Reason)
}

// Usage
var ve *ValidationError
if errors.As(err, &ve) {
    slog.Error("invalid entry", "key", ve.Key, "reason", ve.Reason)
}
```

Use sentinel errors for simple conditions. Use typed errors when the caller needs access to structured context. Do not define error types speculatively — add them when a caller actually needs to branch on the error kind.

### Exit Codes

- **0**: All operations succeeded
- **1**: One or more operations failed, or invalid input

`main()` MUST be the only place that calls `os.Exit()`. All other functions MUST return errors to their callers.

## Logging

### Logger

`log/slog` (Go standard library) MUST be the only logger used. All other logging methods are prohibited:

- MUST NOT use `log.Println`, `log.Printf`, or any `log` package functions
- MUST NOT use `fmt.Fprintf(os.Stderr, ...)` for log-like messages
- MUST NOT use third-party loggers (zerolog, zap, logrus, etc.)

Initialize the default logger in `main()` using `slog.SetDefault()`. Use the package-level functions (`slog.Info`, `slog.Warn`, `slog.Error`, `slog.Debug`) throughout the codebase. Constructor injection (passing `*slog.Logger` to functions) is not needed for a CLI tool.

```go
// In main()
slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
    Level: slog.LevelInfo,
})))

// Anywhere in the codebase
slog.Info("sync complete", "created", summary.Created, "updated", summary.Updated)
slog.Error("failed to put parameter", "key", key, "error", err)
```

### Log Levels

The default log level MUST be Info (Error + Warn + Info visible; Debug hidden). Debug logging MUST be enabled via the `--debug` CLI flag.

| Level | When to Use | Example |
|-------|-------------|---------|
| Error | Operation failed for a specific item — the item is skipped or the operation aborts | `slog.Error("failed to put parameter", "key", "/myapp/prod/DB_HOST", "error", err)` |
| Warn  | Non-fatal issue that deserves attention but does not block the operation | `slog.Warn("unmanaged key detected", "key", "/myapp/legacy/OLD_KEY")` |
| Info  | Normal operational milestone — confirms what happened | `slog.Info("sync complete", "created", 3, "updated", 1, "deleted", 0)` |
| Debug | Diagnostic detail for troubleshooting — hidden by default | `slog.Debug("API call", "operation", "PutParameter", "key", "/myapp/prod/DB_HOST")` |

**Decision guide**: If the user needs to act on it, use Warn. If the operation cannot continue, use Error. If it's a normal progress update, use Info. If only a developer debugging would care, use Debug.

### Log Format

Log output MUST use text format (`slog.NewTextHandler`). JSON format (`slog.JSONHandler`) MUST NOT be added unless a concrete need arises.

### Sensitive Data

Secret values MUST NOT appear in log messages at any level. This includes:

- Parameter/secret values stored in AWS
- API keys, passwords, tokens
- Any data that would compromise security if exposed in logs

Key paths (e.g., `/myapp/prod/API_KEY`) are NOT considered sensitive and MAY be logged — they describe the location, not the value.

```go
// Good — logs the key path, not the value
slog.Error("failed to put parameter", "key", key, "error", err)

// Bad — logs the secret value
slog.Error("failed to put parameter", "key", key, "value", value, "error", err)
```

### Output Routing

Program output and log messages MUST be routed to separate streams:

**stdout** — Program results only. This is the "data" output that can be piped to other tools.

- Diff lines (`create: /myapp/prod/DB_HOST`, `update: /myapp/prod/DB_PORT`)
- Summary lines (`3 created, 1 updated, 0 deleted, 2 unchanged, 0 failed`)
- Action plans (list of planned operations for approve prompt)
- Export output (YAML content)

**stderr** — Everything else.

- All slog messages (Error, Warn, Info, Debug)
- Interactive prompts (e.g., approve `y/N` confirmation)
- Progress indicators

```go
// Good — program output to stdout
fmt.Fprintf(os.Stdout, "create: %s\n", key)

// Good — log message to stderr (slog writes to stderr by default)
slog.Warn("unmanaged key detected", "key", key)

// Good — interactive prompt to stderr
fmt.Fprintf(os.Stderr, "Proceed? [y/N] ")

// Bad — log message to stdout (breaks piping)
fmt.Fprintf(os.Stdout, "WARNING: unmanaged key %s\n", key)
```

This separation ensures `psm sync ... | grep create` works correctly — only diff lines appear in the pipe, not log noise or prompts.
