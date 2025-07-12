package model

// 司机资料审核表
type LxhDriverCheck struct {
	Id                   int64  `gorm:"column:id;type:int;primaryKey;not null;" json:"id"`
	IdCardFileId         string `gorm:"column:id_card_file_id;type:varchar(100);comment:身份证;default:NULL;" json:"id_card_file_id"`                 // 身份证
	DriverLicenseFileId  string `gorm:"column:driver_license_file_id;type:varchar(100);comment:驾照;default:NULL;" json:"driver_license_file_id"`    // 驾照
	DrivingLicenseFileId string `gorm:"column:driving_license_file_id;type:varchar(100);comment:行驶证;default:NULL;" json:"driving_license_file_id"` // 行驶证
	AvatorFileId         string `gorm:"column:avator_file_id;type:varchar(100);comment:司机头像;default:NULL;" json:"avator_file_id"`                  // 司机头像
	CheckStatus          string `gorm:"column:check_status;type:varchar(5);comment:审核状态;default:NULL;" json:"check_status"`                        // 审核状态
	Remark               string `gorm:"column:remark;type:varchar(50);comment:备注;default:NULL;" json:"remark"`                                     // 备注
}

func (l *LxhDriverCheck) TableName() string {
	return "lxh_driver_check"
}
