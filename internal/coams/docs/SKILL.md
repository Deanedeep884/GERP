---
name: COAMS
description: Content Operating and Management System (Markdown-Native CMS)
---

# COAMS (GERP Module)

## Overview
COAMS is an AI-First, ignorant content engine injected into the GERP Matrix. It does not speak JSON; it speaks native Markdown. When you (the Agent) interact with COAMS, you operate physically isolated Vector shards utilizing `pgvector` stored within AlloyDB.

## Architectural Constraints (DO NOT VIOLATE)
1. **Zero Broken Links:** COAMS mathematically enforces the Agent-Index. Every `doc:uuid` link you author in Markdown MUST exist within the isolated `channel_id` bounds OR the commit will violently reject the transaction.
2. **Ignorant Domain:** The pure Go `internal/coams` layer knows nothing about HTTP or the existence of AI. It operates purely on interfaces. 
3. **Partition Pruning:** You cannot search outside your immediate IAM token boundary. The AlloyDB SQL planner will literally drop vectors outside your `channel_id`.
4. **Verbosity:** All documents carry heavy metadata matching the IAM footprint.

## Executing the Publish Saga Lifecycle
To persist arbitrary markdown documents into COAMS:
```bash
gerp coams sync ./docs
```

This triggers the orchestration layers to extract ASTs, run Graph Verification (`verifier.go`), chunk texts, invoke Vertex AI, save vectors, and automatically regenerate the GraphQL Broadcast schema dynamically.

## Man Pages
For deeper integration, execute `gerp coams gen-man` to render system operator manuals natively.
