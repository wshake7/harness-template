package tools

import (
	"admin/internal/services/orm"
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type DBCrudInput struct {
	SQL         string `json:"sql" jsonschema:"description=The SQL query to execute against the database"`
	OperateType string `json:"operate_type" jsonschema:"description=The type of SQL operation: query, insert, update, or delete"`
}

type DBCrudOutput struct {
	Success  bool   `json:"success" jsonschema:"description=Whether the operation was successful"`
	Message  string `json:"message" jsonschema:"description=Status message"`
	RowsAffected int64  `json:"rows_affected,omitempty" jsonschema:"description=Number of rows affected (for insert/update/delete)"`
	Result   string `json:"result,omitempty" jsonschema:"description=Query result as JSON string (for query operations)"`
}

func NewDBCrudTool() tool.InvokableTool {
	t, err := utils.InferOptionableTool(
		"db_crud",
		"Execute SQL queries against the database and return results in JSON format. Use this tool when you need to query, insert, update or delete data from the database. The results will be formatted as JSON for easy parsing.",
		func(ctx context.Context, input *DBCrudInput, opts ...tool.Option) (output string, err error) {
			db := orm.DB()
			if db == nil {
				out := DBCrudOutput{Success: false, Message: "database not initialized"}
				b, _ := json.MarshalIndent(out, "", "  ")
				return string(b), nil
			}

			if input.OperateType == "query" {
				var results []map[string]any
				if err := db.Raw(input.SQL).Scan(&results).Error; err != nil {
					out := DBCrudOutput{Success: false, Message: fmt.Sprintf("query error: %v", err)}
					b, _ := json.MarshalIndent(out, "", "  ")
					return string(b), nil
				}
				resBytes, _ := json.Marshal(results)
				out := DBCrudOutput{Success: true, Message: "query executed successfully", Result: string(resBytes)}
				b, _ := json.MarshalIndent(out, "", "  ")
				return string(b), nil
			}

			result := db.Exec(input.SQL)
			if result.Error != nil {
				out := DBCrudOutput{Success: false, Message: fmt.Sprintf("exec error: %v", result.Error)}
				b, _ := json.MarshalIndent(out, "", "  ")
				return string(b), nil
			}

			out := DBCrudOutput{
				Success:      true,
				Message:      "operation executed successfully",
				RowsAffected: result.RowsAffected,
			}
			b, _ := json.MarshalIndent(out, "", "  ")
			return string(b), nil
		})
	if err != nil {
		panic(err)
	}
	return t
}
