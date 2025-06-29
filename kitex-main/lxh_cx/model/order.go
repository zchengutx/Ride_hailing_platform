package model

import "time"

type LxhOrder struct {
	Id            int32     `gorm:"column:id;type:int;primaryKey;" json:"id"`
	OrderCode     string    `gorm:"column:order_code;type:varchar(50);comment:订单编号;" json:"order_code"`            // 订单编号
	Amount        float64   `gorm:"column:amount;type:decimal(10, 2);comment:总价;" json:"amount"`                   // 总价
	OrderStatus   string    `gorm:"column:order_status;type:varchar(10);comment:订单状态;" json:"order_status"`        // 订单状态
	PassengerId   int32     `gorm:"column:passenger_id;type:int;comment:乘客ID;" json:"passenger_id"`                // 乘客ID
	StartAddr     string    `gorm:"column:start_addr;type:varchar(50);comment:起始地;" json:"start_addr"`             // 起始地
	EndEnd        string    `gorm:"column:end_end;type:varchar(50);comment:目的地;" json:"end_end"`                   // 目的地
	Driver        int32     `gorm:"column:driver;type:int;comment:司机ID;" json:"driver"`                            // 司机ID
	StartTime     time.Time `gorm:"column:start_time;type:datetime;comment:开始时间;" json:"start_time"`               // 开始时间
	EndTime       time.Time `gorm:"column:end_time;type:datetime;comment:结束时间;" json:"end_time"`                   // 结束时间
	ConfirmBy     string    `gorm:"column:confirm_by;type:varchar(5);comment:取消方（乘客/司机）;" json:"confirm_by"`       // 取消方（乘客/司机）
	ConfirmPerson int32     `gorm:"column:confirm_person;type:int;comment:取消人;" json:"confirm_person"`             // 取消人
	ConfirmReason string    `gorm:"column:confirm_reason;type:varchar(50);comment:取消原因;" json:"confirm_reason"`    // 取消原因
	ConfirmRemark string    `gorm:"column:confirm_remark;type:varchar(50);comment:取消备注;" json:"confirm_remark"`    // 取消备注
	PayStatus     string    `gorm:"column:pay_status;type:varchar(20);comment:支付状态;" json:"pay_status"`            // 支付状态
	PayType       string    `gorm:"column:pay_type;type:varchar(20);comment:支付方式;" json:"pay_type"`                // 支付方式
	OrderType     string    `gorm:"column:order_type;type:varchar(20);comment:订单状态（快车订单/预约订单）;" json:"order_type"` // 订单状态（快车订单/预约订单）
}
type LxhOrderDetail struct {
	Id        int32  `gorm:"column:id;type:int;primaryKey;" json:"id"`
	OrderCode string `gorm:"column:order_code;type:varchar(50);comment:订单编号;" json:"order_code"`         // 订单编号
	TripKey   string `gorm:"column:trip_key;type:varchar(50);comment:mongdb中行程的key;" json:"trip_key"`    // mongdb中行程的key
	DriverKey string `gorm:"column:driver_key;type:varchar(50);comment:司机实际行程路线的key;" json:"driver_key"` // 司机实际行程路线的key
}

func (o *LxhOrder) TableName() string {
	return "lxh_order"
}
func (d *LxhOrderDetail) TableName() string {
	return "lxh_order_detail"
}
