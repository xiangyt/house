package main

import (
	"compress/gzip"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	build, err := getBuildingTable("http://119.97.201.22:8083/spfxmcx/spfcx_fang.aspx?dengJH=%CF%C42000061&houseDengJh=%CF%C40008385")
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
			_, err := getHouseInfo(getHouseInfoUrl(house.GId), house)
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
	//fmt.Println(getHouseInfo(getHouseInfoUrl("1946d6e5-d328-4917-838f-5970e3709181"), &House{}))
}

//http://119.97.201.22:8080/TimeFL.aspx?gid=8d81fd6b-497e-4407-a64f-296f048414b4
const (
	NotSale  = "background-color:#CCFFFF" // 未销（预）售
	Sold     = "background-color:#FF0000" // 已网上销售
	Limit    = "background-color:#000000" // 限制出售
	Mortgage = "background-color:#FFFF00" // 已在建工程抵押
	Seized   = "background-color:#CC0099" // 已查封
	//MortgageSeized = "background-color:#FFFF00" // 已在建工程抵押已查封
)

func getSaleStatus(bgColor string) int {
	switch bgColor {
	case NotSale:
		return 1
	case Sold:
		return 2
	case Limit:
		return 3
	case Mortgage:
		return 4
	case Seized:
		return 5
	}
	return 0
}

func printSaleStatus(status int) {
	switch status {
	case 1:
		fmt.Println("未销（预）售")
	case 2:
		fmt.Println("已网上销售")
	case 3:
		fmt.Println("限制出售")
	case 4:
		fmt.Println("已在建工程抵押")
	case 5:
		fmt.Println("已查封")
	default:
		fmt.Println("未知")
	}
}

type House struct {
	GId              string
	BuildingNum      string // 栋号
	Unit             string // 单元
	Floor            string // 层数
	Room             string // 室号
	Description      string // 房屋坐落
	ConstructionArea string // 预售（现售）建筑面积（平方米）
	UnitPrice        string // 预售（现售）单价（元/平方米）
	TotalPrice       string // 房屋总价款（元）
	RoughcastPrice   string // 其中	毛坯价款（元）
	FurnishPrice     string // 装修价款（元）
	DeliveryStandard string // 交付标准
	Status           int    // 状态
}

func (h *House) print() {
	fmt.Printf("坐落:%s, 楼栋:%s, 单元:%s, 层数:%s, 总价:%s, 状态:%s,\r\n",
		h.Description, h.BuildingNum, h.Unit, h.Floor, h.TotalPrice, h.getStatus())
}

func (h *House) getStatus() string {
	switch h.Status {
	case 1:
		return "未销（预）售"
	case 2:
		return "已网上销售"
	case 3:
		return "限制出售"
	case 4:
		return "已在建工程抵押"
	case 5:
		return "已查封"
	default:
		return "未知"
	}
}

type Building struct {
	Number string // 栋号
	Houses []*House
}

// 获取楼盘表
func getBuildingTable(url string) (*Building, error) {
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

	var houses []*House
	doc.Find("#fwxx table tr").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		if i <= 1 {
			return
		}
		//fmt.Printf("Review tr %d\n", i)
		house := &House{}
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			switch i {
			case 0:
				house.BuildingNum, _ = s.Html()
			case 1:
				house.Unit, _ = s.Html()
			case 2:
				house.Floor, _ = s.Html()
			default:
				var h = *house
				hurl, _ := s.Find("a").Attr("href")
				arr := strings.Split(hurl, "?gid=")
				if len(arr) == 2 {
					h.GId = arr[1]
					h.Room = s.Text()
					bgColor, _ := s.Attr("style")
					h.Status = getSaleStatus(bgColor)
					houses = append(houses, &h)
				}
			}

		})

	})

	//for _, house := range houses {
	//	fmt.Printf("house: %+v\n", house)
	//}

	return &Building{Houses: houses}, nil
}

func getHouseInfoUrl(qid string) string {
	return "http://119.97.201.22:8080/TimeFL.aspx?gid=" + qid
}

