package lms

import (
	"time"

	"github.com/google/uuid"
)

// Course represents a training curriculum or compliance requirement.
type Course struct {
	ID           uuid.UUID `json:"id" spanner:"ID"`
	Title        string    `json:"title" spanner:"Title"`
	IsMandatory  bool      `json:"is_mandatory" spanner:"IsMandatory"`
	ValidForDays int       `json:"valid_for_days" spanner:"ValidForDays"`
	CreatedAt    time.Time `json:"created_at" spanner:"CreatedAt"`
	UpdatedAt    time.Time `json:"updated_at" spanner:"UpdatedAt"`
}

// Enrollment tracks a specific employee's progress through an interleaved Course.
type Enrollment struct {
	CourseID    uuid.UUID  `json:"course_id" spanner:"CourseID"`     // Parent Key
	ID          uuid.UUID  `json:"id" spanner:"ID"`
	EmployeeID  uuid.UUID  `json:"employee_id" spanner:"EmployeeID"` // Golden Thread mapping to HCM
	Status      string     `json:"status" spanner:"Status"`          // e.g., "ENROLLED", "IN_PROGRESS", "COMPLETED", "FAILED"
	Score       *int       `json:"score" spanner:"Score"`            // Nullable grading metric
	CompletedAt *time.Time `json:"completed_at" spanner:"CompletedAt"`
	CreatedAt   time.Time  `json:"created_at" spanner:"CreatedAt"`
	UpdatedAt   time.Time  `json:"updated_at" spanner:"UpdatedAt"`
}

// Certification represents the formally issued legal compliance standing for an Employee.
type Certification struct {
	ID         uuid.UUID `json:"id" spanner:"ID"`
	EmployeeID uuid.UUID `json:"employee_id" spanner:"EmployeeID"` // Golden Thread: HCM Entity
	CourseID   uuid.UUID `json:"course_id" spanner:"CourseID"`     // Origin Course mapping
	IssuedAt   time.Time `json:"issued_at" spanner:"IssuedAt"`
	ExpiresAt  time.Time `json:"expires_at" spanner:"ExpiresAt"`
	Revoked    bool      `json:"revoked" spanner:"Revoked"`
	CreatedAt  time.Time `json:"created_at" spanner:"CreatedAt"`
	UpdatedAt  time.Time `json:"updated_at" spanner:"UpdatedAt"`
}
