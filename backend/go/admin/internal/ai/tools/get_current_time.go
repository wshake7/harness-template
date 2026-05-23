package tools

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type GetCurrentTimeInput struct{}

type GetCurrentTimeOutput struct {
	Success      bool   `json:"success" jsonschema:"description=Indicates whether the time retrieval was successful"`
	Seconds      int64  `json:"seconds" jsonschema:"description=Current Unix timestamp in seconds since epoch"`
	Milliseconds int64  `json:"milliseconds" jsonschema:"description=Current Unix timestamp in milliseconds since epoch"`
	Microseconds int64  `json:"microseconds" jsonschema:"description=Current Unix timestamp in microseconds since epoch"`
	Timestamp    string `json:"timestamp" jsonschema:"description=Human-readable timestamp in format YYYY-MM-DD HH:MM:SS.microseconds"`
	Message      string `json:"message" jsonschema:"description=Status message describing the operation result"`
}

func NewGetCurrentTimeTool() tool.InvokableTool {
	t, err := utils.InferOptionableTool(
		"get_current_time",
		"Get current system time in multiple formats. Returns the current time in seconds (Unix timestamp), milliseconds, and microseconds. Use this tool when you need to retrieve current system time for logging, timing operations, or timestamping events.",
		func(ctx context.Context, input *GetCurrentTimeInput, opts ...tool.Option) (output string, err error) {
			now := time.Now()

			out := GetCurrentTimeOutput{
				Success:      true,
				Seconds:      now.Unix(),
				Milliseconds: now.UnixMilli(),
				Microseconds: now.UnixMicro(),
				Timestamp:    now.Format("2006-01-02 15:04:05.000000"),
				Message:      "Current time retrieved successfully",
			}

			b, err := json.MarshalIndent(out, "", "  ")
			if err != nil {
				return "", err
			}
			return string(b), nil
		})
	if err != nil {
		panic(err)
	}
	return t
}
