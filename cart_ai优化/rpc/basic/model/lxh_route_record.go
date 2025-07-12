package model

import "time"

// 路径记录表
type LxhRouteRecord struct {
	Id            int64     `gorm:"column:id;type:int;primaryKey;not null;" json:"id"`
	OrderId       int64     `gorm:"column:order_id;type:int;index;comment:关联订单ID;default:NULL;" json:"order_id"`
	PassengerId   int64     `gorm:"column:passenger_id;type:int;index;comment:乘客ID;default:NULL;" json:"passenger_id"`
	DriverId      int64     `gorm:"column:driver_id;type:int;index;comment:司机ID;default:NULL;" json:"driver_id"`
	StartAddress  string    `gorm:"column:start_address;type:varchar(500);comment:起点地址;default:NULL;" json:"start_address"`
	StartLng      float64   `gorm:"column:start_lng;type:decimal(11,8);comment:起点经度;default:NULL;" json:"start_lng"`
	StartLat      float64   `gorm:"column:start_lat;type:decimal(11,8);comment:起点纬度;default:NULL;" json:"start_lat"`
	EndAddress    string    `gorm:"column:end_address;type:varchar(500);comment:终点地址;default:NULL;" json:"end_address"`
	EndLng        float64   `gorm:"column:end_lng;type:decimal(11,8);comment:终点经度;default:NULL;" json:"end_lng"`
	EndLat        float64   `gorm:"column:end_lat;type:decimal(11,8);comment:终点纬度;default:NULL;" json:"end_lat"`
	Distance      float64   `gorm:"column:distance;type:decimal(10,2);comment:距离(公里);default:0;" json:"distance"`
	EstimatedTime int       `gorm:"column:estimated_time;type:int;comment:预估时间(分钟);default:0;" json:"estimated_time"`
	ActualTime    int       `gorm:"column:actual_time;type:int;comment:实际用时(分钟);default:0;" json:"actual_time"`
	RoutePoints   string    `gorm:"column:route_points;type:longtext;comment:路径坐标点(JSON);default:NULL;" json:"route_points"`
	RouteStatus   string    `gorm:"column:route_status;type:varchar(20);comment:路径状态(planning/running/completed/cancelled);default:planning;" json:"route_status"`
	StartTime     time.Time `gorm:"column:start_time;type:datetime;comment:开始时间;default:NULL;" json:"start_time"`
	EndTime       time.Time `gorm:"column:end_time;type:datetime;comment:结束时间;default:NULL;" json:"end_time"`
	CreatedAt     time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"updated_at"`
}

func (l *LxhRouteRecord) TableName() string {
	return "lxh_route_record"
}
