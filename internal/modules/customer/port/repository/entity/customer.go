package entity

// CustomersTableName TableName
var CustomersTableName = "customers"

type CustomerStatus int

const (
	ActiveCustomerStatus   CustomerStatus = 1 // đã trả lời
	InactiveCustomerStatus CustomerStatus = 2 // chưa trả lời
)

type ServiceType int

const (
	TuyenDung ServiceType = 1
	LienHe    ServiceType = 2
	KhoaHoc   ServiceType = 3
)

// Customer struct
type Customer struct {
	ID           int     `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	CustomerName string  `gorm:"column:customer_name;type:text;not null" mapstructure:"customer_name"`
	Email        string  `gorm:"column:email;type:text;not null" mapstructure:"email"`
	PhoneNumber  string  `gorm:"column:phone_number;type:text;not null" mapstructure:"phone_number"`
	CompanyName  *string `gorm:"column:company_name;type:text" mapstructure:"company_name"`
	Message      *string `gorm:"column:message;type:text" mapstructure:"message"`
	Note         *string `gorm:"column:note;type:text" mapstructure:"note"`
	ServiceType  int     `gorm:"column:service_type;type:int;not null" mapstructure:"service_type"`
	Status       int     `gorm:"column:status;type:int;not null;default:2" mapstructure:"status"`
}

// TableName TableName
func (i *Customer) TableName() string {
	return CustomersTableName
}
