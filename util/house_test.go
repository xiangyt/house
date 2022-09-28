package util

import (
	"fmt"
	"sync"
	"testing"
)

func TestGetHouseInfo(t *testing.T) {
	//fmt.Println(GetHouseInfo(GetHouseInfoUrl("1946d6e5-d328-4917-838f-5970e3709181"), &House{}))

	houses, err := GetFangTable("夏2000061", "夏0008385")
	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup
	for _, house := range houses {
		house := house
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := GetHouseInfo(house)
			if err != nil {
				fmt.Println(err)
				return
			}
		}()
	}
	wg.Wait()

	for _, house := range houses {
		fmt.Printf("%+v\r\n", house)
	}
}
