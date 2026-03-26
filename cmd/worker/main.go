package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/spanner"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"gerp/internal/finance"
	"gerp/internal/pipeline"
	"gerp/internal/scm"
)

func main() {
	log.Println("⚙️ Booting GERP Temporal Worker Engine...")

	ctx := context.Background()

	// Step 1: Initialize Cloud Spanner
	spannerDB := os.Getenv("SPANNER_DATABASE")
	if spannerDB == "" {
		spannerDB = "projects/gerp-local-dev/instances/gerp-instance/databases/gerp-db"
	}

	spannerClient, err := spanner.NewClient(ctx, spannerDB)
	if err != nil {
		log.Fatalf("🚨 Failed to connect to Spanner: %v", err)
	}
	defer spannerClient.Close()
	log.Println("✅ Bound to Cloud Spanner Emulator/Instance.")

	// Step 2: Initialize Temporal Client
	temporalTarget := os.Getenv("TEMPORAL_HOST_PORT")
	if temporalTarget == "" {
		temporalTarget = "localhost:7233"
	}

	temporalClient, err := client.Dial(client.Options{
		HostPort: temporalTarget,
	})
	if err != nil {
		log.Fatalf("🚨 Unable to connect to Temporal cluster: %v", err)
	}
	defer temporalClient.Close()
	log.Println("✅ Bound to Temporal Orchestration Queue.")

	// Step 3: Instantiate pure domain services
	finSvc := finance.NewService(spannerClient)
	scmSvc := scm.NewService(spannerClient)

	// Step 4: Instantiate Temporal bound activities
	finActivities := finance.NewActivities(finSvc)
	scmActivities := scm.NewActivities(scmSvc)

	// Step 5: Configure Worker Execution Queue
	w := worker.New(temporalClient, "GERP_GLOBAL_QUEUE", worker.Options{})

	// Step 6: Register the centralized Fulfillment Saga Orchestration
	w.RegisterWorkflow(pipeline.GlobalFulfillmentSaga)

	// Step 7: Register bounded Domain Activities executing physical Spanner writes
	w.RegisterActivity(scmActivities.AllocateInventoryActivity)
	w.RegisterActivity(scmActivities.ReverseInventoryActivity)
	w.RegisterActivity(finActivities.ChargeLedgerActivity)

	// Step 8: Open the port to execute the matrix orchestrations securely
	log.Println("🚀 GERP Worker successfully initialized! Listening for Saga Dispatches on GERP_GLOBAL_QUEUE...")
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("🚨 Worker Execution Halted: %v", err)
	}
}
