package model

import "time"

// 微信用户表
type LxhWechatUser struct {
	Id         int64     `gorm:"column:id;type:int;primaryKey;not null;" json:"id"`
	Openid     string    `gorm:"column:openid;type:varchar(64);uniqueIndex;not null;comment:微信用户唯一标识;" json:"openid"`
	Nickname   string    `gorm:"column:nickname;type:varchar(100);comment:微信昵称;default:NULL;" json:"nickname"`
	Headimgurl string    `gorm:"column:headimgurl;type:varchar(500);comment:头像URL;default:NULL;" json:"headimgurl"`
	Sex        int       `gorm:"column:sex;type:tinyint;comment:性别(0未知1男2女);default:0;" json:"sex"`
	Country    string    `gorm:"column:country;type:varchar(50);comment:国家;default:NULL;" json:"country"`
	Province   string    `gorm:"column:province;type:varchar(50);comment:省份;default:NULL;" json:"province"`
	City       string    `gorm:"column:city;type:varchar(50);comment:城市;default:NULL;" json:"city"`
	Language   string    `gorm:"column:language;type:varchar(20);comment:语言;default:NULL;" json:"language"`
	Privilege  string    `gorm:"column:privilege;type:text;comment:特权信息(JSON);default:NULL;" json:"privilege"`
	UnionId    string    `gorm:"column:unionid;type:varchar(64);comment:开放平台统一标识;default:NULL;" json:"unionid"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"updated_at"`
}

func (l *LxhWechatUser) TableName() string {
	return "lxh_wechat_user"
}
