package model

import "time"

type OperationLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	AdminID     uint      `gorm:"column:admin_id" json:"-"`
	Username    string    `gorm:"column:username" json:"username"`
	Path        string    `gorm:"column:path" json:"path"`
	IP          string    `gorm:"column:ip" json:"ip"`
	RequestData string    `gorm:"column:request_data" json:"request_data"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
}

func (a OperationLog) GetID() uint {
	return a.ID
}
