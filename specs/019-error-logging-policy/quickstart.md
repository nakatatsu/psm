# Quickstart: Define Error Handling and Logging Policy

## Deliverables

This feature produces two documentation changes and zero code changes:

1. **Rewrite** `documents/backend/error-handling-and-logging.md` with the following structure:

```markdown
# Error Handling and Logging

## Relationship to Constitution
[Reference constitution line 54 as authority for core rules]

## Error Handling
### Return Conventions
[(..., error) pattern, per constitution]
### Panic Prohibition
[No panic anywhere, extends constitution's "library code" to all code]
### Error Wrapping
[fmt.Errorf with %w, context at each call site, per constitution]
### Exit Codes
[0 = success, 1 = failure]

## Logging
### Logger
[slog only, default logger pattern via slog.SetDefault]
### Log Levels
[Error/Warn/Info/Debug table with usage guidelines and examples]
### Default Level
[Info]
### Output Routing
[stdout = program output, stderr = slog + prompts]
### Log Format
[Text format via slog.TextHandler]
### Sensitive Data
[Never log secret values, tokens, passwords. Key paths are permitted]
```

2. **Delete** `documents/backend/coding-standard.md` (redundant — see research R6)

## Verification

- Read the completed document and confirm SC-001 (log level determinable in 30 seconds)
- Check SC-002 (all Issue #19 acceptance criteria covered)
- Check SC-003 (every rule uses MUST/MUST NOT with examples)
- Confirm no code changes were made
