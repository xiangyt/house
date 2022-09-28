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

// GetZone 获取楼盘表(栋表)
func GetZone(zoneId string) (*model.Zone, error) {

	url := fmt.Sprintf("http://%s:8083/spfxmcx/spfcx_mx.aspx?DengJh=%s", IP, mahonia.NewEncoder("GBK").ConvertString(zoneId))

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

	var zone = model.Zone{
		Id: zoneId,
	}
	doc.Find("#table_mx tbody tr td").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 2:
			zone.Name = strings.TrimSpace(decoder.ConvertString(s.Text()))
		case 4:
			zone.Position = strings.TrimSpace(decoder.ConvertString(s.Text()))
		case 15:
			count := strings.TrimSpace(decoder.ConvertString(s.Text()))
			zone.BuildingCount, _ = strconv.ParseInt(strings.Replace(count, "栋", "", 1), 10, 64)
		case 17:
			count := strings.TrimSpace(decoder.ConvertString(s.Text()))
			zone.HouseCount, _ = strconv.ParseInt(strings.Replace(count, "套", "", 1), 10, 64)
		case 19:
			zone.Enterprise = strings.TrimSpace(decoder.ConvertString(s.Text()))
		case 21:
			zone.PhoneNumber = strings.TrimSpace(decoder.ConvertString(s.Text()))
		}
	})

	return &zone, nil
}

//项目名称	武汉美的君兰半岛A-3项目
//项目坐落	大桥新区办事处黄家湖大道
//项目
//基本情况	不动产权证书（国有土地使用证）证号	鄂（2019）武汉市江夏不动产权第0036299号	用地面积	33611.53㎡
//建设工程规划许可证号	武规（夏）建[2019]114号	建筑面积	详见栋列表
//房屋栋数	14栋	房屋套数	494套
//开发企业	武汉市鼎辉房地产开发有限公司
//联系电话	——
//合同备案办理部门	江夏区房管局
