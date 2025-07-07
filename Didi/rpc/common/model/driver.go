package model

import (
	"gorm.io/gorm"
	"time"
)

type Driver struct {
	Id                    uint64         `gorm:"column:id;type:int UNSIGNED;primaryKey;not null;" json:"id"`
	DriverId              uint64         `gorm:"column:driver_id;type:int UNSIGNED;comment:司机ID;not null;" json:"driver_id"`                                // 司机ID
	Mobile                string         `gorm:"column:mobile;type:char(11);comment:手机号码;not null;" json:"mobile"`                                          // 手机号码
	IdCard                string         `gorm:"column:id_card;type:char(18);comment:身份证号码;not null;" json:"id_card"`                                       // 身份证号码
	CarNum                string         `gorm:"column:car_num;type:char(8);comment:车牌号码;not null;" json:"car_num"`                                         // 车牌号码
	CarType               string         `gorm:"column:car_type;type:varchar(10);comment:车型;not null;" json:"car_type"`                                     // 车型
	DrivingLicenseNumber  string         `gorm:"column:driving_license_number;type:char(18);comment:驾驶证号;not null;" json:"driving_license_number"`          // 驾驶证号
	DrivingAge            uint64         `gorm:"column:driving_age;type:int UNSIGNED;comment:驾龄;default:NULL;" json:"driving_age"`                          // 驾龄
	BackgroundCheckStatus string         `gorm:"column:background_check_status;type:varchar(10);comment:背景查询状态;default:正常;" json:"background_check_status"` // 背景查询状态
	DriverStatus          string         `gorm:"column:driver_status;type:varchar(10);comment:司机状态;default:离线;" json:"driver_status"`                       // 司机状态
	Rating                string         `gorm:"column:rating;type:varchar(10);comment:评分;default:NULL;" json:"rating"`                                     // 评分
	TotalOrders           uint64         `gorm:"column:total_orders;type:int UNSIGNED;comment:总接单数;default:NULL;" json:"total_orders"`                      // 总接单数
	CreatedAt             time.Time      `gorm:"column:created_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"created_at"`
	UpdatedAt             time.Time      `gorm:"column:updated_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"column:deleted_at;type:datetime(3);default:NULL;" json:"deleted_at"`
}

func (d *Driver) TableName() string {
	return "driver"
}
