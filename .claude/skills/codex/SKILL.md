---
name: codex
description: >
  Run OpenAI Codex CLI to delegate coding tasks.
  Invoke with /codex followed by a prompt.
  Use when the user wants to use Codex for a task, says "codex", or wants to compare Claude and Codex outputs.
---

# codex — OpenAI Codex CLI

Run OpenAI Codex CLI from within Claude Code to delegate tasks to Codex.

## Prerequisites

- `OPENAI_API_KEY` environment variable must be set
- Codex CLI must be installed (included in DevContainer)

## Usage

| Command                  | What it does                        |
| ------------------------ | ----------------------------------- |
| `/codex <prompt>`        | Run Codex with the given prompt     |
| `/codex`                 | Show usage help                     |

## Security Notice

`--dangerously-bypass-approvals-and-sandbox` フラグにより Codex はサンドボックスなしで実行されます。
**ユーザーの明示的な許可なしに `/codex` を実行してはいけません。**
CLAUDE.md でプロジェクト単位の事前承認が記載されている場合のみ、許可なしで実行できます。

## Execution

1. **Check prerequisites**: Verify `codex` command is available and `OPENAI_API_KEY` is set.
   - If `codex` is not found: ERROR "Codex CLI がインストールされていません。DevContainer を rebuild してください。"
   - If `OPENAI_API_KEY` is not set: ERROR "OPENAI_API_KEY が設定されていません。`.env` ファイルまたは環境変数で設定してください。"

2. **Parse arguments**: The text after `/codex` is the prompt to pass to Codex.
   - If no prompt is provided, display the usage table above and exit.

3. **Run Codex**: Execute the Codex CLI in non-interactive mode using the Bash tool:
   ```bash
   codex exec --dangerously-bypass-approvals-and-sandbox "<prompt>"
   ```
   - `--dangerously-bypass-approvals-and-sandbox` is required because DevContainer does not support unprivileged user namespaces (bwrap sandbox fails). The DevContainer itself provides isolation, so this is safe.
   - Set a generous timeout (e.g. 300000ms) as Codex may take time to respond.

4. **Report results**: Show the Codex output to the user. If Codex returns an error, display the error message clearly.
