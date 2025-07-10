package model

import "time"

// 地理编码缓存表
type LxhGeocodingCache struct {
	Id        int64     `gorm:"column:id;type:int;primaryKey;not null;" json:"id"`
	Address   string    `gorm:"column:address;type:varchar(500);uniqueIndex;not null;comment:详细地址;" json:"address"`
	Lng       float64   `gorm:"column:lng;type:decimal(11,8);comment:经度;default:NULL;" json:"lng"`
	Lat       float64   `gorm:"column:lat;type:decimal(11,8);comment:纬度;default:NULL;" json:"lat"`
	Precise   int       `gorm:"column:precise;type:int;comment:精确度;default:0;" json:"precise"`
	Level     string    `gorm:"column:level;type:varchar(50);comment:地址级别;default:NULL;" json:"level"`
	Status    int       `gorm:"column:status;type:int;comment:API返回状态;default:0;" json:"status"`
	Source    string    `gorm:"column:source;type:varchar(50);comment:数据来源;default:baidu;" json:"source"`
	HitCount  int       `gorm:"column:hit_count;type:int;comment:命中次数;default:1;" json:"hit_count"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"updated_at"`
}

func (l *LxhGeocodingCache) TableName() string {
	return "lxh_geocoding_cache"
}
