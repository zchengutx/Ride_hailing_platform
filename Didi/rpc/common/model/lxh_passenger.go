package model

// 乘客表
type LxhPassenger struct {
	Id       int64  `gorm:"column:id;type:int;primaryKey;not null;" json:"id"`
	Name     string `gorm:"column:name;type:varchar(20);comment:姓名;default:NULL;" json:"name"`           // 姓名
	NickName string `gorm:"column:nick_name;type:varchar(20);comment:昵称;default:NULL;" json:"nick_name"` // 昵称
	FileId   int64  `gorm:"column:file_id;type:int;comment:头像文件ID;default:NULL;" json:"file_id"`         // 头像文件ID
	Mobile   string `gorm:"column:mobile;type:char(11);comment:手机号;default:NULL;" json:"mobile"`         // 手机号
	Age      int64  `gorm:"column:age;type:int;comment:年龄;default:NULL;" json:"age"`                     // 年龄
	Sex      string `gorm:"column:sex;type:varchar(2);comment:性别;default:NULL;" json:"sex"`              // 性别
	Mileage  string `gorm:"column:mileage;type:varchar(10);comment:里程数;default:NULL;" json:"mileage"`    // 里程数
}

func (p *LxhPassenger) TableName() string {
	return "lxh_passenger"
}
