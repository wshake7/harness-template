package mixin

// AutoIncrementID 是 GORM 可复用的 mixin，包含自增主键字段。
// 在模型中嵌入：
//
//	type User struct {
//	    mixin.AutoIncrementID
//	    Name string
//	}
type AutoIncrementID struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
}
