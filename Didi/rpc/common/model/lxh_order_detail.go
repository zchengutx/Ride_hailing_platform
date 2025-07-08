package model

// 订单扩展表
type LxhOrderDetail struct {
	Id        int32  `gorm:"column:id;type:int;primaryKey;not null;" json:"id"`
	OrderCode string `gorm:"column:order_code;type:varchar(50);comment:订单编号;default:NULL;" json:"order_code"`         // 订单编号
	TripKey   string `gorm:"column:trip_key;type:varchar(50);comment:mongdb中行程的key;default:NULL;" json:"trip_key"`    // mongdb中行程的key
	DriverKey string `gorm:"column:driver_key;type:varchar(50);comment:司机实际行程路线的key;default:NULL;" json:"driver_key"` // 司机实际行程路线的key
}

func (o *LxhOrderDetail) TableName() string {
	return "lxh_order_detail"
}
