package model

import (
	"gorm.io/gorm"
	"time"
)

type DriverApplication struct {
	Id                   uint64         `gorm:"column:id;type:int UNSIGNED;primaryKey;not null;" json:"id"`
	Driver               uint64         `gorm:"column:driver;type:int UNSIGNED;comment:司机;not null;" json:"driver"`                                       // 司机
	Mobile               string         `gorm:"column:mobile;type:char(11);comment:联系电话;not null;" json:"mobile"`                                         // 联系电话
	IdCard               string         `gorm:"column:id_card;type:char(18);comment:身份证号码;not null;" json:"id_card"`                                      // 身份证号码
	Age                  uint64         `gorm:"column:age;type:int;comment:年龄;default:NULL;" json:"age"`                                                  // 年龄
	Sex                  string         `gorm:"column:sex;type:varchar(10);comment:性别;default:NULL;" json:"sex"`                                          // 性别
	Address              string         `gorm:"column:address;type:varchar(100);comment:家庭地址;not null;" json:"address"`                                   // 家庭地址
	DrivingLicenseNumber string         `gorm:"column:driving_license_number;type:char(50);comment:驾驶证号;not null;" json:"driving_license_number"`         // 驾驶证号
	QuasiDrivingType     string         `gorm:"column:quasi_driving_type;type:varchar(10);comment:准驾车型;not null;" json:"quasi_driving_type"`              // 准驾车型
	DrivingAge           uint64         `gorm:"column:driving_age;type:int;comment:驾龄;default:NULL;" json:"driving_age"`                                  // 驾龄
	CarNum               string         `gorm:"column:car_num;type:char(8);comment:车牌号;not null;" json:"car_num"`                                         // 车牌号
	CarType              string         `gorm:"column:car_type;type:varchar(10);comment:车辆品牌类型;not null;" json:"car_type"`                                // 车辆品牌类型
	VehicleMileage       string         `gorm:"column:vehicle_mileage;type:varchar(50);comment:车辆行驶里程;not null;" json:"vehicle_mileage"`                  // 车辆行驶里程
	ServingTheCity       string         `gorm:"column:serving_the_city;type:varchar(50);comment:服务城市;not null;" json:"serving_the_city"`                  // 服务城市
	CreatedAt            time.Time      `gorm:"column:created_at;type:datetime(3);comment:注册时间;not null;default:CURRENT_TIMESTAMP(3);" json:"created_at"` // 注册时间
	AuditStatus          string         `gorm:"column:audit_status;type:varchar(10);comment:审核状态;default:待审核;" json:"audit_status"`                       // 审核状态
	Common               string         `gorm:"column:common;type:varchar(100);comment:评语;default:NULL;" json:"common"`                                   // 评语
	UpdatedAt            time.Time      `gorm:"column:updated_at;type:datetime(3);comment:更新时间;not null;default:CURRENT_TIMESTAMP(3);" json:"updated_at"` // 更新时间
	DeletedAt            gorm.DeletedAt `gorm:"column:deleted_at;type:datetime(3);comment:删除时间;default:NULL;" json:"deleted_at"`                          // 删除时间
}
