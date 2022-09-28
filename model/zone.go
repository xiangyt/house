package model

type Zone struct {
	Id            string      `gorm:"primaryKey;column:id"`
	Name          string      // 项目名称
	Position      string      // 项目坐落
	BuildingCount int64       // 房屋栋数
	HouseCount    int64       // 房屋套数
	Enterprise    string      // 开发企业
	PhoneNumber   string      // 联系电话
	Buildings     []*Building `gorm:"foreignKey:ZoneId"`
}

func (z Zone) TableName() string {
	return "zones"
}
