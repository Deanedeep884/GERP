---
name: "internal/mdm"
description: "Performs universal record translation for mdm subsystem. Use when doing cross-domain identity mapping, or when the user mentions the golden record, master data, or vendor translation."
compatibility: ["spanner:latest", "go:1.21"]
metadata:
  gerp-domain: "mdm"
  gerp-status: "DRAFT"
  gerp-reviewer: "TODO"
  gerp-id: "skill-mdm-01"
---

# Master Data Management (internal/mdm)

This module is the Universal Translator of the matrix. It maps isolated domain entities (e.g., a "Customer" in Revenue and a "Vendor" in Supply Chain) back to a single Golden Record representing a real-world company.

## Constraints
- **Zero Foreign Keys:** `LocalID` acts as a pure pointer to another domain's physical row. 
- **Physical Hierarchy:** `EntityMapping` is interleaved under `GlobalEntity` in Cloud Spanner to ensure instant mapping resolution without SQL joins.
