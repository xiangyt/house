package util

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/xiangyt/house/model"
	"net/http"
	"strconv"
	"strings"
)

// http://119.97.201.22:8083/spfxmcx/spfcx_lpb.aspx?DengJh=%CF%C42000061

// GetLouPanTable 获取楼盘表(栋表)
func GetLouPanTable(zoneId string) ([]*model.Building, error) {
	url := fmt.Sprintf("http://%s:8083/spfxmcx/spfcx_lpb.aspx?DengJh=%s", IP, mahonia.NewEncoder("GBK").ConvertString(zoneId))

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	decoder := mahonia.NewDecoder("GBK")

	//content, _ := doc.Html()
	//fmt.Println(decoder.ConvertString(content))

	var buildings []*model.Building
	doc.Find(".box table:nth-child(2) tbody tr").Each(func(i int, s *goquery.Selection) {
		//fmt.Printf("Review tr %d\n", i)
		building := &model.Building{
			ZoneId: zoneId,
		}
		//content, _ := s.Html()
		//fmt.Println(decoder.ConvertString(content))
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			switch i {
			case 0:
				building.Name = decoder.ConvertString(s.Text())
				href, _ := s.Find("a").Attr("href")
				arr := strings.Split(decoder.ConvertString(href), "houseDengJh=")
				if len(arr) == 2 {
					building.Id = arr[1]
				}
			case 1:
				num, _ := strconv.Atoi(s.Text())
				building.Floor = int32(num)
			case 2:
				num, _ := strconv.Atoi(s.Text())
				building.HouseCount = int32(num)
			case 3:
				building.Area, _ = strconv.ParseFloat(s.Text(), 64)
			case 4:
				building.Mapping = decoder.ConvertString(s.Text())
			}

		})
		if building.Id != "" {
			buildings = append(buildings, building)
		}
	})

	//for _, house := range houses {
	//	fmt.Printf("house: %+v\n", house)
	//}

	return buildings, nil
}
