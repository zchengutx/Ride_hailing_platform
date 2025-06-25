package model

// 乘客信息表
type LxhPassenger struct {
	Id             int32  `gorm:"column:id;type:int;primaryKey;not null;" json:"id"`
	Name           string `gorm:"column:name;type:varchar(20);comment:姓名;default:NULL;" json:"name"`           // 姓名
	NickName       string `gorm:"column:nick_name;type:varchar(20);comment:昵称;default:NULL;" json:"nick_name"` // 昵称
	FileId         int32  `gorm:"column:file_id;type:int;comment:头像文件ID;default:NULL;" json:"file_id"`         // 头像文件ID
	Mobile         string `gorm:"column:mobile;type:char(11);comment:手机号;default:NULL;" json:"mobile"`         // 手机号
	Age            int32  `gorm:"column:age;type:int;comment:年龄;default:NULL;" json:"age"`                     // 年龄
	Sex            string `gorm:"column:sex;type:varchar(2);comment:性别;default:NULL;" json:"sex"`              // 性别
	Mileage        string `gorm:"column:mileage;type:varchar(10);comment:里程数;default:NULL;" json:"mileage"`    // 里程数
	WechatOpenID   string `gorm:"column:wechat_open_id;type:varchar(100);" json:"wechat_open_id"`
	WechatUnionID  string `gorm:"column:wechat_union_id;type:varchar(100);" json:"wechat_union_id"`
	WechatNickname string `gorm:"column:wechat_nickname;type:varchar(100);" json:"wechat_nickname"`
	WechatAvatar   string `gorm:"column:wechat_avatar;type:varchar(255);" json:"wechat_avatar"`
}

// 微信用户信息
type WechatUserInfo struct {
	OpenID     string   `json:"openid"`
	UnionID    string   `json:"unionid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Language   string   `json:"language"`
}

// 微信访问令牌响应
type WechatAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
	ErrorCode    int    `json:"errcode"`
	ErrorMsg     string `json:"errmsg"`
}
