GERP(1)                                    GERP Matrix Man Page                                    GERP(1)

NAME
       gerp - Google ERP Master Control Plane & Native Operator CLI

SYNOPSIS
       make [COMMAND]
       gerp [COMMAND] [ARGS...]

DESCRIPTION
       The GERP control plane utilizes Makefile targets to orchestrate a distributed, 7-domain ERP 
       matrix structurally bound by Temporal Sagas and isolated Spanner data layers.
       
       The compiled `gerp` native CLI allows system administrators and AI Autonomous Agents to
       mechanically route direct queries against the Spanner telemetry or asynchronously fire
       complex Temporal execution traces.

ARCHITECTURE
       The physical runtimes consist of:
       - Cloud Spanner Emulator (Local instance binding all 7 domains natively).
       - Temporal Server (Local instance managing `GERP_GLOBAL_QUEUE`).
       - Go GraphQL Gateway (Port 8080: Binds the BFF execution).
       - Go Temporal Worker (Background: Listens for Saga dispatches).
       - Go MCP Server (STDIO: Links external Model Context Protocol systems).

CONFIGURATION
       The `gerp` CLI relies natively on a `.gerp.yaml` configuration profile mapped to the 
       user's local `$HOME` directory or explicitly bound via the following environment variables:
       
       GERP_GRAPHQL_ENDPOINT      Target API Gateway (Default: http://localhost:8080/query)
       GERP_TEMPORAL_HOST         Temporal Control Queue (Default: localhost:7233)
       GERP_SPANNER_DB            Sovereign Storage Domain Bounds

MAKE COMMANDS
       make up
           Boot the local Docker infrastructure encompassing Spanner, Temporal, and associated matrix layers.

       make init-db
           Create the Spanner emulator instance and execute the physical DDL files for all 7 downstream domains, permanently interleaving the schema state.

       make run-gateway
           Boot the GraphQL Backend-For-Frontend (BFF). Exposes the GraphQL Playground and routing execution path on http://localhost:8080.

       make run-worker
           Boot the Temporal Saga orchestrator. Listens dynamically to the Temporal Server queue and executes multi-domain asynchronous mutations.

       make generate
           Regenerate the physical Go bindings for the unified GraphQL API.

GERP NATIVE COMMANDS
       gerp status
           Pings the GERP matrix execution limits and strictly formats the current operational bounds for the sysadmin.

       gerp orders create
           Triggers the Global Fulfillment Saga via GraphQL. Safely executing cross-domain mutations across SCM and Finance.

       gerp audit view [target_record_id]
           Bypasses the API Gateway to directly query Spanner ComplianceAudits locally. Deep inspection specifically formulated for Legal operations.

       gerp add coams
           Injects the Content Operating & Management System into the local matrix. Seeds the QuanuX central Knowledge Vector and provisions pgvector AlloyDB shards dynamically.

       gerp coams sync [directory]
           Executes the Publish Saga Lifecycle across a directory of raw Markdown. Validates mathematical link integrity via the Agent-Index and broadcasts the dynamically generated GraphQL schema.

       gerp coams gen-man
           Autonomously generates extensive UNIX manual pages for COAMS sub-commands, seamlessly mapping back into the QuanuX Knowledge Vector bounds.

SEE ALSO
       docker-compose(1), temporal(1), gcloud-spanner(1)
