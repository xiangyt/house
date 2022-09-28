package service

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xiangyt/house/constants"
	"github.com/xiangyt/house/db"
	"github.com/xiangyt/house/model"
	"github.com/xiangyt/house/util"
	"gorm.io/gorm"
	"sync"
	"time"
)

func RefreshZone() {

}

func RefreshBuilding() {
	logrus.Info("RefreshBuilding Start")
	var zones []*model.Zone
	if err := db.GetDB().Find(&zones).Error; err != nil {
		logrus.Errorf("get all zone failed. err:%s", err)
		return
	}

	for _, zone := range zones {
		buildings, err := util.GetLouPanTable(zone.Id)
		if err != nil {
			logrus.Errorf("get buildings by zone[%s] failed. err:%s", zone.Id, err)
			continue
		}

		for _, building := range buildings {
			db.GetDB().Save(building)
		}
	}
	logrus.Info("RefreshBuilding End")
}

func RefreshHouse() {
	logrus.Info("RefreshHouse Start")
	var buildings []*model.Building
	if err := db.GetDB().Find(&buildings).Error; err != nil {
		logrus.Errorf("get all building failed. err:%s", err)
		return
	}

	var wg sync.WaitGroup
	for _, building := range buildings {
		houses, err := util.GetFangTable(building.ZoneId, building.Id)
		if err != nil {
			logrus.Errorf("get houses by building[%s] failed. err:%s", building.Id, err)
			continue
		}

		for _, house := range houses {
			house := house
			wg.Add(1)
			go func() {
				defer wg.Done()
				var h model.House
				if house.GId == "00000000-0000-0000-0000-000000000000" {
					house.GId = fmt.Sprintf("%s-%s-%s-%s", house.BuildingNum, house.Unit, house.Floor, house.Room)
				}
				if err := db.GetDB().Where("gid = ?", house.GId).Find(&h).Error; err != nil {
					if !errors.Is(err, gorm.ErrRecordNotFound) {
						return
					}
					logrus.Debugf("RefreshHouse Find house failed. GId:%+v err:%s", house.GId, err)
				}

				if h.Status == constants.Sold {
					return
				}

				if h.GId == "" {
					if house, err = util.GetHouseInfo(house); err != nil {
						logrus.Debugf("RefreshHouse GetHouseInfo failed. house:%+v err:%s", house, err)
						return
					}
					if house.Status == constants.Sold && house.SaleTime == 0 {
						house.SaleTime = time.Now().Unix()
					}
					if err := db.GetDB().Save(house).Error; err != nil {
						logrus.Errorf("save house failed. house:%+v err:%s", house, err)
					}
				} else if h.Status != constants.Sold && house.Status == constants.Sold {
					h.SaleTime = time.Now().Unix()
					if err := db.GetDB().Save(house).Error; err != nil {
						logrus.Errorf("save house failed. house:%+v err:%s", house, err)
					}
				}
			}()
		}
	}
	wg.Wait()

	logrus.Info("RefreshHouse End")
}
