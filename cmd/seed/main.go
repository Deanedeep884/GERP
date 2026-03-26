// cmd/seed/main.go
package main

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
)

// Deterministic UUIDs so we can hardcode our GraphQL test mutations!
var (
	GlobalCorpID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	CustomerID   = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	EmployeeID   = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	WarehouseID  = uuid.MustParse("44444444-4444-4444-4444-444444444444")
	ProductID    = uuid.MustParse("55555555-5555-5555-5555-555555555555")
	LotID        = uuid.MustParse("66666666-6666-6666-6666-666666666666")
	AccountAR    = uuid.MustParse("77777777-7777-7777-7777-777777777777") // Accounts Receivable (Debit)
	AccountRev   = uuid.MustParse("88888888-8888-8888-8888-888888888888") // Product Revenue (Credit)
)

func main() {
	log.Println("🌱 Booting GERP Genesis Seeder...")

	ctx := context.Background()
	spannerDB := os.Getenv("SPANNER_DATABASE")
	if spannerDB == "" {
		spannerDB = "projects/gerp-local-dev/instances/gerp-instance/databases/gerp-db"
	}

	client, err := spanner.NewClient(ctx, spannerDB)
	if err != nil {
		log.Fatalf("🚨 Failed to connect to Spanner Emulator: %v", err)
	}
	defer client.Close()

	now := time.Now().UTC()
	var mutations []*spanner.Mutation

	// 1. MDM & REVENUE: The Customer Golden Thread
	mutations = append(mutations, 
		spanner.InsertMap("GlobalEntities", map[string]interface{}{
			"ID": GlobalCorpID.String(), "LegalName": "Acme Corp", "TaxID": "US-12345", "CountryCode": "US", "CreatedAt": now, "UpdatedAt": now,
		}),
		spanner.InsertMap("EntityMappings", map[string]interface{}{
			"GlobalEntityID": GlobalCorpID.String(), "Domain": "revenue", "LocalID": CustomerID.String(), "CreatedAt": now, "UpdatedAt": now,
		}),
		spanner.InsertMap("Customers", map[string]interface{}{
			"ID": CustomerID.String(), "Name": "Acme Corp (North America)", "MasterDataID": GlobalCorpID.String(), "CreditLimit": int64(50000000), "CreatedAt": now, "UpdatedAt": now,
		}),
	)

	// 2. EAM & SCM: The Physical Infrastructure and Inventory
	mutations = append(mutations,
		spanner.InsertMap("Assets", map[string]interface{}{
			"ID": WarehouseID.String(), "Name": "Tokyo Central Fulfillment", "Type": "WAREHOUSE", "Status": "ONLINE", "CreatedAt": now, "UpdatedAt": now,
		}),
		spanner.InsertMap("Products", map[string]interface{}{
			"ID": ProductID.String(), "SKU": "GERP-SERVER-01", "Name": "Enterprise Server Rack", "IsActive": true, "CreatedAt": now, "UpdatedAt": now,
		}),
		spanner.InsertMap("InventoryLots", map[string]interface{}{
			"ID": LotID.String(), "ProductID": ProductID.String(), "WarehouseID": WarehouseID.String(), "Quantity": int64(150), "CostBasis": int64(100000), "CreatedAt": now, "UpdatedAt": now,
		}),
	)

	// 3. FINANCE: The Ledger Accounts
	mutations = append(mutations,
		spanner.InsertMap("Accounts", map[string]interface{}{
			"ID": AccountAR.String(), "Name": "Accounts Receivable (A/R)", "Type": "ASSET", "CreatedAt": now, "UpdatedAt": now,
		}),
		spanner.InsertMap("Accounts", map[string]interface{}{
			"ID": AccountRev.String(), "Name": "Hardware Revenue", "Type": "REVENUE", "CreatedAt": now, "UpdatedAt": now,
		}),
	)

	// Commit the Genesis Block
	_, err = client.Apply(ctx, mutations)
	if err != nil {
		log.Fatalf("🚨 Failed to seed Spanner: %v", err)
	}

	log.Println("✅ Genesis Seed Complete. The Matrix is populated.")
	log.Println("==================================================")
	log.Println("TEST DATA FOR GRAPHQL PLAYGROUND:")
	log.Printf("Customer ID: %s", CustomerID)
	log.Printf("Lot ID:      %s", LotID)
	log.Printf("Debit A/C:   %s", AccountAR)
	log.Printf("Credit A/C:  %s", AccountRev)
	log.Println("==================================================")
}
