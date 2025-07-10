package model

import "time"

// 热门地点表
type LxhHotLocation struct {
	Id           int64     `gorm:"column:id;type:int;primaryKey;not null;" json:"id"`
	LocationName string    `gorm:"column:location_name;type:varchar(200);comment:地点名称;default:NULL;" json:"location_name"`
	Address      string    `gorm:"column:address;type:varchar(500);comment:详细地址;default:NULL;" json:"address"`
	Lng          float64   `gorm:"column:lng;type:decimal(11,8);comment:经度;default:NULL;" json:"lng"`
	Lat          float64   `gorm:"column:lat;type:decimal(11,8);comment:纬度;default:NULL;" json:"lat"`
	Province     string    `gorm:"column:province;type:varchar(50);comment:省份;default:NULL;" json:"province"`
	City         string    `gorm:"column:city;type:varchar(50);index;comment:城市;default:NULL;" json:"city"`
	District     string    `gorm:"column:district;type:varchar(50);comment:区县;default:NULL;" json:"district"`
	Category     string    `gorm:"column:category;type:varchar(50);comment:地点分类(商场/医院/学校/酒店等);default:NULL;" json:"category"`
	SearchCount  int       `gorm:"column:search_count;type:int;comment:搜索次数;default:0;" json:"search_count"`
	ClickCount   int       `gorm:"column:click_count;type:int;comment:点击次数;default:0;" json:"click_count"`
	OrderCount   int       `gorm:"column:order_count;type:int;comment:下单次数;default:0;" json:"order_count"`
	HotScore     float64   `gorm:"column:hot_score;type:decimal(8,2);comment:热度得分;default:0;" json:"hot_score"`
	Status       int       `gorm:"column:status;type:tinyint;comment:状态(0禁用1启用);default:1;" json:"status"`
	StatDate     string    `gorm:"column:stat_date;type:date;index;comment:统计日期;default:NULL;" json:"stat_date"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"updated_at"`
}

func (l *LxhHotLocation) TableName() string {
	return "lxh_hot_location"
}
