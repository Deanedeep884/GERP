---
name: "internal/pipeline"
description: "Performs cross-domain transaction orchestration for pipeline subsystem. Use when doing workflow sagas, or when the user mentions temporal, rollback compensations, or fulfillment logic."
compatibility: ["temporal:latest", "go:1.21"]
metadata:
  gerp-domain: "pipeline"
  gerp-status: "DRAFT"
  gerp-reviewer: "TODO"
  gerp-id: "skill-pipeline-01"
---

# Temporal Orchestrator (internal/pipeline)

This module acts as the nervous system for the GERP matrix. It safely coordinates complex transaction sagas across the strictly isolated structural domains.

## Constraints
- **Zero Raw SQL:** The pipeline never talks to the Spanner databases directly. It only invokes strict Activity Contracts bound to the domains.
- **Compensating Rollbacks:** Every state mutation must be caught and gracefully reversed using its sibling reversing activity if downstream errors occur.
