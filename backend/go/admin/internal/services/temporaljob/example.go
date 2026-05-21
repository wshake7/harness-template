package temporaljob

import "go.temporal.io/sdk/workflow"

const PrintCountWorkflowName = "PrintCountWorkflow"

type PrintCountInput struct {
	Count int `json:"count"`
}

type PrintCountResult struct {
	Count int `json:"count"`
}

func PrintCountWorkflow(ctx workflow.Context, input PrintCountInput) (*PrintCountResult, error) {
	workflow.GetLogger(ctx).Info("print count", "count", input.Count)
	return &PrintCountResult{Count: input.Count}, nil
}
