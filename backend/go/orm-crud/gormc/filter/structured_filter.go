package filter

import (
	"github.com/bytedance/sonic"
	"go.uber.org/zap"
	"strings"

	"gorm.io/gorm"

	"go-common/utils/stringcase"

	paginationV1 "orm-crud/api/gen/go/pagination/v1"
)

// StructuredFilter 基于 FilterExpr 的 GORM 过滤器
type StructuredFilter struct {
	Unmarshal func([]byte, any) error
	Marshal   func(any) ([]byte, error)
	processor *Processor
}

func NewStructuredFilter() *StructuredFilter {
	return &StructuredFilter{
		Unmarshal: sonic.Unmarshal,
		Marshal:   sonic.Marshal,
		processor: NewProcessor(),
	}
}

// BuildSelectors 将 FilterExpr 转为一组可应用于 *gorm.DB 的闭包
func (sf StructuredFilter) BuildSelectors(expr *paginationV1.FilterExpr) ([]func(*gorm.DB) *gorm.DB, error) {
	var sels []func(*gorm.DB) *gorm.DB

	if expr == nil {
		// 返回空 slice 以保持兼容测试（也可返回 nil）
		return sels, nil
	}

	// 未指定类型视为跳过（测试期望返回 nil）
	if expr.GetType() == paginationV1.ExprType_EXPR_TYPE_UNSPECIFIED {
		zap.L().Warn("Skipping unspecified FilterExpr")
		return nil, nil
	}

	sel, err := sf.buildFilterSelector(expr)
	if err != nil {
		return nil, err
	}
	if sel != nil {
		sels = append(sels, sel)
	}
	return sels, nil
}

// buildFilterSelector 将单个 FilterExpr 转为 *gorm.DB 闭包（递归处理组）
func (sf StructuredFilter) buildFilterSelector(expr *paginationV1.FilterExpr) (func(*gorm.DB) *gorm.DB, error) {
	if expr == nil {
		zap.L().Warn("Skipping nil FilterExpr")
		return nil, nil
	}
	if expr.GetType() == paginationV1.ExprType_EXPR_TYPE_UNSPECIFIED {
		zap.L().Warn("Skipping unspecified FilterExpr")
		return nil, nil
	}

	// helper: 将单个 Condition 应用到 db 上
	applyCond := func(db *gorm.DB, cond *paginationV1.FilterCondition) *gorm.DB {
		if db == nil || cond == nil {
			return db
		}
		val := ""
		switch cond.ValueOneof.(type) {
		case *paginationV1.FilterCondition_Value:
			val = cond.GetValue()
		default:
		}

		// 支持 JSON 字段 (e.g. preferences.daily_email)
		if strings.Contains(cond.GetField(), ".") {
			parts := strings.SplitN(cond.GetField(), ".", 2)
			col := stringcase.ToSnakeCase(parts[0])
			jsonKey := parts[1]
			exprStr, _ := sf.processor.JsonbFieldExpr(db, jsonKey, col)
			if exprStr == "" {
				return db
			}
			return sf.processor.Process(db, cond.GetOp(), exprStr, val, cond.GetValues())
		}

		col := stringcase.ToSnakeCase(cond.GetField())

		// 需要 CAST 的字符串模糊匹配操作
		stringLikeOps := map[paginationV1.Operator]bool{
			paginationV1.Operator_CONTAINS:     true,
			paginationV1.Operator_ICONTAINS:    true,
			paginationV1.Operator_STARTS_WITH:  true,
			paginationV1.Operator_ISTARTS_WITH: true,
			paginationV1.Operator_ENDS_WITH:    true,
			paginationV1.Operator_IENDS_WITH:   true,
			paginationV1.Operator_LIKE:         true,
			paginationV1.Operator_ILIKE:        true,
			paginationV1.Operator_NOT_LIKE:     true,
		}

		if stringLikeOps[cond.GetOp()] {
			isPostgres := db.Dialector.Name() == "postgres"

			var castCol string
			if isPostgres {
				castCol = "CAST(" + col + " AS TEXT)"
			} else {
				castCol = "CAST(" + col + " AS CHAR)"
			}

			switch cond.GetOp() {
			case paginationV1.Operator_ICONTAINS:
				if isPostgres {
					return db.Where(castCol+" ILIKE ?", "%"+val+"%")
				}
				return db.Where(castCol+" LIKE ?", "%"+val+"%")

			case paginationV1.Operator_CONTAINS:
				return db.Where(castCol+" LIKE ?", "%"+val+"%")

			case paginationV1.Operator_ISTARTS_WITH:
				if isPostgres {
					return db.Where(castCol+" ILIKE ?", val+"%")
				}
				return db.Where(castCol+" LIKE ?", val+"%")

			case paginationV1.Operator_STARTS_WITH:
				return db.Where(castCol+" LIKE ?", val+"%")

			case paginationV1.Operator_IENDS_WITH:
				if isPostgres {
					return db.Where(castCol+" ILIKE ?", "%"+val)
				}
				return db.Where(castCol+" LIKE ?", "%"+val)

			case paginationV1.Operator_ENDS_WITH:
				return db.Where(castCol+" LIKE ?", "%"+val)

			case paginationV1.Operator_ILIKE:
				if isPostgres {
					return db.Where(castCol+" ILIKE ?", val)
				}
				return db.Where(castCol+" LIKE ?", val)

			case paginationV1.Operator_LIKE:
				return db.Where(castCol+" LIKE ?", val)

			case paginationV1.Operator_NOT_LIKE:
				return db.Where(castCol+" NOT LIKE ?", val)
			}
		}

		return sf.processor.Process(db, cond.GetOp(), col, val, cond.GetValues())
	}

	// 构造闭包
	closure := func(db *gorm.DB) *gorm.DB {
		if db == nil {
			return db
		}

		switch expr.GetType() {
		case paginationV1.ExprType_AND:
			for _, cond := range expr.GetConditions() {
				db = applyCond(db, cond)
			}
			for _, g := range expr.GetGroups() {
				subSel, err := sf.buildFilterSelector(g)
				if err != nil {
					zap.S().Errorf("buildFilterSelector sub-group error: %v", err)
					continue
				}
				if subSel != nil {
					db = subSel(db)
				}
			}
			return db

		case paginationV1.ExprType_OR:
			orDB := db.Session(&gorm.Session{NewDB: true})
			first := true

			for _, cond := range expr.GetConditions() {
				c := cond
				if first {
					orDB = applyCond(orDB, c)
					first = false
				} else {
					branch := applyCond(db.Session(&gorm.Session{NewDB: true}), c)
					orDB = orDB.Or(branch)
				}
			}

			for _, g := range expr.GetGroups() {
				subSel, err := sf.buildFilterSelector(g)
				if err != nil {
					zap.S().Errorf("buildFilterSelector sub-group error: %v", err)
					continue
				}
				if subSel == nil {
					continue
				}
				branch := subSel(db.Session(&gorm.Session{NewDB: true}))
				if first {
					orDB = branch
					first = false
				} else {
					orDB = orDB.Or(branch)
				}
			}

			return db.Where(orDB)

		default:
			return db
		}
	}

	return closure, nil
}
