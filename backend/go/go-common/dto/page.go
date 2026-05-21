package dto

import (
	"go-common/utils/str"
)

type BasicPageDto struct {
	PageNum  int `json:"pageNum" binding:"lte=100"`                 // 当前页码
	PageSize int `json:"pageSize" binding:"required,gte=1,lte=100"` // 每页展示数量
}

func (d BasicPageDto) Offset() int {
	return d.PageSize * d.PageNum
}

type OrderPageDto struct {
	BasicPageDto
	OrderBy   string `binding:"required"`                // 排序字段
	OrderType string `binding:"required,oneof=DESC ASC"` // 排序方式 DESC/ASC
}

func (dto *OrderPageDto) GetOrder() string {
	return str.Join(dto.OrderBy, " ", dto.OrderType)
}

type IdPageDto struct {
	LastId   uint64 `json:"lastId"`
	PageSize int    `json:"pageSize" binding:"required,gte=1,lte=100"` // 每页展示数量
}
