package model

import "kitex_main/config"

type Auth struct {
	Id          int32  `gorm:"column:id;type:int;primaryKey;" json:"id"`
	AppName     string `gorm:"column:app_name;type:varchar(20);" json:"app_name"`
	AppUnionid  string `gorm:"column:app_unionid;type:varchar(60);" json:"app_unionid"`
	AccessToken string `gorm:"column:access_token;type:varchar(60);" json:"access_token"`
	AppRemark   string `gorm:"column:app_remark;type:varchar(100);" json:"app_remark"`
}

// 定义结构体名称
func (u *Auth) TableName() string {
	return "auth"
}

// 查询用户是否存在
func (u *Auth) FindAuthOpenId(OpenId string) (err error) {
	err = config.DB.Where("app_unionid = ?", OpenId).First(u).Error
	return
}

// 创建Auth
func (u *Auth) CreateAuth() (err error) {
	err = config.DB.Create(u).Error
	return
}
