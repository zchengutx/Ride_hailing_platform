package model

import "time"

// 用户收藏地点表
type LxhFavoriteLocation struct {
	Id           int64     `gorm:"column:id;type:int;primaryKey;not null;" json:"id"`
	PassengerId  int64     `gorm:"column:passenger_id;type:int;index;not null;comment:乘客ID;" json:"passenger_id"`
	LocationType string    `gorm:"column:location_type;type:varchar(20);comment:地点类型(home/company/custom);default:custom;" json:"location_type"`
	LocationName string    `gorm:"column:location_name;type:varchar(100);comment:地点名称;default:NULL;" json:"location_name"`
	Address      string    `gorm:"column:address;type:varchar(500);comment:详细地址;default:NULL;" json:"address"`
	Lng          float64   `gorm:"column:lng;type:decimal(11,8);comment:经度;default:NULL;" json:"lng"`
	Lat          float64   `gorm:"column:lat;type:decimal(11,8);comment:纬度;default:NULL;" json:"lat"`
	Province     string    `gorm:"column:province;type:varchar(50);comment:省份;default:NULL;" json:"province"`
	City         string    `gorm:"column:city;type:varchar(50);comment:城市;default:NULL;" json:"city"`
	District     string    `gorm:"column:district;type:varchar(50);comment:区县;default:NULL;" json:"district"`
	UsageCount   int       `gorm:"column:usage_count;type:int;comment:使用次数;default:0;" json:"usage_count"`
	IsDefault    int       `gorm:"column:is_default;type:tinyint;comment:是否默认地址(0否1是);default:0;" json:"is_default"`
	SortOrder    int       `gorm:"column:sort_order;type:int;comment:排序号;default:0;" json:"sort_order"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"updated_at"`
}

func (l *LxhFavoriteLocation) TableName() string {
	return "lxh_favorite_location"
}
