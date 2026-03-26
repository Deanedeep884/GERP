---
name: "cmd/worker"
description: "Performs distributed sequence orchestration for worker subsystem. Use when doing cross-domain worker queues, or when the user mentions temporal engine, background jobs, or saga runtimes."
compatibility: ["temporal:latest", "go:1.21"]
metadata:
  gerp-domain: "worker"
  gerp-status: "DRAFT"
  gerp-reviewer: "TODO"
  gerp-id: "skill-worker-01"
---

# Temporal Orchestrator Worker (cmd/worker)

This module is the physical execution engine that listens to the Temporal cluster queue (`GERP_GLOBAL_QUEUE`) and securely triggers cross-domain operations in Cloud Spanner. 

## Constraints
- **Compute Isolation:** The worker executes Sagas mathematically by orchestrating decoupled `Activity` functions (`AllocateInventoryActivity`, `ChargeLedgerActivity`).
