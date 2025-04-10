package models

import "time"

type BaseModel struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement;" json:"id,omitempty"`
}

// CommonTimestampsField 时间戳
type CommonTimestampsField struct {
	CreatedAt time.Time `gorm:"colume:created_at;index;" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"colume:updated_at;index;" json:"updated_at,omitempty"`
}
