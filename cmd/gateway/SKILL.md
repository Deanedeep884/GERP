---
name: "cmd/gateway"
description: "Performs api exposure and graph stitching for gateway subsystem. Use when doing cross-domain graphql resolving, or when the user mentions the bff, frontend queries, or the universal api."
compatibility: ["gqlgen:latest", "go:1.21"]
metadata:
  gerp-domain: "gateway"
  gerp-status: "DRAFT"
  gerp-reviewer: "TODO"
  gerp-id: "skill-gateway-01"
---

# GraphQL Gateway (cmd/gateway & internal/transport/graphql)

This module acts as the Backend-For-Frontend (BFF). It receives optimized GraphQL queries from the web clients and fans out requests across the isolated micro-domains, stitching the Golden Threads (UUIDs) back together in memory.

## Constraints
- **Zero Raw SQL:** The Gateway must never touch Cloud Spanner directly. It can only call the domain `Service` interfaces.
- **Dataloading:** Resolvers must batch their physical reads via N+1 safe graph loaders when traversing down the Golden Thread.
