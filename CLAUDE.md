# workspace Development Guidelines

## MANDATORY

NEVER RUN `sudo /usr/local/bin/init-firewall.sh`! IT DESTROY SESSION!

Always save `.claude/` files (skills, settings, etc.) in the repository (`/workspace/.claude/`), never in `~/.claude/`. Files in `~/.claude/` are lost when the container is destroyed.

## Active Technologies

- Go 1.26 (1.26.1) (001-psm)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for Go 1.26 (1.26.1)

## Code Style

Go 1.26 (1.26.1): Follow standard conventions

## Recent Changes

- 001-psm: Added Go 1.26 (1.26.1)

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
