package model

// 司机表
type LxhDriver struct {
	Id        int32   `gorm:"column:id;type:int;primaryKey;not null;" json:"id"`
	Name      string  `gorm:"column:name;type:varchar(20);comment:姓名;default:NULL;" json:"name"`                 // 姓名
	NickName  string  `gorm:"column:nick_name;type:varchar(20);comment:昵称;default:NULL;" json:"nick_name"`       // 昵称
	AllAcount float64 `gorm:"column:all_acount;type:decimal(10, 2);comment:总收益;default:NULL;" json:"all_acount"` // 总收益
	CarAge    int32   `gorm:"column:car_age;type:int;comment:车龄;default:NULL;" json:"car_age"`                   // 车龄
	Status    string  `gorm:"column:status;type:varchar(10);comment:是否开启接单;default:NULL;" json:"status"`         // 是否开启接单
	Mobile    string  `gorm:"column:mobile;type:char(11);comment:联系电话;default:NULL;" json:"mobile"`              // 联系电话
	FileId    int32   `gorm:"column:file_id;type:int;comment:头像文件ID;default:NULL;" json:"file_id"`               // 头像文件ID
}

func (LxhDriver) TableName() string {
	return "lxh_driver"
}
