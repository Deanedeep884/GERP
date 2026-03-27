package pipeline_test

import (
	"errors"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"

	"gerp/internal/coams"
	"gerp/internal/pipeline"
)

type PublishSagaTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

func (s *PublishSagaTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

// Test 1: The Integrity Rejection (Simulating a Broken Link)
func (s *PublishSagaTestSuite) TestPublishSaga_RollbackOnBrokenLink() {
	channelID := "engineering"
	rawMarkdown := []byte("# Arch \n [See old spec](doc:invalid-uuid-404)")
	authorID := "agent-007"

	// 1. AST Parser succeeds
	s.env.OnActivity(pipeline.ExtractMarkdownActivity, mock.Anything, channelID, rawMarkdown, authorID).
		Return(coams.ParseResult{Edges: []coams.Edge{{IsExternal: false}}}, nil)

	// 2. Verifier mathematically rejects the link
	s.env.OnActivity(pipeline.VerifyGraphActivity, mock.Anything, channelID, mock.Anything).
		Return(errors.New("ErrReferentialIntegrity: target doc:invalid-uuid-404 not found in partition"))

	// 3. The Compensating Action MUST be called to alert the Agent/CLI (assuming it's added)
	// s.env.OnActivity(pipeline.AlertAgentCompensationActivity, ...).Return(nil)

	// Execute the Saga
	s.env.ExecuteWorkflow(pipeline.CoamsPublishSaga, channelID, rawMarkdown, authorID)

	// Assertions
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.Contains(err.Error(), "ErrReferentialIntegrity")
	
	// CRITICAL: Ensure the Vector Upsert was NEVER called
	s.env.AssertNotCalled(s.T(), "VectorizeChunksActivity")
	s.env.AssertNotCalled(s.T(), "PersistCoamsStorageActivity")
}
