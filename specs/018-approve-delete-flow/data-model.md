# Data Model: Add Approve Flow and Replace --prune with --delete

## Updated Entities

### Config (modified)

| Field | Type | Description |
|-------|------|-------------|
| Subcommand | string | "sync" or "export" |
| Store | string | "ssm" or "sm" |
| Profile | string | AWS profile name (optional) |
| ~~Prune~~ | ~~bool~~ | **REMOVED** |
| DryRun | bool | Show plan without executing |
| File | string | Path to sync YAML file |
| DeleteFile | string | Path to delete patterns YAML file (optional) |
| SkipApprove | bool | Bypass approval prompt |
| Debug | bool | Enable debug-level logging |

### Action (unchanged)

| Field | Type | Description |
|-------|------|-------------|
| Key | string | AWS parameter/secret key |
| Type | ActionType | Create, Update, Delete, Skip |
| Value | string | Value for Put operations |
| Error | error | Failure error (nil on success) |

### Summary (unchanged)

| Field | Type | Description |
|-------|------|-------------|
| Created | int | Keys created |
| Updated | int | Keys updated |
| Deleted | int | Keys deleted |
| Unchanged | int | Keys unchanged (skip) |
| Failed | int | Keys that failed |

## New Entities

### Delete Pattern File (input format)

A YAML file containing a list of regex pattern strings.

```yaml
- "^/myapp/legacy/"
- "^/myapp/deprecated-.*"
```

Validation rules:
- Must be a valid YAML list of strings
- Each string must be a valid Go (RE2) regular expression
- Empty list is valid (results in no deletions)

### Deletion Candidate

An AWS key that:
1. Matches at least one regex in the delete pattern file
2. Is NOT present in the sync YAML

Not a persisted entity — computed during planning.

### Unmanaged Key

An AWS key that:
1. Is NOT in the sync YAML
2. Does NOT match any regex in the delete pattern file

Displayed as warnings. Not a persisted entity — computed during planning.

## State Transitions

### Sync Flow (updated)

```
parseArgs → loadFiles → GetAll(AWS) → plan(sync) → planDeletes(if --delete)
  → detectConflicts → [ABORT if conflicts]
  → displayPlan → [approve prompt if not --skip-approve and changes exist]
  → execute (if approved) → summary
```

### Dry-Run Flow

```
parseArgs → loadFiles → GetAll(AWS) → plan(sync) → planDeletes(if --delete)
  → detectConflicts → [ABORT if conflicts]
  → displayPlan → summary (no prompt, no execute)
```
