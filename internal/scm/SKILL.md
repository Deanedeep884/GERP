---
name: "internal/scm"
description: "Performs supply chain and inventory management for scm subsystem. Use when doing cross-domain stock allocations, or when the user mentions physical products, lots, or warehouse counts."
compatibility: ["spanner:latest", "go:1.21"]
metadata:
  gerp-domain: "scm"
  gerp-status: "DRAFT"
  gerp-reviewer: "TODO"
  gerp-id: "skill-scm-01"
---

# Supply Chain & Manufacturing (internal/scm)

This module commands the definition of physical goods and controls absolute inventory truth across the global warehouse footprint.

## Constraints
- **Zero Foreign Keys:** `WarehouseID` natively points to the Enterprise Asset Management (`eam`) domain entirely via the Go-supervised Golden Thread.
- **Strict Mutability:** Inventory Lots cannot drop below 0 quantities at the database transaction layer.
