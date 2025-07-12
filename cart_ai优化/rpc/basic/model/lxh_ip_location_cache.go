package model

import "time"

// IP定位缓存表
type LxhIpLocationCache struct {
	Id           int64     `gorm:"column:id;type:int;primaryKey;not null;" json:"id"`
	IpAddress    string    `gorm:"column:ip_address;type:varchar(45);uniqueIndex;not null;comment:IP地址;" json:"ip_address"`
	Address      string    `gorm:"column:address;type:varchar(500);comment:详细地址;default:NULL;" json:"address"`
	Province     string    `gorm:"column:province;type:varchar(50);comment:省份;default:NULL;" json:"province"`
	City         string    `gorm:"column:city;type:varchar(50);comment:城市;default:NULL;" json:"city"`
	District     string    `gorm:"column:district;type:varchar(50);comment:区县;default:NULL;" json:"district"`
	Street       string    `gorm:"column:street;type:varchar(100);comment:街道;default:NULL;" json:"street"`
	StreetNumber string    `gorm:"column:street_number;type:varchar(100);comment:门牌号;default:NULL;" json:"street_number"`
	CityCode     int       `gorm:"column:city_code;type:int;comment:城市代码;default:0;" json:"city_code"`
	Lng          float64   `gorm:"column:lng;type:decimal(11,8);comment:经度;default:NULL;" json:"lng"`
	Lat          float64   `gorm:"column:lat;type:decimal(11,8);comment:纬度;default:NULL;" json:"lat"`
	Status       int       `gorm:"column:status;type:int;comment:API返回状态;default:0;" json:"status"`
	Source       string    `gorm:"column:source;type:varchar(50);comment:数据来源;default:baidu;" json:"source"`
	HitCount     int       `gorm:"column:hit_count;type:int;comment:命中次数;default:1;" json:"hit_count"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"updated_at"`
}

func (l *LxhIpLocationCache) TableName() string {
	return "lxh_ip_location_cache"
}
