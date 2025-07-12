package model

import "time"

// MongoDB订单结构体
// 字段与MySQL订单结构体基本一致，便于后续迁移
// bson标签用于MongoDB文档映射

type MongoOrder struct {
	OrderId     string    `bson:"order_id"` // 订单唯一ID（建议用uuid或雪花id）
	PassengerId int64     `bson:"passenger_id"`
	StartAddr   string    `bson:"start_addr"`
	EndAddr     string    `bson:"end_addr"`
	DriverId    int64     `bson:"driver_id,omitempty"` // 司机ID，可为空
	OrderStatus string    `bson:"order_status"`        // 订单状态
	OrderType   string    `bson:"order_type"`          // 订单类型
	PayStatus   string    `bson:"pay_status"`          // 支付状态
	Amount      float64   `bson:"amount"`              // 金额
	StartTime   time.Time `bson:"start_time"`
	EndTime     time.Time `bson:"end_time,omitempty"`
	CreateTime  time.Time `bson:"create_time"`
	UpdateTime  time.Time `bson:"update_time"`
	// 可根据业务需求扩展更多字段
}
