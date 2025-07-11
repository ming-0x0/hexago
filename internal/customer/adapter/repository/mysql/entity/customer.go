package entity

import (
	"github.com/ming-0x0/hexago/internal/shared/entity"
	"github.com/ming-0x0/hexago/internal/shared/undefined"
)

var CustomerTable = "customers"

type Customer struct {
	ID           string                      `gorm:"column:id;primaryKey;type:char(26);not null"`
	CustomerName string                      `gorm:"column:customer_name;type:varchar(255);not null"`
	Email        string                      `gorm:"column:email;type:varchar(255);not null"`
	PhoneNumber  string                      `gorm:"column:phone_number;type:char(10);not null"`
	CompanyName  undefined.Undefined[string] `gorm:"column:company_name;type:varchar(255)"`
	Message      undefined.Undefined[string] `gorm:"column:message;type:varchar(1000)"`
	Note         undefined.Undefined[string] `gorm:"column:note;type:varchar(1000)"`
	ServiceType  int64                       `gorm:"column:service_type;type:tinyint(1) unsigned;not null"`
	Status       int64                       `gorm:"column:status;type:tinyint(1) unsigned;not null;default:2"`
	entity.BaseEntityWithDeleted
}
