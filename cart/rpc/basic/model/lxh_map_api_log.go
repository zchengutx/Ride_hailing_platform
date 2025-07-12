package model

import "time"

// 地图API调用日志表
type LxhMapApiLog struct {
	Id            int64     `gorm:"column:id;type:int;primaryKey;not null;" json:"id"`
	ApiType       string    `gorm:"column:api_type;type:varchar(50);index;comment:API类型(geocoding/reverse_geocoding/ip_location/distance);default:NULL;" json:"api_type"`
	RequestParams string    `gorm:"column:request_params;type:text;comment:请求参数(JSON);default:NULL;" json:"request_params"`
	ResponseData  string    `gorm:"column:response_data;type:text;comment:响应数据(JSON);default:NULL;" json:"response_data"`
	ResponseCode  int       `gorm:"column:response_code;type:int;comment:响应状态码;default:0;" json:"response_code"`
	ResponseTime  int       `gorm:"column:response_time;type:int;comment:响应时间(毫秒);default:0;" json:"response_time"`
	IpAddress     string    `gorm:"column:ip_address;type:varchar(45);comment:客户端IP;default:NULL;" json:"ip_address"`
	UserAgent     string    `gorm:"column:user_agent;type:varchar(500);comment:用户代理;default:NULL;" json:"user_agent"`
	UserId        int64     `gorm:"column:user_id;type:int;index;comment:用户ID;default:NULL;" json:"user_id"`
	ErrorMessage  string    `gorm:"column:error_message;type:text;comment:错误信息;default:NULL;" json:"error_message"`
	IsSuccess     int       `gorm:"column:is_success;type:tinyint;comment:是否成功(0失败1成功);default:1;" json:"is_success"`
	CreatedAt     time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;default:CURRENT_TIMESTAMP;" json:"created_at"`
}

func (l *LxhMapApiLog) TableName() string {
	return "lxh_map_api_log"
}

// BaiduMapClient 百度地图客户端结构体
// 用于 utils/baidumap.go 和 mapserivce/handler.go
type BaiduMapClient struct {
	APIKey  string
	APIHost string
}
