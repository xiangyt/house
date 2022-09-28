package util

import (
	"fmt"
	"testing"
)

func TestGetLouPanTable(t *testing.T) {
	// http://119.97.201.22:8083/spfxmcx/spfcx_lpb.aspx?DengJh=%CF%C42000061
	builds, err := GetLouPanTable("夏2000061")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, build := range builds {
		fmt.Printf("%+v\r\n", build)
	}
}

func TestGetFangTable(t *testing.T) {

	build, err := GetFangTable("夏2000061", "夏0008378")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(build)
}
