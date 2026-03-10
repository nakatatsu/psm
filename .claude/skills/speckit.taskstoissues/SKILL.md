---
description: Read tasks.md and create GitHub issues for each task with proper labels, dependencies, and file paths. Use this after tasks have been generated to push them into GitHub for team tracking and project management.
allowed-tools: github/github-mcp-server/issue_write
---

## User Input

```text
$ARGUMENTS
```

Consider the user input before proceeding (if not empty).

## Outline

1. Run `.specify/scripts/bash/check-prerequisites.sh --json --require-tasks --include-tasks` from repo root and parse FEATURE_DIR and AVAILABLE_DOCS list. All paths must be absolute. For single quotes in args like "I'm Groot", use escape syntax: e.g 'I'\''m Groot' (or double-quote if possible: "I'm Groot").

2. Extract the path to **tasks.md** from the script output and read it.

3. Get the Git remote:

```bash
git config --get remote.origin.url
```

Only proceed if the remote is a GitHub URL. If it is not, stop and tell the user this command only works with GitHub repositories.

4. For each task in tasks.md, use the GitHub MCP server to create a new issue in the repository matching the Git remote. Never create issues in a repository that does not match the remote URL -- this prevents accidentally polluting other repos.

   Each issue should include:
   - **Title**: The task description from tasks.md (e.g., "T005 Implement authentication middleware")
   - **Body**: The full task details including file paths, dependencies on other tasks (referencing their issue numbers), the phase it belongs to, and any parallel [P] or story [US*] labels from tasks.md
   - **Labels**: Apply labels based on task metadata -- phase name, story label (if present), and a "parallel" label for [P] tasks where appropriate
