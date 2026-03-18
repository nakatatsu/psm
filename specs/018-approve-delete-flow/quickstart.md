# Quickstart: Add Approve Flow and Replace --prune with --delete

## Basic Sync (with approval)

```bash
# Displays plan, prompts y/N, then executes
psm sync --store ssm params.yml

# stdout:
# create: /myapp/prod/NEW_KEY
# update: /myapp/prod/CHANGED_KEY
# 1 created, 1 updated, 0 deleted, 3 unchanged, 0 failed

# stderr:
# Proceed? [y/N]
```

## Dry-Run (no prompt, no execution)

```bash
psm sync --store ssm --dry-run params.yml

# stdout:
# (dry-run) create: /myapp/prod/NEW_KEY
# (dry-run) update: /myapp/prod/CHANGED_KEY
# 1 created, 1 updated, 0 deleted, 3 unchanged, 0 failed (dry-run)
```

## CI/Automation (skip approval)

```bash
psm sync --store ssm --skip-approve params.yml
```

## Delete Obsolete Keys

```bash
# Create delete patterns file
cat > needless.yml << 'EOF'
- "^/myapp/legacy/"
- "^/myapp/deprecated-.*"
EOF

# Sync + delete (displays plan with deletions, prompts approval)
psm sync --store ssm --delete needless.yml params.yml

# stdout:
# create: /myapp/prod/NEW_KEY
# delete: /myapp/legacy/OLD_KEY
# delete: /myapp/deprecated-feature-x
# 1 created, 0 updated, 2 deleted, 3 unchanged, 0 failed

# stderr:
# time=... level=WARN msg="unmanaged key detected" key=/other/team/THEIR_KEY
# Proceed? [y/N]
```

## Conflict Detection (abort)

```bash
# If /myapp/prod/API_KEY exists in both params.yml and matches a delete pattern:
psm sync --store ssm --delete needless.yml params.yml

# stderr:
# error: conflict detected: deletion candidate /myapp/prod/API_KEY exists in sync YAML
# error: aborting — no changes made

# exit code: 1
```

## Debug Logging

```bash
psm sync --store ssm --debug params.yml

# stderr includes:
# time=... level=DEBUG msg="API call" operation=GetParametersByPath path=/
# time=... level=DEBUG msg="regex match" pattern="^/myapp/legacy/" key=/myapp/legacy/OLD matched=true
```

## Migration from --prune

```bash
# Old (removed):
psm sync --store ssm --prune params.yml
# error: --prune has been removed. Use --delete <file> with regex patterns instead.

# New:
echo '- ".*"' > delete-all.yml
psm sync --store ssm --delete delete-all.yml params.yml
# (Note: ".*" matches everything — equivalent to old --prune but explicit and reviewable)
```

## Verification Checklist

- [ ] `psm sync --store ssm params.yml` shows plan and prompts before executing
- [ ] Typing `N` or Enter at prompt cancels without changes (exit 0)
- [ ] Typing `y` executes all planned changes
- [ ] `--skip-approve` executes immediately without prompt
- [ ] `--dry-run` shows plan without prompt or execution
- [ ] `--prune` returns error with migration message
- [ ] `--delete needless.yml` deletes only matching keys not in sync YAML
- [ ] Conflict (key in both sync YAML and delete match) aborts entirely
- [ ] Unmanaged keys shown as warnings
- [ ] `--debug` shows debug-level slog output
- [ ] Non-terminal stdin without `--skip-approve` exits 0 (no changes)
- [ ] No changes case: summary shown, no prompt
