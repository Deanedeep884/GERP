package lms

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

// Service defines the boundary for educational compliance and certification management.
type Service interface {
	GetCourseWithEnrollments(ctx context.Context, courseID uuid.UUID) (*Course, []*Enrollment, error)
}

type lmsService struct {
	client *spanner.Client
}

// NewService provisions the LMS service with the dedicated Spanner client.
func NewService(client *spanner.Client) Service {
	return &lmsService{client: client}
}

// GetCourseWithEnrollments guarantees an ACID snapshot of a course and its active student load.
func (s *lmsService) GetCourseWithEnrollments(ctx context.Context, courseID uuid.UUID) (*Course, []*Enrollment, error) {
	txn := s.client.ReadOnlyTransaction()
	defer txn.Close()

	// 1. Snapshot read of the Parent Course Record
	row, err := txn.ReadRow(ctx, "Courses", spanner.Key{courseID.String()}, []string{
		"ID", "Title", "IsMandatory", "ValidForDays", "CreatedAt", "UpdatedAt",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read course log: %w", err)
	}

	var course Course
	if err := row.ToStruct(&course); err != nil {
		return nil, nil, fmt.Errorf("failed to decode course payload: %w", err)
	}

	// 2. Consistent snapshot query for interleaved Enrollments
	stmt := spanner.Statement{
		SQL: `SELECT CourseID, ID, EmployeeID, Status, Score, CompletedAt, CreatedAt, UpdatedAt 
              FROM Enrollments 
              WHERE CourseID = @course_id`,
		Params: map[string]interface{}{
			"course_id": courseID.String(),
		},
	}
	
	iter := txn.Query(ctx, stmt)
	defer iter.Stop()

	var enrollments []*Enrollment
	for {
		eRow, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("enrollments queries failed: %w", err)
		}
		var enr Enrollment
		if err := eRow.ToStruct(&enr); err != nil {
			return nil, nil, fmt.Errorf("failed to decode enrollment struct: %w", err)
		}
		enrollments = append(enrollments, &enr)
	}

	return &course, enrollments, nil
}