func getHouseInfo(url string, house *House) (*House, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "119.97.201.22:8080")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "http://119.97.201.22:8083/")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36")

	res, err := (&http.Client{}).Do(req)
	//resp, err := http.Get(serviceUrl + "/topic/query/false/lsj")
	if err != nil {
		return house, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return house, errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}
	reader, err := gzip.NewReader(res.Body)
	if err != nil {
		return house, err
	}
	defer reader.Close()
	//b, _ := ioutil.ReadAll(reader)
	//fmt.Println(string(b))

	//	reader := bytes.NewReader([]byte(`
	//<!DOCTYPE html>
	//
	//<html xmlns="http://www.w3.org/1999/xhtml">
	//<head><meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	//    <!-- 最新版本的 Bootstrap 核心 CSS 文件 -->
	//    <link rel="stylesheet" href="https://cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap.min.css" />
	//
	//    <!-- 可选的 Bootstrap 主题文件（一般不用引入） -->
	//    <link rel="stylesheet" href="https://cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap-theme.min.css" /><title>
	//
	//</title></head>
	//<body>
	//    <form method="post" action="xsysjg.aspx?gid=1946d6e5-d328-4917-838f-5970e3709181" id="form1">
	//<div class="aspNetHidden">
	//<input type="hidden" name="__VIEWSTATE" id="__VIEWSTATE" value="/wEPDwUKLTkxODg0NDI3NQ9kFgICAw9kFhhmDxYCHgRUZXh0BT3mrabmsYnnvo7nmoTlkJvlhbDljYrlsptBLTPpobnnm64xMOagizLljZXlhYMxLTLlsYLvvIgy77yJ5Y+3ZAIBDxYCHwAFG+atpuaIv+W8gOmihOWUrlsyMDIxXTE4MeWPt2QCAg8WAh8ABQYxNDEuNjZkAgMPFgIfAAUIMjU0MzkuMDBkAgQPFgIfAAUKMzYwMzY4OC43NGQCBQ8WAh8ABQozNjAzNjg4Ljc0ZAIGDxYCHwAFATBkAgcPFgIfAAUJ5q+b5Z2v5oi/ZAIIDxYCHwAFPeatpuaxiee+jueahOWQm+WFsOWNiuWym0EtM+mhueebrjEw5qCLMuWNleWFgzEtMuWxgu+8iDLvvInlj7dkAgkPFgIfAAUb5q2m5oi/5byA6aKE5ZSuWzIwMjFdMTgx5Y+3ZAIKDxYCHwAFBjE0MS42NmQCCw8WAh8ABQnmr5vlna/miL9kZKW404I+GzO/nK9w0Fui5bAH/x6z/RekxAacASX2Hc5O" />
	//</div>
	//
	//<div class="aspNetHidden">
	//
	//        <input type="hidden" name="__VIEWSTATEGENERATOR" id="__VIEWSTATEGENERATOR" value="C58710BB" />
	//</div>
	//        <div class="container">
	//            <h1>商品房预售（现售）方案信息查询</h1>
	//            <table class="table table-condensed table-hover" style="">
	//                 <thead>
	//                        <tr>
	//                            <th>#</th>
	//                            <th>楼盘指标</th>
	//                            <th>详细信息</th>
	//                        </tr>
	//                      </thead>
	//                <tbody>
	//                    <tr><td>1</td><td>房屋座落：</td><td>武汉美的君兰半岛A-3项目10栋2单元1-2层（2）号</td></tr>
	//                    <tr><td>2</td><td>预售（现售）许可证号：</td><td>武房开预售[2021]181号</td></tr>
	//                    <tr><td>3</td><td>预售（现售）建筑面积（平方米）：</td><td>141.66</td></tr>
	//                    <tr><td>4</td><td>预售（现售）单价（元/平方米）：</td><td>25439.00</td></tr>
	//                    <tr><td>5</td><td>房屋总价款（元）：</td><td>3603688.74</td></tr>
	//                    <tr><td rowspan="2">其中</td><td> 毛坯价款（元）：</td><td>3603688.74</td></tr>
	//                    <tr><td> 装修价款（元）：</td><td>0</td></tr>
	//                    <tr><td>6</td><td>交付标准：</td><td>毛坯房</td></tr>
	//                </tbody>
	//            </table>
	//            <table class="table table-condensed table-hover" style="display:none">
	//                 <thead>
	//                        <tr>
	//                            <th>#</th>
	//                            <th>楼盘指标</th>
	//                            <th>详细信息</th>
	//                        </tr>
	//                      </thead>
	//                <tbody>
	//                    <tr><td>1</td><td>房屋座落：</td><td>武汉美的君兰半岛A-3项目10栋2单元1-2层（2）号</td></tr>
	//                    <tr><td>2</td><td>预售（现售）许可证号：</td><td>武房开预售[2021]181号</td></tr>
	//                    <tr><td>3</td><td>预售（现售）建筑面积（平方米）：</td><td>141.66</td></tr>
	//                    <tr><td>4</td><td>交付标准：</td><td>毛坯房</td></tr>
	//                    <tr><td></td><td>该房屋为企业自留房</td><td></td></tr>
	//                </tbody>
	//            </table>
	//            <p>备注：</p>
	//            <p>1.2016年1月1日后取得预售许可的商品房项目中，准售房屋的单价可在预售许可证发放后的2个工作日后在楼盘表中查询。</p>
	//            <p>2.因我市新建商品房买卖合同示范文本2020年5月25日启用，自2020年5月25日起取得预售许可的商品房项目信息按照新版本进行公示。</p>
	//            <p>3.因网络通讯故障、房地产开发企业上报信息不及时、上报信息填写错误等原因，可能会造成已备案的合同信息查询不到。请同具体的房地产开发企业及办理业务的区房产管理局联系。</p>
	//            <p>4.本网站数据维护时间为每天凌晨1:00--5:00。数据维护期间，可能出现页面无法响应的问题，请避开此时段进行查询。</p>
	//        </div>
	//    </form>
	//</body>
	//</html>
	//`))
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return house, err
	}

	doc.Find(".container table tr").Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(j int, s1 *goquery.Selection) {
			//fmt.Println(i, j, s1.Text())
			switch {
			case i == 1 && j == 2:
				house.Description = s1.Text()
			case i == 3 && j == 2:
				house.ConstructionArea = s1.Text()
			case i == 4 && j == 2:
				house.UnitPrice = s1.Text()
			case i == 5 && j == 2:
				house.TotalPrice = s1.Text()
			case i == 6 && j == 2:
				house.RoughcastPrice = s1.Text()
			case i == 7 && j == 1:
				house.FurnishPrice = s1.Text()
			case i == 8 && j == 2:
				house.DeliveryStandard = s1.Text()
			}
		})
	})

	return house, nil
}

