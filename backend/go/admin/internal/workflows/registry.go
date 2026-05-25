package workflows

import (
	"admin/internal/services/temporaljob"
	exampleworkflow "admin/internal/workflows/example"
	knowledgedocument "admin/internal/workflows/knowledge_document"

	"go.temporal.io/sdk/workflow"
)

func WorkflowTypeOptions() map[string]string {
	return map[string]string{
		exampleworkflow.PrintCountWorkflowName: exampleworkflow.PrintCountWorkflowName,
		knowledgedocument.IndexWorkflowName:    knowledgedocument.IndexWorkflowName,
	}
}

func RegisterWorker(w temporaljob.WorkerRegistry) {
	w.RegisterWorkflowWithOptions(exampleworkflow.PrintCountWorkflow, workflow.RegisterOptions{Name: exampleworkflow.PrintCountWorkflowName})
	w.RegisterWorkflowWithOptions(knowledgedocument.IndexWorkflow, workflow.RegisterOptions{Name: knowledgedocument.IndexWorkflowName})
	w.RegisterActivity(knowledgedocument.IndexActivity)
}
