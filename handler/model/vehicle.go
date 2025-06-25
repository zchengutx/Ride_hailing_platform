package model

import "time"

// 车辆表
type Vehicle struct {
	Id          string    `gorm:"column:id;type:varchar(191);primaryKey;not null;" json:"id"`
	DriverID    string    `gorm:"column:driver_id;type:varchar(191);not null;" json:"driver_id"`
	PlateNumber string    `gorm:"column:plate_number;type:varchar(20);not null;" json:"plate_number"`
	Brand       string    `gorm:"column:brand;type:varchar(50);not null;" json:"brand"`
	Model       string    `gorm:"column:model;type:varchar(50);not null;" json:"model"`
	Color       string    `gorm:"column:color;type:varchar(20);not null;" json:"color"`
	Seats       int8      `gorm:"column:seats;type:tinyint;default:4;" json:"seats"`
	VehicleType string    `gorm:"column:vehicle_type;type:varchar(20);default:'economy';" json:"vehicle_type"`
	Status      int8      `gorm:"column:status;type:tinyint;default:1;" json:"status"`
	CreateTime  time.Time `gorm:"column:create_time;type:datetime(3);default:CURRENT_TIMESTAMP(3);" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime(3);default:CURRENT_TIMESTAMP(3);" json:"update_time"`

	// 关联
	LxhDriver *LxhDriver `gorm:"foreignKey:DriverID;references:Id" json:"driver,omitempty"`
}

func (Vehicle) TableName() string {
	return "vehicle"
}
