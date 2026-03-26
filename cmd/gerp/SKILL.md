---
name: "cmd/gerp"
description: "Performs native command line terminal operations for gerp subsystem. Use when doing matrix cli executions, or when the user mentions cobra, terminal operators, viper configs, or the gerp binary."
compatibility: ["cobra:latest", "viper:latest", "go:1.21"]
metadata:
  gerp-domain: "cli"
  gerp-status: "DRAFT"
  gerp-reviewer: "TODO"
  gerp-id: "skill-cli-01"
---

# GERP CLI Operator (cmd/gerp & internal/cli)

This module builds the definitive native operator terminal `gerp`. It bypasses the BFF web boundary, giving system administrators direct, low-latency control and query capability into the core Spanner domains and Temporal execution queues using `spf13/cobra` (routing) and `spf13/viper` (configuration).

## Constraints
- **Config Bound:** The binary dynamically binds to `.gerp.yaml` in the user's home directory across different shell environments.
