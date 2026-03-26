package hcm

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

// Service defines the boundary for human capital and payroll operations.
type Service interface {
	GetEmployeeWithPayroll(ctx context.Context, employeeID uuid.UUID) (*Employee, []*PayrollRun, error)
	InsertEmployee(ctx context.Context, emp *Employee) error
	InsertPayrollRun(ctx context.Context, run *PayrollRun) error
}

type hcmService struct {
	client *spanner.Client
}

// NewService provisions the HCM service with the dedicated Spanner client.
func NewService(client *spanner.Client) Service {
	return &hcmService{client: client}
}

// GetEmployeeWithPayroll returns the core employee identity and their interleaved payroll history.
// It leverages a ReadOnlyTransaction snapshot to perfectly isolate the multi-table reads in high-TPS environments.
func (s *hcmService) GetEmployeeWithPayroll(ctx context.Context, employeeID uuid.UUID) (*Employee, []*PayrollRun, error) {
	// Snapshot Isolation: Ensure new payroll commits mid-read never cause dirty/incomplete states
	txn := s.client.ReadOnlyTransaction()
	defer txn.Close()

	// 1. Snapshot read of the Parent Employee Record
	row, err := txn.ReadRow(ctx, "Employees", spanner.Key{employeeID.String()}, []string{
		"ID", "FirstName", "LastName", "Role", "Email", "IsActive", "CreatedAt", "UpdatedAt",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read employee: %w", err)
	}

	var emp Employee
	if err := row.ToStruct(&emp); err != nil {
		return nil, nil, fmt.Errorf("failed to decode employee payload: %w", err)
	}

	// 2. Consistent snapshot query for all dependent Payroll data
	stmt := spanner.Statement{
		SQL: `SELECT EmployeeID, ID, GrossPay, NetPay, FinanceRefID, PayPeriodEnd, CreatedAt 
              FROM PayrollRuns 
              WHERE EmployeeID = @employee_id
              ORDER BY PayPeriodEnd DESC`,
		Params: map[string]interface{}{
			"employee_id": employeeID.String(),
		},
	}
	
	iter := txn.Query(ctx, stmt)
	defer iter.Stop()

	var payrolls []*PayrollRun
	for {
		prRow, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("payroll run query traversal failed: %w", err)
		}
		var run PayrollRun
		if err := prRow.ToStruct(&run); err != nil {
			return nil, nil, fmt.Errorf("failed to decode payroll struct: %w", err)
		}
		payrolls = append(payrolls, &run)
	}

	return &emp, payrolls, nil
}

// InsertEmployee commits a new human capital record.
func (s *hcmService) InsertEmployee(ctx context.Context, emp *Employee) error {
	mut, err := spanner.InsertStruct("Employees", emp)
	if err != nil {
		return err
	}
	_, err = s.client.Apply(ctx, []*spanner.Mutation{mut})
	return err
}

// InsertPayrollRun persists a payroll event. The FinanceRefID acts as the Golden Thread.
func (s *hcmService) InsertPayrollRun(ctx context.Context, run *PayrollRun) error {
	mut, err := spanner.InsertStruct("PayrollRuns", run)
	if err != nil {
		return err
	}
	_, err = s.client.Apply(ctx, []*spanner.Mutation{mut})
	return err
}
