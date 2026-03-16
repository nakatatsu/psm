# Output Format Checklist: Distinguish dry-run output

**Purpose**: Validate completeness and consistency of spec/plan requirements
**Created**: 2026-03-16
**Feature**: specs/008-dryrun-output/spec.md

## Requirement Completeness

- [ ] CHK001 Is the dry-run prefix format explicitly specified? [Completeness, Plan Design]
- [ ] CHK002 Is the summary line format for dry-run explicitly specified? [Completeness, Plan Design]
- [ ] CHK003 Are all action types (create, update, delete, skip) covered for dry-run output? [Completeness, Spec FR-002]
- [ ] CHK004 Is backward compatibility for non-dry-run output explicitly required? [Completeness, Spec FR-003]

## Requirement Clarity

- [ ] CHK005 Is the exact output format unambiguous (prefix position, parentheses, spacing)? [Clarity, Plan]
- [ ] CHK006 Is the skip action's dry-run behavior defined (skip produces no output in either mode)? [Clarity, Gap]

## Requirement Consistency

- [ ] CHK007 Is the prefix format consistent across all action types and summary? [Consistency, Plan]
- [ ] CHK008 Is the plan's output format consistent with the spec's FR-001 and FR-002? [Consistency, Plan vs Spec]

## Acceptance Criteria Quality

- [ ] CHK009 Are acceptance scenarios testable with string matching? [Measurability, Spec AS 1-3]
- [ ] CHK010 Is the non-dry-run backward compatibility testable? [Measurability, Spec AS-2]

## Scenario Coverage

- [ ] CHK011 Is the --prune + --dry-run combination covered? [Coverage, Spec Edge Cases]
- [ ] CHK012 Is the zero-changes dry-run scenario covered? [Coverage, Spec AS-3]

## Notes

- 12 items total
- Small, focused change with clear requirements
