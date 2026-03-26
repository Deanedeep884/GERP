---
name: "cmd/mcp"
description: "Performs autonomous ai tool exposure for mcp subsystem. Use when doing native json-rpc bindings, or when the user mentions model context protocol, agentic environments, or stdio tool servers."
compatibility: ["mcp:latest", "go:1.21"]
metadata:
  gerp-domain: "mcp"
  gerp-status: "DRAFT"
  gerp-reviewer: "TODO"
  gerp-id: "skill-mcp-01"
---

# Model Context Protocol Server (cmd/mcp & internal/mcp)

This module builds the AI Brain Interface for GERP. It exposes the FAANG-grade architecture via the official Model Context Protocol (MCP) using JSON-RPC over STDIO. This allows AI Swarms (like Claude, Gemini, or custom orchestrators) to securely interact with the physical databases and Temporal Sagas autonomously.

## Constraints
- **STDIO Transport:** All data is structured over standard input/output streams natively.
- **Config Bound:** The server strictly utilizes `cli.InitConfig()` to guarantee the AI agent is subject to the identical Spanner and Temporal environments as human sysadmins.
