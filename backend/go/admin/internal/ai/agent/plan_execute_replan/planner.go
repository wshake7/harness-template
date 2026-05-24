package plan_execute_replan

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/planexecute"
	"github.com/cloudwego/eino/components/model"
)

func NewPlanner(ctx context.Context, planModel model.ToolCallingChatModel) (adk.Agent, error) {
	return planexecute.NewPlanner(ctx, &planexecute.PlannerConfig{
		ToolCallingChatModel: planModel,
	})
}
