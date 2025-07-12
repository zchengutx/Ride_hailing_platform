package model

import (
	"kitex_main/config"
	"time"
)

type Trips struct {
	Id            int32     `gorm:"column:id;type:int;comment:主键id;primaryKey;" json:"id"`                    // 主键id
	CreatedAt     time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;" json:"created_at"`        // 创建时间
	UserId        int32     `gorm:"column:user_id;type:int;comment:用户id;" json:"user_id"`                     // 用户id
	StartingPoint string    `gorm:"column:starting_point;type:varchar(50);comment:起点;" json:"starting_point"` // 起点
	Terminal      string    `gorm:"column:terminal;type:varchar(50);comment:终点;" json:"terminal"`             // 终点
	EndTime       string    `gorm:"column:end_time;type:varchar(50);comment:结束时间;" json:"end_time"`         // 结束时间
}

func (t *Trips) TableName() string {
	return "trips"
}

// 添加行程
func (t *Trips) AddTrip() (err error) {
	err = config.DB.Create(t).Error
	return
}

// 查询行程是否存在
func (t *Trips) FindTripsId(id int) (err error) {
	err = config.DB.Where("id = ?", id).First(&t).Error
	return
}

// 修改行程
func (t *Trips) UpdateTrip(id int) (err error) {
	err = config.DB.Where("id = ?", id).Updates(t).Error
	return
}
