---
name: "internal/finance"
description: "Performs ledger transactions and auditing for finance subsystem. Use when doing cross-domain Spanner transactions, or when the user mentions financial ledger, accounting, or journaling."
compatibility: ["spanner:latest", "go:1.21"]
metadata:
  gerp-domain: "finance"
  gerp-status: "DRAFT"
  gerp-reviewer: "TODO"
  gerp-id: "skill-finance-01"
---

# Finance Domain (internal/finance)

This module manages the core financial ledger, journal entries, and account balances natively in Cloud Spanner. It strictly isolates its database transactions from other domains like HCM or MDM.

## Constraints
- **Zero Foreign Keys:** Never hard-link to HCM, MDM, or SCM tables. Use `uuid.UUID`.
- **Stateless Workers:** All cross-domain operations must be wrapped in Temporal sagas.
- **Amounts:** Always use `int64` minor units (cents).
