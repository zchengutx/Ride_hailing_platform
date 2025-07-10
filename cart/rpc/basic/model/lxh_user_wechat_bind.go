package model

import "time"

// 用户微信绑定表
type LxhUserWechatBind struct {
	Id          int64     `gorm:"column:id;type:int;primaryKey;not null;" json:"id"`
	PassengerId int64     `gorm:"column:passenger_id;type:int;index;not null;comment:乘客ID;" json:"passenger_id"`
	Openid      string    `gorm:"column:openid;type:varchar(64);uniqueIndex;not null;comment:微信用户标识;" json:"openid"`
	BindType    int       `gorm:"column:bind_type;type:tinyint;comment:绑定类型(1微信授权2手机绑定);default:1;" json:"bind_type"`
	BindStatus  int       `gorm:"column:bind_status;type:tinyint;comment:绑定状态(0未绑定1已绑定2已解绑);default:1;" json:"bind_status"`
	BindTime    time.Time `gorm:"column:bind_time;type:datetime;comment:绑定时间;default:CURRENT_TIMESTAMP;" json:"bind_time"`
	UnbindTime  time.Time `gorm:"column:unbind_time;type:datetime;comment:解绑时间;default:NULL;" json:"unbind_time"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"updated_at"`
}

func (l *LxhUserWechatBind) TableName() string {
	return "lxh_user_wechat_bind"
}
