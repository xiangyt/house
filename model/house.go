package model

import "time"

type House struct {
	GId              string  `gorm:"primaryKey"`
	BuildingNum      string  // 栋号
	Unit             string  // 单元
	Floor            string  // 层数
	Room             string  // 室号
	Description      string  // 房屋坐落
	ConstructionArea float64 // 预售（现售）建筑面积（平方米）
	UnitPrice        float64 // 预售（现售）单价（元/平方米）
	TotalPrice       float64 // 房屋总价款（元）
	RoughcastPrice   float64 // 其中	毛坯价款（元）
	FurnishPrice     float64 // 装修价款（元）
	DeliveryStandard string  // 交付标准
	Status           int     // 状态
	SaleTime         time.Time
}

func (h House) TableName() string {
	return "houses"
}
