# Survey: Add Approve Flow and Replace --prune with --delete

**Date**: 2026-03-17
**Spec**: specs/018-approve-delete-flow/spec.md

## Summary

The spec correctly identifies two real problems: no approval gate and an overly broad `--prune`. The proposed approach (approve prompt + regex-based delete file) is sound and well-scoped. Key findings: (1) the existing plan/execute separation in the code naturally supports inserting an approval step, (2) terminal detection for non-interactive environments can be done with stdlib only, (3) Go's regex engine (RE2) has limitations that should be documented, (4) the Config struct and parseArgs need significant rework to handle new flags. No fundamental problems with the direction.

## S1: Problem Validation — Are Approve and Delete-Replacement the Right Solutions?

**Category**: Problem Reframing — Problem Definition
**Finding**: The problems are real and well-defined:
- **No approval**: `execute()` in `sync.go:46` runs changes immediately after `plan()`. There is no step between planning and execution. Any operator mistake is irreversible.
- **Dangerous prune**: `plan()` in `sync.go:29-40` deletes every AWS key not in YAML. For an account with 500 parameters from 10 different teams, one YAML file could wipe 490 parameters.

Both problems are "one mistake = disaster" class. The spec's approach (insert approval between plan/execute, replace broad prune with scoped regex) addresses root causes, not symptoms.

**Recommendation**: Proceed as specified. No reframing needed.
**Evidence**: Code review of `sync.go:12-43` (plan) and `sync.go:46-132` (execute) — no safety gate exists between them.

## S2: Architecture Fit — Plan/Execute Separation Already Exists

**Category**: Solution Evaluation — Approach Alternatives
**Finding**: The current architecture already separates planning from execution:
- `plan()` returns `[]Action` (sync.go:12)
- `execute()` takes `[]Action` and runs them (sync.go:46)

Inserting an approval step between these two is a natural extension — display the action list, prompt, then call execute. No architectural refactoring required.

However, `execute()` currently handles both display and execution in a single pass. The approval flow needs to display first (all actions), then prompt, then execute (non-skip actions only). This means the "display" logic needs to be extracted from `execute()` or the approval display needs to be a separate step before `execute()`.

**Recommendation**: Add a `displayPlan()` function that renders the action list to stdout. Call it before the approval prompt. Keep `execute()` for the actual execution. This preserves the existing architecture while adding the new capability.
**Evidence**: `sync.go:60-84` — display and execution are interleaved in the same loop.

## S3: Terminal Detection Without New Dependencies

**Category**: Constraints and Tradeoffs
**Finding**: The spec assumes the tool can detect non-interactive environments (piped stdin). Two approaches:

| Approach | Dependency | Method |
|----------|------------|--------|
| `golang.org/x/term` | Third-party (golang.org/x) | `term.IsTerminal(int(os.Stdin.Fd()))` |
| stdlib `os.Stat` | None | `fi, _ := os.Stdin.Stat(); fi.Mode()&os.ModeCharDevice != 0` |

The stdlib approach works on Linux and macOS. It avoids a new dependency per Constitution Principle I.

**Recommendation**: Use `os.Stdin.Stat()` with `os.ModeCharDevice` check. No new dependency needed.
**Evidence**: Go stdlib documentation for `os.FileInfo.Mode()`. `os.ModeCharDevice` is set when stdin is a terminal. Widely used pattern in Go CLI tools.

## S4: Go Regex (RE2) Limitations

**Category**: Risk & Failure Modes
**Finding**: Go's `regexp` package uses RE2, which does not support:
- Lookahead/lookbehind (`(?=...)`, `(?!...)`)
- Backreferences (`\1`)
- Possessive quantifiers

Operators familiar with PCRE (Perl, Python, JavaScript regex) may write patterns that fail to compile. FR-010 already handles this (invalid regex = immediate error), but the error message should be user-friendly and mention RE2 limitations.

**Recommendation**: Document RE2 limitations in README. Ensure error messages from `regexp.Compile` are surfaced clearly with the offending pattern name. No spec change needed.
**Evidence**: Go `regexp` package documentation: "The syntax of the regular expressions accepted is the same general syntax used by Perl, Python, and other languages. More precisely, it is the syntax accepted by RE2."

## S5: Config Struct and parseArgs Rework

**Category**: Integration Impact
**Finding**: The current `Config` struct (store.go:31-38) and `parseArgs` (main.go:88-134) need significant changes:

Current Config:
```go
type Config struct {
    Subcommand string
    Store      string
    Profile    string
    Prune      bool      // → remove
    DryRun     bool
    File       string    // single file
}
```

New Config needs:
- Remove `Prune`
- Add `DeleteFile string` (path to delete patterns YAML)
- Add `SkipApprove bool`
- Add `Debug bool`
- `File` remains the sync YAML (positional argument)

`parseArgs` currently enforces "exactly one file argument" (main.go:121-123). With `--delete`, this remains true — the delete file is a flag value, not a positional argument. No ambiguity.

**Recommendation**: Straightforward struct modification. No architectural concern.
**Evidence**: `main.go:88-134`, `store.go:31-38`.

## S6: Existing Test Impact

**Category**: Integration Impact
**Finding**: Multiple test files reference `Prune`:
- `main_test.go:108` — `got.Prune` assertion
- `main_test.go:25-38` — parseArgs test cases with `--prune`
- `sync_test.go:163,174` — plan function tests with prune=true/false
- `ssm_test.go:231,277` — integration tests for prune behavior

All of these need to be updated: prune test cases removed, new test cases for `--delete`, `--skip-approve`, `--debug`, and approval flow added.

**Recommendation**: Plan for comprehensive test updates. The test-first constitution requirement means new tests must be written before implementation.
**Evidence**: `grep -r "Prune\|prune" *_test.go` across all test files.

## S7: Delete File Parsing — Reuse or New Parser?

**Category**: Solution Evaluation — Cost & Complexity
**Finding**: The delete file is a YAML list of strings. The existing `parseYAML()` in `yaml.go` parses YAML mappings (key-value pairs), not lists. A new parser is needed for the delete file format.

The parsing is simple: `yaml.Unmarshal(data, &[]string{})` directly. Then compile each string as a regex with `regexp.Compile`. No need for the `yaml.Node` API used in `parseYAML`.

**Recommendation**: Add a `parseDeletePatterns(data []byte) ([]*regexp.Regexp, error)` function. Keep it separate from `parseYAML` — different input format, different output type.
**Evidence**: `yaml.go:11-89` — existing parser handles mappings, not lists.

## Items Requiring PoC

None. All components use well-understood Go stdlib features (`regexp`, `os.Stdin.Stat`, `bufio.Scanner` for prompt input, `flag` for new flags).

## Constitution Impact

No amendments required. The feature uses:
- Standard library only (Principle I)
- No speculative features (Principle II)
- Test-first approach applicable (Principle III)

## Recommendation

Proceed to `/speckit.plan`. No blockers identified. Key points for planning:

1. **S2**: Extract display logic from `execute()` into a separate `displayPlan()` function
2. **S3**: Use `os.Stdin.Stat()` for terminal detection — no new dependencies
3. **S5**: Config struct rework is straightforward
4. **S6**: Plan for comprehensive test updates (test-first)
5. **S7**: New `parseDeletePatterns()` function, separate from existing YAML parser
