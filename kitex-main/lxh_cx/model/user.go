package model

type LxhPassenger struct {
	Id       int32  `gorm:"column:id;type:int;primaryKey;" json:"id"`
	Name     string `gorm:"column:name;type:varchar(20);comment:姓名;" json:"name"`           // 姓名
	NickName string `gorm:"column:nick_name;type:varchar(20);comment:昵称;" json:"nick_name"` // 昵称
	FileId   int32  `gorm:"column:file_id;type:int;comment:头像文件ID;" json:"file_id"`         // 头像文件ID
	Mobile   string `gorm:"column:mobile;type:char(11);comment:手机号;" json:"mobile"`         // 手机号
	Age      int32  `gorm:"column:age;type:int;comment:年龄;" json:"age"`                     // 年龄
	Sex      string `gorm:"column:sex;type:varchar(2);comment:性别;" json:"sex"`              // 性别
	Mileage  string `gorm:"column:mileage;type:varchar(10);comment:里程数;" json:"mileage"`    // 里程数
}

type LxhAuth struct {
	Id          int32  `gorm:"column:id;type:int;primaryKey;" json:"id"`
	AppName     string `gorm:"column:app_name;type:varchar(20);" json:"app_name"`
	AppUnionid  string `gorm:"column:app_unionid;type:varchar(60);" json:"app_unionid"`
	AccessToken string `gorm:"column:access_token;type:varchar(60);" json:"access_token"`
	AppRemark   string `gorm:"column:app_remark;type:varchar(100);" json:"app_remark"`
}

type LxhUserAuth struct {
	Id          int32 `gorm:"column:id;type:int;primaryKey;" json:"id"`
	PassengerId int32 `gorm:"column:passenger_id;type:int;comment:乘客id;not null;default:0;" json:"passenger_id"` // 乘客id
	AuthId      int32 `gorm:"column:auth_id;type:int;comment:认证id;not null;default:0;" json:"auth_id"`           // 认证id
}

func (u *LxhPassenger) TableName() string {
	return "lxh_passenger"
}
func (ua *LxhUserAuth) TableName() string {
	return "lxh_user_auth"
}
func (a *LxhAuth) TableName() string {
	return "lxh_auth"
}
