package hcm

import (
	"time"

	"github.com/google/uuid"
)

// Employee represents the core human capital entity physically stored in the HCM domain.
type Employee struct {
	ID        uuid.UUID `json:"id" spanner:"ID"`
	FirstName string    `json:"first_name" spanner:"FirstName"`
	LastName  string    `json:"last_name" spanner:"LastName"`
	Role      string    `json:"role" spanner:"Role"`
	Email     string    `json:"email" spanner:"Email"` // Synchronization link to Google Workspace
	IsActive  bool      `json:"is_active" spanner:"IsActive"`
	CreatedAt time.Time `json:"created_at" spanner:"CreatedAt"`
	UpdatedAt time.Time `json:"updated_at" spanner:"UpdatedAt"`
}

// PayrollRun acts as the payroll snapshot, heavily interleaved under an Employee.
type PayrollRun struct {
	ID           uuid.UUID `json:"id" spanner:"ID"`
	EmployeeID   uuid.UUID `json:"employee_id" spanner:"EmployeeID"`         // Parent Key
	GrossPay     int64     `json:"gross_pay" spanner:"GrossPay"`             // Minor units (cents)
	NetPay       int64     `json:"net_pay" spanner:"NetPay"`                 // Minor units (cents), after deductions
	FinanceRefID uuid.UUID `json:"finance_ref_id" spanner:"FinanceRefID"`   // The Golden Thread linking this run to the finance ledger lock
	PayPeriodEnd time.Time `json:"pay_period_end" spanner:"PayPeriodEnd"`
	CreatedAt    time.Time `json:"created_at" spanner:"CreatedAt"`
}
