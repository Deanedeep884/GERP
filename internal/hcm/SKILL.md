---
name: "internal/hcm"
description: "Performs human capital and payroll management for hcm subsystem. Use when doing cross-domain employee workflows, or when the user mentions personnel, payroll runs, or organizational changes."
compatibility: ["spanner:latest", "go:1.21"]
metadata:
  gerp-domain: "hcm"
  gerp-status: "DRAFT"
  gerp-reviewer: "TODO"
  gerp-id: "skill-hcm-01"
---

# Human Capital Management (internal/hcm)

This module handles the core organizational hierarchy, employee tracking, and payroll histories. 

## Constraints
- **Zero Foreign Keys:** Never hard-link the `FinanceRefID` to the finance ledger using SQL. Temporal guarantees the integrity.
- **Amounts:** Payroll fields (`GrossPay`, `NetPay`) must be expressed in `int64` minor units (cents).
