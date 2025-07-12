package model

import "kitex_main/config"

type LxhPricingRules struct {
	Id         int32  `gorm:"column:id;type:int;comment:主键id;primaryKey;" json:"id"`             // 主键id
	CarType    string `gorm:"column:car_type;type:varchar(20);comment:车型;" json:"car_type"`      // 车型
	Kilometers string `gorm:"column:kilometers;type:varchar(20);comment:公里数;" json:"kilometers"` // 公里数
	Price      string `gorm:"column:price;type:varchar(20);comment:价格;" json:"price"`            // 价格
	Excess     string `gorm:"column:excess;type:varchar(20);comment:超额计价;" json:"excess"`        // 超额计价
}

func (l *LxhPricingRules) TableName() string {
	return "lxh_pricing_rules"
}

func (l *LxhPricingRules) FindLxhPricingRulesCarType(CarType string) (err error) {
	err = config.DB.Where("car_type=?", CarType).First(l).Error
	return
}
