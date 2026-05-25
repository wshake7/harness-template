package knowledgedocument

import (
	"admin/internal/appsvc"
	"admin/internal/config"
	"admin/internal/services/objectstore"
	"admin/internal/services/orm/query"
	"context"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const IndexWorkflowName = "KnowledgeDocumentIndexWorkflow"

type IndexInput struct {
	DocumentID uint64 `json:"documentID"`
}

func IndexWorkflow(ctx workflow.Context, input IndexInput) error {
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 2,
			BackoffCoefficient: 2,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	return workflow.ExecuteActivity(ctx, IndexActivity, input).Get(ctx, nil)
}

func IndexActivity(ctx context.Context, input IndexInput) error {
	indexer := appsvc.NewKnowledgeDocumentIndexer(query.Q, config.Current(), objectstore.Current())
	return indexer.IndexDocument(ctx, input.DocumentID)
}
