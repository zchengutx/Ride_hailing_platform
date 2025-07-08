package model

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	Id               uint64         `gorm:"column:id;type:int UNSIGNED;primaryKey;not null;" json:"id"`
	UserId           uint64         `gorm:"column:user_id;type:int UNSIGNED;comment:用户;not null;" json:"user_id"`                           // 用户
	DriverId         uint64         `gorm:"column:driver_id;type:int UNSIGNED;comment:司机;not null;" json:"driver_id"`                       // 司机
	StartLocation    string         `gorm:"column:start_location;type:varchar(50);comment:起始地;not null;" json:"start_location"`             // 起始地
	EndLocation      string         `gorm:"column:end_location;type:varchar(50);comment:目的地;not null;" json:"end_location"`                 // 目的地
	StartLng         string         `gorm:"column:start_lng;type:varchar(100);comment:起点经度;default:NULL;" json:"start_lng"`                 // 起点经度
	StartLat         string         `gorm:"column:start_lat;type:varchar(100);comment:起点纬度;default:NULL;" json:"start_lat"`                 // 起点纬度
	EndLng           string         `gorm:"column:end_lng;type:varchar(100);comment:目的地经度;default:NULL;" json:"end_lng"`                    // 目的地经度
	EndLat           string         `gorm:"column:end_lat;type:varchar(100);comment:目的地纬度;default:NULL;" json:"end_lat"`                    // 目的地纬度
	CartType         string         `gorm:"column:cart_type;type:varchar(10);comment:服务类型;default:快车;" json:"cart_type"`                    // 服务类型
	Status           string         `gorm:"column:status;type:varchar(10);comment:订单状态;default:待接单;" json:"status"`                         // 订单状态
	Price            float64        `gorm:"column:price;type:decimal(10, 2);comment:订单金额;default:NULL;" json:"price"`                       // 订单金额
	FavourableStatus string         `gorm:"column:favourable_status;type:varchar(10);comment:优惠状态;default:无;" json:"favourable_status"`     // 优惠状态
	FavourablePrice  float64        `gorm:"column:favourable_price;type:decimal(10, 2);comment:优惠金额;default:NULL;" json:"favourable_price"` // 优惠金额
	RealityPrice     float64        `gorm:"column:reality_price;type:decimal(10, 2);comment:实际金额;default:NULL;" json:"reality_price"`       // 实际金额
	PayStatus        string         `gorm:"column:pay_status;type:varchar(10);comment:支付状态;default:未支付;" json:"pay_status"`                 // 支付状态
	PayType          string         `gorm:"column:pay_type;type:varchar(10);comment:支付类型;default:微信;" json:"pay_type"`                      // 支付类型
	CreatedAt        time.Time      `gorm:"column:created_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;type:date;default:NULL;" json:"deleted_at"`
}

func (o *Order) TableName() string {
	return "order"
}
