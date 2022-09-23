package util

import (
	"fmt"
	"sync"
	"testing"
)

func TestGetHouseInfo(t *testing.T) {
	fmt.Println(GetHouseInfo(GetHouseInfoUrl("1946d6e5-d328-4917-838f-5970e3709181"), &House{}))
}

func TestGetBuildingTable(t *testing.T) {
	build, err := GetBuildingTable("http://119.97.201.22:8083/spfxmcx/spfcx_fang.aspx?dengJH=%CF%C42000061&houseDengJh=%CF%C40008385")
	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup
	for _, house := range build.Houses {
		house := house
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := GetHouseInfo(GetHouseInfoUrl(house.GId), house)
			if err != nil {
				fmt.Println(err)
				return
			}
		}()
	}
	wg.Wait()

	for _, house := range build.Houses {
		house.print()
	}
}
