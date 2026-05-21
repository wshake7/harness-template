package money

import (
	"github.com/shopspring/decimal"
	"strconv"
)

// Ceil2 保留两位小数并四舍五入
func Ceil2(v string) decimal.Decimal {
	d, _ := decimal.NewFromString(v)
	return d.Mul(decimal.NewFromInt(100)).Ceil().Div(decimal.NewFromInt(100))
}

// Ceil2Float64 浮点数保留两位小数并四舍五入
func Ceil2Float64(v float64) decimal.Decimal {
	// 保留足够精度避免丢失信息
	s := strconv.FormatFloat(v, 'f', -1, 64)

	d, _ := decimal.NewFromString(s)
	return d.Mul(decimal.NewFromInt(100)).Ceil().Div(decimal.NewFromInt(100))
}
