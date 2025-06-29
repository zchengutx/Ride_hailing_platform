package model

type LxhDriver struct {
	Id        int32   `gorm:"column:id;type:int;primaryKey;" json:"id"`
	Name      string  `gorm:"column:name;type:varchar(20);comment:姓名;" json:"name"`                 // 姓名
	NickName  string  `gorm:"column:nick_name;type:varchar(20);comment:昵称;" json:"nick_name"`       // 昵称
	AllAcount float64 `gorm:"column:all_acount;type:decimal(10, 2);comment:总收益;" json:"all_acount"` // 总收益
	CarAge    int32   `gorm:"column:car_age;type:int;comment:车龄;" json:"car_age"`                   // 车龄
	Status    string  `gorm:"column:status;type:varchar(10);comment:是否开启接单;" json:"status"`         // 是否开启接单
	Mobile    string  `gorm:"column:mobile;type:char(11);comment:联系电话;" json:"mobile"`              // 联系电话
	FileId    int32   `gorm:"column:file_id;type:int;comment:头像文件ID;" json:"file_id"`               // 头像文件ID
}
type LxhDriverCheck struct {
	Id                   int32  `gorm:"column:id;type:int;primaryKey;" json:"id"`
	IdCardFileId         string `gorm:"column:id_card_file_id;type:varchar(100);comment:身份证;" json:"id_card_file_id"`                 // 身份证
	DriverLicenseFileId  string `gorm:"column:driver_license_file_id;type:varchar(100);comment:驾照;" json:"driver_license_file_id"`    // 驾照
	DrivingLicenseFileId string `gorm:"column:driving_license_file_id;type:varchar(100);comment:行驶证;" json:"driving_license_file_id"` // 行驶证
	AvatorFileId         string `gorm:"column:avator_file_id;type:varchar(100);comment:司机头像;" json:"avator_file_id"`                  // 司机头像
	CheckStatus          string `gorm:"column:check_status;type:varchar(5);comment:审核状态;" json:"check_status"`                        // 审核状态
	Remark               string `gorm:"column:remark;type:varchar(50);comment:备注;" json:"remark"`                                     // 备注
	DriverId             int32  `gorm:"column:driver_id;type:int;comment:司机id;not null;default:0;" json:"driver_id"`                  // 司机id
}

func (d *LxhDriver) TableName() string {
	return "lxh_driver"
}

func (c *LxhDriverCheck) TableName() string {
	return "lxh_driver_check"
}
