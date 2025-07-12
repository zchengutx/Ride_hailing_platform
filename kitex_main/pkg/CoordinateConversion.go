package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type LocationResponse struct {
	Status int `json:"status"`
	Result struct {
		Location struct {
			Lng float64 `json:"lng"`
			Lat float64 `json:"lat"`
		} `json:"location"`
		FormattedAddress string `json:"formatted_address"`
		Edz              struct {
			Name string `json:"name"`
		} `json:"edz"`
		Business         string        `json:"business"`
		BusinessInfo     []interface{} `json:"business_info"` // 空数组类型
		AddressComponent struct {
			Country         string `json:"country"`
			CountryCode     int    `json:"country_code"`
			CountryCodeIso  string `json:"country_code_iso"`
			CountryCodeIso2 string `json:"country_code_iso2"`
			Province        string `json:"province"`
			City            string `json:"city"`
			CityLevel       int    `json:"city_level"`
			District        string `json:"district"`
			Town            string `json:"town"`
			TownCode        string `json:"town_code"`
			Distance        string `json:"distance"`
			Direction       string `json:"direction"`
			Adcode          string `json:"adcode"`
			Street          string `json:"street"`
			StreetNumber    string `json:"street_number"`
		} `json:"addressComponent"`
		Pois []struct {
			Addr      string `json:"addr"`
			Cp        string `json:"cp"`
			Direction string `json:"direction"`
			Distance  string `json:"distance"`
			Name      string `json:"name"`
			PoiType   string `json:"poiType"`
			Point     struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"point"`
			Tag             string `json:"tag"`
			Tel             string `json:"tel"`
			Uid             string `json:"uid"`
			Zip             string `json:"zip"`
			PopularityLevel string `json:"popularity_level"`
			AoiName         string `json:"aoi_name"`
			ParentPoi       struct {
				Name  string `json:"name"`
				Tag   string `json:"tag"`
				Addr  string `json:"addr"`
				Point struct {
					X float64 `json:"x"`
					Y float64 `json:"y"`
				} `json:"point"`
				Direction       string `json:"direction"`
				Distance        string `json:"distance"`
				Uid             string `json:"uid"`
				PopularityLevel string `json:"popularity_level"`
			} `json:"parent_poi"`
		} `json:"pois"`
		Roads      []interface{} `json:"roads"` // 空数组类型
		PoiRegions []struct {
			DirectionDesc string  `json:"direction_desc"`
			Name          string  `json:"name"`
			Tag           string  `json:"tag"`
			Uid           string  `json:"uid"`
			Distance      string  `json:"distance"`
			RegionArea    float64 `json:"region_area"`
		} `json:"poiRegions"`
		SematicDescription  string `json:"sematic_description"`
		FormattedAddressPoi string `json:"formatted_address_poi"`
		CityCode            int    `json:"cityCode"`
	} `json:"result"`
}

func CoordinateConversion(Coordinate []string) (StartAddr, EndEnd string) {
	var Coordinates []string

	for i := 0; i < len(Coordinate); i++ {

		// 此处填写您在控制台-应用管理-创建应用后获取的AK
		ak := "fbZPlNxNhl49M4bg6zHwLTP0DpJym7eS"

		// 服务地址
		host := "https://api.map.baidu.com"

		// 接口地址
		uri := "/reverse_geocoding/v3"

		// 设置请求参数
		params := url.Values{
			"ak":             []string{ak},
			"output":         []string{"json"},
			"coordtype":      []string{"bd09ll"},
			"extensions_poi": []string{"1"},
			"location":       []string{Coordinate[i]},
		}

		// 发起请求
		request, err := url.Parse(host + uri + "?" + params.Encode())
		if nil != err {
			fmt.Printf("host error: %v", err)
			return
		}

		resp, err1 := http.Get(request.String())
		fmt.Printf("url: %s\n", request.String())
		defer resp.Body.Close()
		if err1 != nil {
			fmt.Printf("request error: %v", err1)
			return
		}
		body, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			fmt.Printf("response error: %v", err2)
		}

		var Data LocationResponse

		json.Unmarshal(body, &Data)

		Coordinates = append(Coordinates, Data.Result.FormattedAddress)
	}

	StartAddr = Coordinates[0]
	EndEnd = Coordinates[1]

	return StartAddr, EndEnd
}
