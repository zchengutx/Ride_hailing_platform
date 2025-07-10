package model

type Region struct {
	Code    int64  `gorm:"column:code;type:int;comment:行政区划代码;primaryKey;not null;" json:"code"`               // 行政区划代码
	Name    string `gorm:"column:name;type:varchar(40);comment:行政区划名称;default:NULL;" json:"name"`              // 行政区划名称
	Pcode   int64  `gorm:"column:pcode;type:int;comment:上级区划代码;default:NULL;" json:"pcode"`                    // 上级区划代码
	Sname   string `gorm:"column:sname;type:varchar(40);comment:地名简称;default:NULL;" json:"sname"`              // 地名简称
	Level   int64  `gorm:"column:level;type:int;comment:行政区划等级（1：省、直辖市；2：市州；3：区县）;default:NULL;" json:"level"` // 行政区划等级（1：省、直辖市；2：市州；3：区县）
	Mername string `gorm:"column:mername;type:varchar(100);comment:组合名称;default:NULL;" json:"mername"`         // 组合名称
	Pinyin  string `gorm:"column:pinyin;type:varchar(100);comment:拼音;default:NULL;" json:"pinyin"`             // 拼音
}

func (r *Region) TableName() string {
	return "region"
}
