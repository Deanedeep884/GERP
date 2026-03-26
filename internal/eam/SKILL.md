---
name: "internal/eam"
description: "Performs enterprise asset and maintenance tracking for eam subsystem. Use when doing cross-domain infrastructure tracking, or when the user mentions warehouses, delivery trucks, or physical assets."
compatibility: ["spanner:latest", "go:1.21"]
metadata:
  gerp-domain: "eam"
  gerp-status: "DRAFT"
  gerp-reviewer: "TODO"
  gerp-id: "skill-eam-01"
---

# Enterprise Asset Management (internal/eam)

This module represents the physical footprint of the GERP matrix: warehouses, delivery trucks, and IoT infrastructure.

## Constraints
- **Zero Foreign Keys:** `FinanceAssetID` points to `internal/finance`, and `TechnicianID` points to `internal/hcm`. Both use implicit UUIDs managed by Temporal.
- **Physical Hierarchy:** MaintenanceLogs are strictly interleaved under Assets in Cloud Spanner to ensure ACID snapshot reads of an asset's total operational health.
