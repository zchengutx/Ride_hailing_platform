package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id        uint64         `gorm:"column:id;type:int UNSIGNED;primaryKey;not null;" json:"id"`
	Mobile    string         `gorm:"column:mobile;type:char(11);comment:手机号;not null;" json:"mobile"` // 手机号
	UserName  string         `gorm:"column:user_name;type:varchar(30);not null;" json:"user_name"`
	NikeName  string         `gorm:"column:nike_name;type:varchar(30);not null;" json:"nike_name"`
	Sex       string         `gorm:"column:sex;type:varchar(5);comment:性别;default:NULL;" json:"sex"`   // 性别
	Age       uint64         `gorm:"column:age;type:int UNSIGNED;comment:年龄;default:NULL;" json:"age"` // 年龄
	Mileage   string         `gorm:"column:mileage;type:varchar(10);not null;" json:"mileage"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime(3);default:NULL;" json:"deleted_at"`
}

func (u *User) TableName() string {
	return "user"
}
