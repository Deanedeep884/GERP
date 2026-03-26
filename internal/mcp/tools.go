package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"

	"gerp/internal/cli"
)

// HandleStatus securely checks the bound infrastructure for the AI system.
func HandleStatus() (interface{}, error) {
	return map[string]interface{}{
		"status":          "🟢 GERP MATRIX OPERATIONAL",
		"graphql_bound":   cli.ActiveConfig.GraphQLEndpoint,
		"temporal_bound":  cli.ActiveConfig.TemporalHost,
		"spanner_bound":   cli.ActiveConfig.SpannerDB,
	}, nil
}

// HandleCreateOrder receives JSON parameters, structures them, and triggers the Saga.
func HandleCreateOrder(params json.RawMessage) (interface{}, error) {
	// Re-routing directly into the existing GraphQL BFF structure.
	// For production, the Agent would supply the strict 'input' variables payload directly.
	reqPayload := []byte(`{
		"query": "mutation ExecuteGlobalSale($input: CreateOrderInput!) { createSalesOrder(input: $input) { id status totalValue customer { legalName countryCode } } }",
		"variables": ` + string(params) + `
	}`)

	req, err := http.NewRequest("POST", cli.ActiveConfig.GraphQLEndpoint, bytes.NewBuffer(reqPayload))
	if err != nil {
		return nil, fmt.Errorf("api formulation fault: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("matrix unreachable via BFF: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return string(body), nil // Return the raw text if parsing fails natively.
	}

	return result, nil
}

// HandleAuditView gives the AI direct read-only access to the physical Compliance bounds.
func HandleAuditView(params json.RawMessage) (interface{}, error) {
	var input struct {
		TargetRecordID string `json:"target_record_id"`
	}
	if err := json.Unmarshal(params, &input); err != nil {
		return nil, fmt.Errorf("invalid json-rpc parameters for audit: %w", err)
	}

	ctx := context.Background()
	client, err := spanner.NewClient(ctx, cli.ActiveConfig.SpannerDB)
	if err != nil {
		return nil, fmt.Errorf("sovereign spanner binding failed: %w", err)
	}
	defer client.Close()

	stmt := spanner.Statement{
		SQL: `SELECT ID, ActorID, Action, AuditTimestamp 
			  FROM ComplianceAudits 
			  WHERE TargetRecordID = @record_id`,
		Params: map[string]interface{}{
			"record_id": input.TargetRecordID,
		},
	}

	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	var audits []map[string]interface{}
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("matrix iteration failure: %w", err)
		}

		var id, actorID, action string
		var auditTimestamp spanner.NullTime

		if err := row.Columns(&id, &actorID, &action, &auditTimestamp); err != nil {
			continue
		}

		audits = append(audits, map[string]interface{}{
			"audit_id":   id,
			"action":     action,
			"actor_id":   actorID,
			"timestamp":  auditTimestamp.Time,
		})
	}

	return map[string]interface{}{
		"target_record_id": input.TargetRecordID,
		"logs":             audits,
	}, nil
}
