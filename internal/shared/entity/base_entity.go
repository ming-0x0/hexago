package entity

import (
	"time"

	"gorm.io/gorm"
)

type BaseEntity struct {
	CreatedBy string    `gorm:"column:created_by;not null;type:char(26)"`
	CreatedAt time.Time `gorm:"column:created_at;not null;type:datetime;default:current_timestamp"`
	UpdatedBy string    `gorm:"column:updated_by;not null;type:char(26)"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;type:datetime;default:current_timestamp"`
}

type BaseEntityWithDeleted struct {
	BaseEntity
	DeletedBy string         `gorm:"column:deleted_by;type:char(26)"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index"`
}
