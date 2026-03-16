# Starter Template Checklist: psm Example Project

**Purpose**: Validate completeness, clarity, and consistency of spec/plan requirements for the example project template
**Created**: 2026-03-16
**Feature**: specs/007-docker/spec.md

## Requirement Completeness

- [ ] CHK001 Are all required tools (psm, SOPS, age, AWS CLI v2) explicitly listed with installation methods in the plan? [Completeness, Spec FR-002, Plan]
- [ ] CHK002 Is the full file listing of `example/` specified with each file's purpose? [Completeness, Spec FR-001, Plan]
- [ ] CHK003 Are the DevContainer configuration requirements (devcontainer.json fields) enumerated? [Completeness, Spec FR-002, Plan]
- [ ] CHK004 Is the content structure of `secrets.yaml` (sample format) defined? [Completeness, Spec FR-003]
- [ ] CHK005 Is the content structure of `.sops.yaml.example` defined? [Completeness, Spec FR-003, Survey S4]
- [ ] CHK006 Are all README sections (key generation → encryption → SSO → sync) specified? [Completeness, Spec FR-005]
- [ ] CHK007 Is the `test.sh` content (manual verification commands) enumerated? [Completeness, Plan]

## Requirement Clarity

- [ ] CHK008 Is "lightweight" in NFR-001 quantified — which tools are excluded and why? [Clarity, Spec NFR-001]
- [ ] CHK009 Is the AWS credential mounting strategy clearly specified (host `~/.aws` mount vs environment variables)? [Clarity, Survey S6, Plan]
- [ ] CHK010 Is the `.sops.yaml.example` → `.sops.yaml` rename workflow clearly described? [Clarity, Survey S4]
- [ ] CHK011 Are the expected psm secrets YAML format requirements specified (flat key-value, path prefixes)? [Clarity, Survey S5]

## Requirement Consistency

- [ ] CHK012 Are tool versions in the plan consistent with those in the existing `docker/Dockerfile`? [Consistency, Plan vs docker/Dockerfile]
- [ ] CHK013 Is the base image choice consistent between plan and survey? [Consistency, Plan vs Survey S2]
- [ ] CHK014 Is the single-container decision in the plan consistent with the survey recommendation? [Consistency, Plan vs Survey S3]

## Acceptance Criteria Quality

- [ ] CHK015 Is acceptance scenario 1 (DevContainer opens, all tools work) independently verifiable? [Measurability, Spec AS-1]
- [ ] CHK016 Is acceptance scenario 2 (age key + sops encrypt) independently verifiable without AWS? [Measurability, Spec AS-2]
- [ ] CHK017 Is acceptance scenario 3 (sops decrypt + psm sync) specific about which Parameter Store path to verify? [Measurability, Spec AS-3]

## Scenario Coverage

- [ ] CHK018 Are requirements defined for the "copy to existing repo" scenario (merging into existing project)? [Coverage, Spec Edge Case 3]
- [ ] CHK019 Are error message requirements defined for missing age key and missing AWS auth scenarios? [Coverage, Spec Edge Cases 1-2]
- [ ] CHK020 Is the "independent from psm dev repo" requirement (SC-003) testable with specific criteria? [Coverage, Spec SC-003]

## Dependencies and Assumptions

- [ ] CHK021 Is the assumption that `aws sso login` works inside DevContainer validated or flagged for PoC? [Assumption, Survey PoC]
- [ ] CHK022 Is the assumption that users have AWS SSO configured documented as a prerequisite? [Assumption, Spec Assumptions]

## Notes

- Check items off as completed: `[x]`
- Focus areas: Completeness and Clarity (standard depth, reviewer audience)
- This is a template/reference project; no Go code involved
