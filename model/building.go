package model

type Building struct {
	Id         string `gorm:"primaryKey,autoIncrement"`
	ZoneId     string
	Name       string   // 楼栋名称
	Floor      int32    // 总层数
	HouseCount int32    // 总套数
	Area       float64  // 建筑面积(㎡)
	Mapping    string   // 测绘机构
	Houses     []*House `gorm:"foreignKey:BuildingId"`
}

func (b Building) TableName() string {
	return "buildings"
}
