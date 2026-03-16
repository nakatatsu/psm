# Implementation Plan: Reference Dockerfile for psm

**Branch**: `007-docker` | **Date**: 2026-03-15 | **Spec**: specs/007-docker/spec.md

## Summary

SOPS + AWS CLI v2 + psm を同梱した参考用 Dockerfile を提供。ユーザーが自分でビルドする形式。

## Technical Context

**Base image**: debian:bookworm-slim
**Tools**: psm (GitHub Releases), SOPS v3.12.1 (GitHub Releases), AWS CLI v2 (official installer)
**Multi-platform**: linux/amd64, linux/arm64 (TARGETARCH)
**All versions pinned as build ARG**

## Constitution Check

| #   | Principle        | Pass?                   |
| --- | ---------------- | ----------------------- |
| I   | Simplicity First | Yes — single Dockerfile |
| II  | YAGNI            | Yes — reference only    |
| III | Test-First       | N/A — Dockerfile        |

## Project Structure

```text
docker/
├── Dockerfile       # Reference Dockerfile
└── README.md        # Usage instructions
```
