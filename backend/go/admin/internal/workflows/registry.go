package workflows

import (
	"admin/internal/services/temporaljob"
	exampleworkflow "admin/internal/workflows/example"

	"go.temporal.io/sdk/workflow"
)

func WorkflowTypeOptions() map[string]string {
	return map[string]string{
		exampleworkflow.PrintCountWorkflowName: exampleworkflow.PrintCountWorkflowName,
	}
}

func RegisterWorker(w temporaljob.WorkerRegistry) {
	w.RegisterWorkflowWithOptions(exampleworkflow.PrintCountWorkflow, workflow.RegisterOptions{Name: exampleworkflow.PrintCountWorkflowName})
}
