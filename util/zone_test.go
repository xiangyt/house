package util

import (
	"fmt"
	"testing"
)

func TestGetZone(t *testing.T) {

	zone, err := GetZone("夏2000061")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(*zone)
}
