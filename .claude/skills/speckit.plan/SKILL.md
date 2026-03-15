---
description: Execute the implementation planning workflow to produce design artifacts -- research.md, data-model.md, contracts/, and quickstart.md. Use this after a spec is written and the user wants to move from requirements to technical design.
---

## User Input

```text
$ARGUMENTS
```

Consider the user input before proceeding (if not empty).

## Outline

1. **Setup**: Run `.specify/scripts/bash/setup-plan.sh --json` from repo root and parse the JSON output for FEATURE_SPEC, IMPL_PLAN, SPECS_DIR, and BRANCH. Use absolute paths throughout. For single quotes in args, use: `"I'm Groot"` (double-quote) or `'I'\''m Groot'` (escape).

2. **Load context**: Read FEATURE_SPEC and `.specify/memory/constitution.md`. Load the IMPL_PLAN template (already copied by the setup script).

3. **Execute the plan workflow** following the IMPL_PLAN template structure:
   - Fill Technical Context, marking unknowns as "NEEDS CLARIFICATION".
   - Fill the Constitution Check section from the constitution.
   - Evaluate gates -- ERROR if violations are unjustified.
   - Phase 0: generate research.md (resolve all NEEDS CLARIFICATION items).
   - Phase 1: generate data-model.md, contracts/, and quickstart.md.
   - Phase 1: update agent context by running the agent script.
   - Re-evaluate the Constitution Check post-design.

4. **Stop and report**: This command ends after Phase 1 planning. Report the branch, IMPL_PLAN path, and all generated artifacts.

## Phase 0: Outline and Research

1. **Extract unknowns** from Technical Context:
   - Each NEEDS CLARIFICATION item becomes a research task.
   - Each dependency becomes a best-practices task.
   - Each integration becomes a patterns task.

2. **Generate and dispatch research tasks**:
   ```text
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **Consolidate findings** in `research.md`:
   - Decision: what was chosen
   - Rationale: why it was chosen
   - Alternatives considered: what else was evaluated

Output: research.md with all NEEDS CLARIFICATION items resolved.

## Phase 1: Design and Contracts

Prerequisites: research.md is complete.

1. **Extract entities from feature spec** into `data-model.md`:
   - Entity name, fields, relationships
   - Validation rules derived from requirements
   - State transitions if applicable

2. **Define interface contracts** (if the project has external interfaces) in `/contracts/`:
   - Identify what interfaces the project exposes to users or other systems.
   - Document contracts in the format appropriate for the project type: public APIs for libraries, command schemas for CLI tools, endpoints for web services, grammars for parsers, UI contracts for applications.
   - Skip this step if the project is purely internal (build scripts, one-off tools, etc.).

3. **Update agent context**:
   - Run `.specify/scripts/bash/update-agent-context.sh claude`.
   - The script detects the AI agent in use and updates the appropriate context file.
   - Add only new technology from the current plan; preserve manual additions between markers.

Output: data-model.md, /contracts/*, quickstart.md, agent-specific context file.

ERROR on any gate failures or unresolved clarifications that remain after research.
