---
name: "internal/lms"
description: "Performs corporate training and certification tracking for lms subsystem. Use when doing cross-domain compliance checks, or when the user mentions courses, enrollments, or legal certifications."
compatibility: ["spanner:latest", "go:1.21"]
metadata:
  gerp-domain: "lms"
  gerp-status: "DRAFT"
  gerp-reviewer: "TODO"
  gerp-id: "skill-lms-01"
---

# Learning Management System (internal/lms)

This module bridges Human Capital and Legal compliance by ensuring employees are trained, certified, and legally compliant to operate physical assets.

## Constraints
- **Zero Foreign Keys:** `EmployeeID` points to `internal/hcm`. Managed strictly via UUID Golden Threads.
- **Physical Hierarchy:** Enrollments are strictly interleaved under Courses in Cloud Spanner to ensure Snapshot Consistency.
