package tools

import (
	"admin/internal/services/httpc"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type PrometheusAlert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	State       string            `json:"state"`
	ActiveAt    string            `json:"activeAt"`
	Value       string            `json:"value"`
}

type PrometheusAlertsResult struct {
	Status string `json:"status"`
	Data   struct {
		Alerts []PrometheusAlert `json:"alerts"`
	} `json:"data"`
	Error     string `json:"error,omitempty"`
	ErrorType string `json:"errorType,omitempty"`
}

type SimplifiedAlert struct {
	AlertName   string `json:"alert_name" jsonschema:"description=Alert name from Prometheus labels.alertname"`
	Description string `json:"description" jsonschema:"description=Alert description from annotations.description"`
	State       string `json:"state" jsonschema:"description=Alert state: firing or pending"`
	ActiveAt    string `json:"active_at" jsonschema:"description=Alert activation time in RFC3339 format"`
	Duration    string `json:"duration" jsonschema:"description=Duration since alert activated"`
}

type PrometheusAlertsOutput struct {
	Success bool              `json:"success" jsonschema:"description=Whether the query was successful"`
	Alerts  []SimplifiedAlert `json:"alerts,omitempty" jsonschema:"description=Active alerts with deduplicated alert names"`
	Message string            `json:"message,omitempty" jsonschema:"description=Status message"`
	Error   string            `json:"error,omitempty" jsonschema:"description=Error message if query failed"`
}

func queryPrometheusAlerts(baseURL string) (PrometheusAlertsResult, error) {
	var result PrometheusAlertsResult

	resp, err := httpc.Client.R().Get(fmt.Sprintf("%s/api/v1/alerts", baseURL))
	if err != nil {
		return result, fmt.Errorf("failed to query Prometheus alerts: %v", err)
	}
	if resp.IsError() {
		return result, fmt.Errorf("prometheus returned status %d: %s", resp.StatusCode(), resp.String())
	}
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return result, fmt.Errorf("failed to parse response: %v", err)
	}
	return result, nil
}

func calculateDuration(activeAtStr string) string {
	activeAt, err := time.Parse(time.RFC3339Nano, activeAtStr)
	if err != nil {
		return "unknown"
	}
	duration := time.Since(activeAt)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh%dm%ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

func NewPrometheusAlertsQueryTool() tool.InvokableTool {
	t, err := utils.InferOptionableTool(
		"query_prometheus_alerts",
		"Query active alerts from Prometheus alerting system. This tool retrieves all currently active/firing alerts including their labels, annotations, state, and values. Use this tool when you need to check what alerts are currently firing or investigate alert conditions.",
		func(ctx context.Context, input *struct{}, opts ...tool.Option) (output string, err error) {
			result, err := queryPrometheusAlerts("http://127.0.0.1:9090")
			if err != nil {
				out := PrometheusAlertsOutput{
					Success: false,
					Error:   err.Error(),
					Message: "Failed to query Prometheus alerts",
				}
				b, _ := json.MarshalIndent(out, "", "  ")
				return string(b), err
			}

			seen := make(map[string]bool)
			simplified := make([]SimplifiedAlert, 0)
			for _, alert := range result.Data.Alerts {
				name := alert.Labels["alertname"]
				if seen[name] {
					continue
				}
				seen[name] = true
				simplified = append(simplified, SimplifiedAlert{
					AlertName:   name,
					Description: alert.Annotations["description"],
					State:       alert.State,
					ActiveAt:    alert.ActiveAt,
					Duration:    calculateDuration(alert.ActiveAt),
				})
			}

			out := PrometheusAlertsOutput{
				Success: true,
				Alerts:  simplified,
				Message: fmt.Sprintf("Successfully retrieved %d active alerts", len(simplified)),
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
