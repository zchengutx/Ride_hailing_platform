package model

import "time"

// 逆地理编码缓存表
type LxhReverseGeocodingCache struct {
	Id               int64     `gorm:"column:id;type:int;primaryKey;not null;" json:"id"`
	LocationKey      string    `gorm:"column:location_key;type:varchar(100);uniqueIndex;not null;comment:坐标键值;" json:"location_key"`
	Lng              float64   `gorm:"column:lng;type:decimal(11,8);comment:经度;default:NULL;" json:"lng"`
	Lat              float64   `gorm:"column:lat;type:decimal(11,8);comment:纬度;default:NULL;" json:"lat"`
	FormattedAddress string    `gorm:"column:formatted_address;type:varchar(500);comment:格式化地址;default:NULL;" json:"formatted_address"`
	Business         string    `gorm:"column:business;type:varchar(200);comment:商圈信息;default:NULL;" json:"business"`
	Country          string    `gorm:"column:country;type:varchar(50);comment:国家;default:NULL;" json:"country"`
	Province         string    `gorm:"column:province;type:varchar(50);comment:省份;default:NULL;" json:"province"`
	City             string    `gorm:"column:city;type:varchar(50);comment:城市;default:NULL;" json:"city"`
	District         string    `gorm:"column:district;type:varchar(50);comment:区县;default:NULL;" json:"district"`
	Street           string    `gorm:"column:street;type:varchar(100);comment:街道;default:NULL;" json:"street"`
	StreetNumber     string    `gorm:"column:street_number;type:varchar(100);comment:门牌号;default:NULL;" json:"street_number"`
	Citycode         string    `gorm:"column:citycode;type:varchar(20);comment:城市编码;default:NULL;" json:"citycode"`
	Status           int       `gorm:"column:status;type:int;comment:API返回状态;default:0;" json:"status"`
	Source           string    `gorm:"column:source;type:varchar(50);comment:数据来源;default:baidu;" json:"source"`
	HitCount         int       `gorm:"column:hit_count;type:int;comment:命中次数;default:1;" json:"hit_count"`
	CreatedAt        time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"updated_at"`
}

func (l *LxhReverseGeocodingCache) TableName() string {
	return "lxh_reverse_geocoding_cache"
}
