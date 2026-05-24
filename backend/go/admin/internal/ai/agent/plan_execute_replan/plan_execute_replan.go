package plan_execute_replan

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/planexecute"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
)

func BuildPlanAgent(ctx context.Context, planModel, execModel, replanModel model.ToolCallingChatModel, tools []tool.BaseTool, query string) (string, []string, error) {
	planAgent, err := NewPlanner(ctx, planModel)
	if err != nil {
		return "", []string{}, err
	}
	executeAgent, err := NewExecutor(ctx, execModel, tools)
	if err != nil {
		return "", []string{}, err
	}
	replanAgent, err := NewRePlanAgent(ctx, replanModel)
	if err != nil {
		return "", []string{}, err
	}
	planExecuteAgent, err := planexecute.New(ctx, &planexecute.Config{
		Planner:       planAgent,
		Executor:      executeAgent,
		Replanner:     replanAgent,
		MaxIterations: 20,
	})
	if err != nil {
		return "", []string{}, fmt.Errorf("build PlanExecuteAgent Error: %v", err)
	}
	r := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: planExecuteAgent,
	})
	iter := r.Query(ctx, query)
	var lastMessage adk.Message
	var detail []string
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		fmt.Println("------------- Event -------------")
		fmt.Printf("%+v\n", event)
		if event.Output != nil {
			lastMessage, _, err = adk.GetMessage(event)
			detail = append(detail, lastMessage.String())
		}
	}
	if lastMessage == nil {
		return "", []string{}, fmt.Errorf("get lastMessage Error")
	}
	return lastMessage.Content, detail, nil
}
