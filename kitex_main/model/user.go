package model

import (
	"kitex_main/config"
)

type LxhPassenger struct {
	Id       int32  `gorm:"column:id;type:int;primaryKey;" json:"id"`
	Name     string `gorm:"column:name;type:varchar(20);comment:姓名;" json:"name"`           // 姓名
	NickName string `gorm:"column:nick_name;type:varchar(20);comment:昵称;" json:"nick_name"` // 昵称
	FileId   int32  `gorm:"column:file_id;type:int;comment:头像文件ID;" json:"file_id"`       // 头像文件ID
	Mobile   string `gorm:"column:mobile;type:char(11);comment:手机号;" json:"mobile"`        // 手机号
	Age      int32  `gorm:"column:age;type:int;comment:年龄;" json:"age"`                     // 年龄
	Sex      string `gorm:"column:sex;type:varchar(2);comment:性别;" json:"sex"`              // 性别
	Mileage  string `gorm:"column:mileage;type:varchar(10);comment:里程数;" json:"mileage"`   // 里程数
}

// 初始化结构体名称
func (u *LxhPassenger) TableName() string {
	return "lxh_passenger"
}

// 查询改手机号是否存在
func (u *LxhPassenger) FindUserMobile(mobile string) (err error) {
	err = config.DB.Where("mobile = ?", mobile).First(u).Error
	return
}

// 添加用户
func (u *LxhPassenger) CreateUser() (err error) {
	err = config.DB.Create(u).Error
	return
}

// 查询用户Id
func (u *LxhPassenger) FindUserId(id int) (err error) {
	err = config.DB.Where("id = ?", id).First(u).Error
	return
}
