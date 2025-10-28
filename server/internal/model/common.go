package model

import (
	"time"
)

type Commonmodel struct {
	ID        uint      `gorm:"column:id;primaryKey;type:uint(12);comment:ID" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updated_at"`
	// DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// GetID 获取ID
func (c Commonmodel) GetID() uint {
	return c.ID
}

// Ids 用于批量操作的ID集合
type IDs struct {
	IDs []uint `json:"ids"  binding:"required"`
}
type ID struct {
	ID uint `json:"id"  binding:"required"`
}
