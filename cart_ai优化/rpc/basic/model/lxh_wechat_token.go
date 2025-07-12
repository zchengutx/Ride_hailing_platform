package model

import "time"

// 微信授权令牌表
type LxhWechatToken struct {
	Id           int64     `gorm:"column:id;type:int;primaryKey;not null;" json:"id"`
	Openid       string    `gorm:"column:openid;type:varchar(64);index;not null;comment:微信用户标识;" json:"openid"`
	AccessToken  string    `gorm:"column:access_token;type:varchar(512);comment:访问令牌;default:NULL;" json:"access_token"`
	RefreshToken string    `gorm:"column:refresh_token;type:varchar(512);comment:刷新令牌;default:NULL;" json:"refresh_token"`
	ExpiresIn    int       `gorm:"column:expires_in;type:int;comment:过期时间(秒);default:0;" json:"expires_in"`
	Scope        string    `gorm:"column:scope;type:varchar(100);comment:授权作用域;default:NULL;" json:"scope"`
	TokenType    string    `gorm:"column:token_type;type:varchar(20);comment:令牌类型;default:Bearer;" json:"token_type"`
	ExpiresAt    time.Time `gorm:"column:expires_at;type:datetime;comment:令牌过期时间;default:NULL;" json:"expires_at"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"updated_at"`
}

func (l *LxhWechatToken) TableName() string {
	return "lxh_wechat_token"
}
