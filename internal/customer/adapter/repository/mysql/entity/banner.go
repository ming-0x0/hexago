package entity

import "github.com/ming-0x0/hexago/internal/shared/entity"

var BannerTable = "banners"

type Banner struct {
	ID   string `gorm:"column:id;primaryKey;type:char(26);not null"`
	Name string `gorm:"column:name;type:varchar(255);not null"`
	entity.BaseEntityWithDeleted
}

func (Banner) TableName() string {
	return BannerTable
}
