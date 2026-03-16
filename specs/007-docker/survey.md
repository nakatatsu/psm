# Survey: Reference Dockerfile for psm

**Date**: 2026-03-15
**Spec**: specs/007-docker/spec.md

## Summary

Dockerfile approach is valid for CI/CD use case. Key findings: (1) psm uses Go AWS SDK internally — AWS CLI binary is a debugging convenience, not a hard dependency. (2) debian:bookworm-slim is the right base (alpine breaks AWS CLI installer). (3) All 3 tools have arm64 binaries. (4) Maintenance burden of 3 external tool versions is the main risk.

## S1: Does psm need AWS CLI?

**Category**: Problem Definition
**Finding**: psm uses aws-sdk-go-v2 directly. It does not shell out to `aws` CLI. AWS CLI is only useful for debugging (`aws sts get-caller-identity`, `aws ssm get-parameter`, etc.). It adds ~100MB to the image.
**Recommendation**: Include AWS CLI as it's part of the stated use case (SOPS + AWS CLI + psm as a set), but note in docs that psm itself doesn't require it.

## S2: Base image

**Category**: Approach Alternatives
**Finding**: alpine breaks AWS CLI v2 official installer (requires glibc). scratch has no shell (can't pipe sops output). amazon/aws-cli is bloated and designed as entrypoint-based.
**Recommendation**: debian:bookworm-slim. Same as DevContainer base, proven to work.

## S3: Multi-platform

**Category**: Feasibility
**Finding**: All 3 tools (psm, SOPS, AWS CLI v2) have linux/amd64 and linux/arm64 support. Dockerfile can use TARGETARCH for platform-aware downloads.
**Recommendation**: Support both amd64 and arm64 via TARGETARCH.

## S4: Maintenance risk

**Category**: Risk
**Finding**: 3 external tools with independent release cycles. Dockerfile rots without active version bumping.
**Recommendation**: Pin all versions as ARG. Document that users should check for updates. This is a reference, not a maintained product.

## Items Requiring PoC

- SOPS latest version number (need web search in plan phase)

## Constitution Impact

None.

## Recommendation

Proceed to plan. Straightforward Dockerfile + README addition.
