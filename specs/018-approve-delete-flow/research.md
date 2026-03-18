# Research: Add Approve Flow and Replace --prune with --delete

**Date**: 2026-03-17

## R1: Approval Prompt Implementation

**Decision**: Use `fmt.Fprintf(os.Stderr, ...)` for the prompt text and `bufio.NewScanner(os.Stdin)` for reading user input. Accept `y` or `Y` as approval, treat everything else (including empty input/Enter) as decline.
**Rationale**: stdlib only. The prompt goes to stderr per the logging policy (stdout is for program output). `bufio.Scanner` is simpler than `fmt.Scanln` for reading a single line with potential empty input.
**Alternatives**:
- `fmt.Scanln` — doesn't handle empty input well (blocks on Enter without text).
- Third-party prompt library (e.g., `survey`, `promptui`) — violates Constitution Principle I.

## R2: Terminal Detection

**Decision**: Use `os.Stdin.Stat()` with `os.ModeCharDevice` bit check to detect if stdin is a terminal.
**Rationale**: stdlib only, no new dependencies. When stdin is not a terminal (piped/redirected) and `--skip-approve` is not set, treat as decline (exit 0, no changes). Survey S3 confirmed this approach.

```go
fi, err := os.Stdin.Stat()
if err != nil || fi.Mode()&os.ModeCharDevice == 0 {
    // not a terminal — treat as declined
}
```

**Alternatives**:
- `golang.org/x/term` — `term.IsTerminal()`. Adds a dependency, rejected per Constitution Principle I.

## R3: Delete Pattern File Parsing

**Decision**: New `parseDeletePatterns(data []byte) ([]*regexp.Regexp, error)` function. Uses `yaml.Unmarshal` into `[]string`, then compiles each string with `regexp.Compile`. Returns compiled patterns for efficient matching.
**Rationale**: Simple YAML list → compiled regex. Separate from `parseYAML` (which handles key-value mappings). Validates all patterns before returning — if any fails to compile, returns error with the offending pattern (FR-010). Survey S7 confirmed this approach.
**Alternatives**: None — straightforward implementation.

## R4: Conflict Detection Strategy

**Decision**: Pre-validation function `detectConflicts(candidates []string, yamlKeys map[string]bool) []string` that returns conflicting keys. Called before any execution. If result is non-empty, abort with error listing all conflicts.
**Rationale**: All-or-nothing validation per FR-009. Must check ALL candidates before executing ANY operation. The function is pure (no side effects) and easily testable.
**Alternatives**: None — the spec mandates this exact behavior.

## R5: Plan Function Signature Change

**Decision**: Change `plan()` to no longer accept `prune bool`. Delete planning moves to a separate `planDeletes(existing map[string]string, yamlKeys map[string]bool, patterns []*regexp.Regexp) ([]Action, []string, []string)` returning (delete actions, conflict keys, unmanaged keys).
**Rationale**: Separation of concerns. Sync planning (create/update/skip) and delete planning (regex matching, conflict detection, unmanaged key identification) are distinct operations with different inputs. Survey S2 confirmed the plan/execute separation supports this.
**Alternatives**: Keep everything in `plan()` — rejected because it would make the function too complex and hard to test independently.

## R6: Display Plan Function

**Decision**: New `displayPlan(actions []Action, stdout io.Writer)` function that renders the action list without executing. Called before the approval prompt. The existing `execute()` function no longer needs to handle display for non-dry-run cases.
**Rationale**: Survey S2 identified that display and execution are currently interleaved in `execute()`. Separating display enables the approve flow: display → prompt → execute. `--dry-run` uses `displayPlan()` only (no prompt, no execute).
**Alternatives**: Modify `execute()` to accept a "display-only" flag — rejected because it adds complexity to an already complex function.

## R7: slog Integration

**Decision**: Initialize `slog.SetDefault()` in `main()` with `slog.LevelInfo` default. When `--debug` is passed, set level to `slog.LevelDebug`. Use `slog.Debug` for diagnostic output (API calls, regex matches), `slog.Warn` for unmanaged keys, `slog.Error` for failures.
**Rationale**: Per the logging policy document (#19). Text format via `slog.NewTextHandler` to stderr.
**Alternatives**: None — mandated by logging policy.

## R8: Config Struct Changes

**Decision**: Update `Config` struct:
- Remove `Prune bool`
- Add `DeleteFile string` (path to delete patterns YAML)
- Add `SkipApprove bool`
- Add `Debug bool`

`parseArgs` changes:
- Remove `--prune` flag registration, add error message if detected in args
- Add `--delete <file>` string flag (sync subcommand only)
- Add `--skip-approve` bool flag (sync subcommand only)
- Add `--debug` bool flag (all subcommands)
- Keep "exactly one file argument" validation — `--delete` is a flag, not positional

**Rationale**: Survey S5 confirmed this is straightforward. No architectural concern.
**Alternatives**: None.
