package model

import "time"

// 订单表
type LxhOrder struct {
	Id            int64     `gorm:"column:id;type:int;primaryKey;not null;" json:"id"`
	OrderCode     string    `gorm:"column:order_code;type:varchar(50);comment:订单编号;default:NULL;" json:"order_code"`            // 订单编号
	Amount        float64   `gorm:"column:amount;type:decimal(10, 2);comment:总价;default:NULL;" json:"amount"`                   // 总价
	OrderStatus   string    `gorm:"column:order_status;type:varchar(10);comment:订单状态;default:NULL;" json:"order_status"`        // 订单状态
	PassengerId   int64     `gorm:"column:passenger_id;type:int;comment:乘客ID;default:NULL;" json:"passenger_id"`                // 乘客ID
	StartAddr     string    `gorm:"column:start_addr;type:varchar(50);comment:起始地;default:NULL;" json:"start_addr"`             // 起始地
	EndEnd        string    `gorm:"column:end_end;type:varchar(50);comment:目的地;default:NULL;" json:"end_end"`                   // 目的地
	Driver        int64     `gorm:"column:driver;type:int;comment:司机ID;default:NULL;" json:"driver"`                            // 司机ID
	StartTime     time.Time `gorm:"column:start_time;type:datetime;comment:开始时间;default:NULL;" json:"start_time"`               // 开始时间
	EndTime       time.Time `gorm:"column:end_time;type:datetime;comment:结束时间;default:NULL;" json:"end_time"`                   // 结束时间
	ConfirmBy     string    `gorm:"column:confirm_by;type:varchar(5);comment:取消方（乘客/司机）;default:NULL;" json:"confirm_by"`       // 取消方（乘客/司机）
	ConfirmPerson int64     `gorm:"column:confirm_person;type:int;comment:取消人;default:NULL;" json:"confirm_person"`             // 取消人
	ConfirmReason string    `gorm:"column:confirm_reason;type:varchar(50);comment:取消原因;default:NULL;" json:"confirm_reason"`    // 取消原因
	ConfirmRemark string    `gorm:"column:confirm_remark;type:varchar(50);comment:取消备注;default:NULL;" json:"confirm_remark"`    // 取消备注
	PayStatus     string    `gorm:"column:pay_status;type:varchar(20);comment:支付状态;default:NULL;" json:"pay_status"`            // 支付状态
	PayType       string    `gorm:"column:pay_type;type:varchar(20);comment:支付方式;default:NULL;" json:"pay_type"`                // 支付方式
	OrderType     string    `gorm:"column:order_type;type:varchar(20);comment:订单状态（快车订单/预约订单）;default:NULL;" json:"order_type"` // 订单状态（快车订单/预约订单）
}

func (o *LxhOrder) TableName() string {
	return "lxh_order"
}
