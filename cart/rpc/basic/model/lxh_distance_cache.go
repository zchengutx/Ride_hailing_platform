package model

import "time"

// 距离计算缓存表
type LxhDistanceCache struct {
	Id            int64     `gorm:"column:id;type:int;primaryKey;not null;" json:"id"`
	RouteKey      string    `gorm:"column:route_key;type:varchar(200);uniqueIndex;not null;comment:路线键值;" json:"route_key"`
	OriginLng     float64   `gorm:"column:origin_lng;type:decimal(11,8);comment:起点经度;default:NULL;" json:"origin_lng"`
	OriginLat     float64   `gorm:"column:origin_lat;type:decimal(11,8);comment:起点纬度;default:NULL;" json:"origin_lat"`
	DestLng       float64   `gorm:"column:dest_lng;type:decimal(11,8);comment:终点经度;default:NULL;" json:"dest_lng"`
	DestLat       float64   `gorm:"column:dest_lat;type:decimal(11,8);comment:终点纬度;default:NULL;" json:"dest_lat"`
	Distance      float64   `gorm:"column:distance;type:decimal(10,2);comment:距离(公里);default:0;" json:"distance"`
	DistanceM     float64   `gorm:"column:distance_m;type:decimal(12,2);comment:距离(米);default:0;" json:"distance_m"`
	EstimatedTime int       `gorm:"column:estimated_time;type:int;comment:预估时间(分钟);default:0;" json:"estimated_time"`
	CalculateType string    `gorm:"column:calculate_type;type:varchar(20);comment:计算方式(straight/driving);default:straight;" json:"calculate_type"`
	HitCount      int       `gorm:"column:hit_count;type:int;comment:命中次数;default:1;" json:"hit_count"`
	CreatedAt     time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"updated_at"`
}

func (l *LxhDistanceCache) TableName() string {
	return "lxh_distance_cache"
}