//<div class="box" style="width:1120px;">
//<input name="hide" type="hidden" id="hide" value="夏2000061" />
//<div class="box_border">
//<table class="f_table">
//<tr>
//<td><a id="href1" href="#" onclick="getDengjh()">&nbsp;&nbsp;&nbsp;楼盘表（栋表）</a></td>
//<td style="color:red">&nbsp;&nbsp;&nbsp;<!--人房关联--></td>
//<td><span>&nbsp;&nbsp;&nbsp;楼盘表（房表）</span></td>
//</tr>
//</table>

//<div id="fwxx"><table class='tab_style'><th>栋号</th><th>单元</th><th>层数</th><th>室号</th><th>室号</th><th>室号</th><th>室号</th><tr><td>10</td><td>/</td><td>1</td><td style=background-color:#CCFFFF><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=00000000-0000-0000-0000-000000000000 target='_blank'>（3）风井</a></td><td style=background-color:#CCFFFF><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=00000000-0000-0000-0000-000000000000 target='_blank'>（4）风井</a></td><td style=background-color:#CCFFFF><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=00000000-0000-0000-0000-000000000000 target='_blank'>（5）风井</a></td><td style=background-color:#CCFFFF><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=00000000-0000-0000-0000-000000000000 target='_blank'>（6）风井</a></td><tr><td>10</td><td>1</td><td>1-2</td><td style=background-color:#FF0000><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=d40d7c44-5885-4d43-802e-b42c6d24d891 target='_blank'>（1）</a></td><td style=background-color:#CC0099><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=02bd3797-c093-4c0d-9a38-8265251a90c9 target='_blank'>（2）</a></td><tr><td>10</td><td>1</td><td>3-4</td><td style=background-color:#CCFFFF><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=84e495b6-070a-4daa-a64c-7fa8436dba5d target='_blank'>（1）</a></td><td style=background-color:#CCFFFF><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=9d23d1d8-1c32-4d4d-866b-52ac1e9238e3 target='_blank'>（2）</a></td><tr><td>10</td><td>1</td><td>5-6</td><td style=background-color:#FF0000><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=59c8f545-52bc-468f-8aec-b888bbb61bdf target='_blank'>（1）</a></td><td style=background-color:#FF0000><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=874957af-bd7e-45b4-839f-90ec1afadd31 target='_blank'>（2）</a></td><tr><td>10</td><td>2</td><td>1-2</td><td style=background-color:#FF0000><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=8d81fd6b-497e-4407-a64f-296f048414b4 target='_blank'>（1）</a></td><td style=background-color:#FF0000><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=23317fab-c352-4b9b-8063-806025085930 target='_blank'>（2）</a></td><tr><td>10</td><td>2</td><td>3-4</td><td style=background-color:#FF0000><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=a71e42bf-d5db-4e5f-a5d7-739b4daa7c2a target='_blank'>（1）</a></td><td style=background-color:#CCFFFF><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=d13ba7e5-0951-4f78-9274-5cf6d083ad59 target='_blank'>（2）</a></td><tr><td>10</td><td>2</td><td>5-6</td><td style=background-color:#FF0000><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=4042f70b-1019-492c-a716-e0f32d0249ca target='_blank'>（1）</a></td><td style=background-color:#FF0000><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=b6102955-9886-4b4d-9b3f-bd66f796e791 target='_blank'>（2）</a></td><tr><td>10</td><td>3</td><td>1-2</td><td style=background-color:#FF0000><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=0a4ece3b-1b8f-405b-99ee-b24f615bbc0b target='_blank'>（1）</a></td><td style=background-color:#FF0000><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=3931ce88-fed0-416b-a407-1f122e872775 target='_blank'>（2）</a></td><tr><td>10</td><td>3</td><td>3-4</td><td style=background-color:#CCFFFF><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=c2d6f910-b8f8-42f0-9dbd-9fbc33d1eda3 target='_blank'>（1）</a></td><td style=background-color:#CCFFFF><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=d31438b5-a3c9-4554-b5ee-8e4a20967c71 target='_blank'>（2）</a></td><tr><td>10</td><td>3</td><td>5-6</td><td style=background-color:#FF0000><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=c690eb0e-c825-484a-9b7e-6fd29687c58e target='_blank'>（1）</a></td><td style=background-color:#FF0000><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=888dc965-c2b2-459f-a2b4-61297ba61625 target='_blank'>（2）</a></td><tr><td>10</td><td>4</td><td>1-2</td><td style=background-color:#FF0000><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=f808139c-6340-47cc-bd35-f32401c6224e target='_blank'>（1）</a></td><td style=background-color:#FF0000><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=7e9292e1-c2fd-4a95-8c53-1ee600ef9e43 target='_blank'>（2）</a></td><tr><td>10</td><td>4</td><td>3-4</td><td style=background-color:#CCFFFF><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=1620c11a-6a52-48a7-9642-b75f68dc7dd1 target='_blank'>（1）</a></td><td style=background-color:#FF0000><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=a9ed6647-2e2f-46b6-aa1c-09db9451e11d target='_blank'>（2）</a></td><tr><td>10</td><td>4</td><td>5-6</td><td style=background-color:#FF0000><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=3b5eebbb-86e7-4f92-a31f-8b08f6f0ebc2 target='_blank'>（1）</a></td><td style=background-color:#FF0000><a href=http://119.97.201.22:8080/TimeFL.aspx?gid=e87b3a09-fc52-4f94-804c-dcd85386f7cf target='_blank'>（2）</a></td></table></div>
//</div>
//<div class="wxts" style="margin-top:30px;">
//<p>备注：1.2016年1月1日后取得预售许可的商品房项目中，准售房屋的单价可在预售许可证发放后的2个工作日后在楼盘表中查询。</p>
//<p style="text-indent:24px;">2.因我市新建商品房买卖合同示范文本2020年5月25日启用，自2020年5月25日起取得预售许可的商品房项目信息按照新版本进行公示。</p>
//<p style="text-indent:24px;">3.因网络通讯故障、房地产开发企业上报信息不及时、上报信息填写错误等原因，可能会造成已备案的合同信息查询不到。请同具体的房地产开发企业及办理业务的区房产管理局联系。</p>
//<p style="text-indent:24px;">4.本网站数据维护时间为每天凌晨1:00--5:00。数据维护期间，可能出现页面无法响应的问题，请避开此时段进行查询。</p>
//</div>
//</div>
//<div class="scroll scroll_1">
//<div class="scroll_child"></div>
//</div>
//</form>
//</body>
//</html>
