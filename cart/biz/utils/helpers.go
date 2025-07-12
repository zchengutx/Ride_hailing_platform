package utils

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"

	driverpb "cart/kitex_gen/cart/driver"
	passengerpb "cart/kitex_gen/cart/passenger"
	"cart/rpc/basic/model"
)

// GetCurrentTimestamp 获取当前时间戳（毫秒）
func GetCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}

// GetEnv 获取环境变量，如果不存在则返回默认值
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ValidateDriverInfo 验证司机基本信息
func ValidateDriverInfo(req *driverpb.DriverPetitionReq) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("姓名不能为空")
	}
	if strings.TrimSpace(req.NickName) == "" {
		return fmt.Errorf("昵称不能为空")
	}
	if strings.TrimSpace(req.Mobile) == "" {
		return fmt.Errorf("手机号不能为空")
	}
	if len(req.Mobile) != 11 {
		return fmt.Errorf("手机号格式不正确")
	}
	if req.CarAge < 0 || req.CarAge > 20 {
		return fmt.Errorf("车龄应在0-20年之间")
	}
	if strings.TrimSpace(req.IdCardFileId) == "" {
		return fmt.Errorf("身份证照片不能为空")
	}
	if strings.TrimSpace(req.DriverLicenseFileId) == "" {
		return fmt.Errorf("驾驶证照片不能为空")
	}
	if strings.TrimSpace(req.DrivingLicenseFileId) == "" {
		return fmt.Errorf("行驶证照片不能为空")
	}
	if strings.TrimSpace(req.AvatarFileId) == "" {
		return fmt.Errorf("头像不能为空")
	}
	return nil
}

// GenerateOrderCode 生成订单编号
func GenerateOrderCode() string {
	timestamp := time.Now().Unix()
	randomNum := rand.Intn(10000)
	return fmt.Sprintf("D%d%04d", timestamp, randomNum)
}

// ConvertOrderToInfo 将订单模型转换为返回信息
func ConvertOrderToInfo(order model.LxhOrder, db *gorm.DB) *driverpb.OrderInfo {
	var passenger model.LxhPassenger
	passengerName := "未知乘客"
	passengerMobile := ""
	if err := db.Where("id = ?", order.PassengerId).First(&passenger).Error; err == nil {
		if passenger.NickName != "" {
			passengerName = passenger.NickName
		} else if passenger.Name != "" {
			passengerName = passenger.Name
		}
		if len(passenger.Mobile) == 11 {
			passengerMobile = passenger.Mobile[:3] + "****" + passenger.Mobile[7:]
		}
	}

	return &driverpb.OrderInfo{
		Id:              order.Id,
		OrderCode:       order.OrderCode,
		Amount:          order.Amount,
		OrderStatus:     order.OrderStatus,
		PassengerId:     order.PassengerId,
		PassengerName:   passengerName,
		PassengerMobile: passengerMobile,
		StartAddr:       order.StartAddr,
		EndAddr:         order.EndEnd,
		StartTime:       order.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:         order.EndTime.Format("2006-01-02 15:04:05"),
		OrderType:       order.OrderType,
	}
}

// GetRecentDestinations 获取用户最近目的地
func GetRecentDestinations(passengerId int16, db *gorm.DB) []string {
	var regions []model.Region
	if err := db.Where("level = 3").Limit(10).Find(&regions).Error; err != nil {
		return []string{"北京市朝阳区国贸中心", "上海市浦东新区陆家嘴", "广州市天河区珠江新城", "深圳市南山区科技园"}
	}
	var destinations []string
	for i, region := range regions {
		if i >= 4 {
			break
		}
		fullPath := GetRegionFullPath(region.Code, db)
		destinations = append(destinations, fullPath)
	}
	if len(destinations) == 0 {
		return []string{"北京市朝阳区国贸中心", "上海市浦东新区陆家嘴", "广州市天河区珠江新城", "深圳市南山区科技园"}
	}
	return destinations
}

// GetRegionFullPath 获取区域的完整路径
func GetRegionFullPath(regionCode int64, db *gorm.DB) string {
	var currentRegion model.Region
	if err := db.Where("code = ?", regionCode).First(&currentRegion).Error; err != nil {
		return "未知地区"
	}
	pathParts := []string{currentRegion.Name}
	if currentRegion.Pcode != 0 {
		var parentRegion model.Region
		if err := db.Where("code = ?", currentRegion.Pcode).First(&parentRegion).Error; err == nil {
			if parentRegion.Pcode != 0 {
				var grandParentRegion model.Region
				if err := db.Where("code = ?", parentRegion.Pcode).First(&grandParentRegion).Error; err == nil {
					pathParts = append([]string{grandParentRegion.Name, parentRegion.Name}, pathParts...)
				} else {
					pathParts = append([]string{parentRegion.Name}, pathParts...)
				}
			} else {
				pathParts = append([]string{parentRegion.Name}, pathParts...)
			}
		}
	}
	fullPath := ""
	for i, part := range pathParts {
		if i > 0 {
			fullPath += "/"
		}
		fullPath += part
	}
	return fullPath
}

// GetNearbyCars 获取附近车辆
func GetNearbyCars(location string, db *gorm.DB) []*passengerpb.CarInfo {
	var drivers []model.LxhDriver
	if err := db.Where("status = ?", "online").Limit(10).Find(&drivers).Error; err != nil {
		return []*passengerpb.CarInfo{}
	}
	if len(drivers) == 0 {
		return []*passengerpb.CarInfo{}
	}
	var nearbyCarInfos []*passengerpb.CarInfo
	carTypes := []string{"经济型", "舒适型", "豪华型", "商务型"}
	platePrefix := []string{"京A", "京B", "京C", "京D", "京E", "沪A", "沪B", "粤A", "粤B"}
	for i, driver := range drivers {
		distance := 0.3 + rand.Float64()*2.7
		estimatedTime := int32(math.Max(2, distance*2))
		carType := carTypes[rand.Intn(len(carTypes))]
		prefix := platePrefix[rand.Intn(len(platePrefix))]
		plateNumber := fmt.Sprintf("%s%05d", prefix, rand.Intn(100000))
		rating := 4.0 + rand.Float64()
		driverName := driver.NickName
		if driverName == "" {
			driverName = driver.Name
		}
		if driverName == "" {
			driverName = fmt.Sprintf("司机%d", driver.Id)
		}
		carInfo := &passengerpb.CarInfo{
			CarId:         int16(driver.Id),
			CarType:       carType,
			LicensePlate:  plateNumber,
			Distance:      distance,
			EstimatedTime: int16(estimatedTime),
			DriverName:    driverName,
			Rating:        rating,
		}
		nearbyCarInfos = append(nearbyCarInfos, carInfo)
		if i >= 4 {
			break
		}
	}
	sort.Slice(nearbyCarInfos, func(i, j int) bool {
		return nearbyCarInfos[i].Distance < nearbyCarInfos[j].Distance
	})
	return nearbyCarInfos
}

// GetServiceInfo 获取服务信息
func GetServiceInfo() []*passengerpb.ServiceInfo {
	return []*passengerpb.ServiceInfo{
		{ServiceName: "快车", ServiceDesc: "经济实惠，快速到达", ServiceIcon: "icon_kuaiche", ServiceUrl: "/service/kuaiche"},
		{ServiceName: "专车", ServiceDesc: "舒适体验，专业服务", ServiceIcon: "icon_zhuanche", ServiceUrl: "/service/zhuanche"},
		{ServiceName: "豪华车", ServiceDesc: "高端体验，尊贵享受", ServiceIcon: "icon_haohua", ServiceUrl: "/service/haohua"},
		{ServiceName: "代驾", ServiceDesc: "安全代驾，放心回家", ServiceIcon: "icon_daijia", ServiceUrl: "/service/daijia"},
	}
}

// GetCurrentLocationFromRegion 根据位置字符串获取完整路径
func GetCurrentLocationFromRegion(locationStr string, db *gorm.DB, defaultLocation string) string {
	if locationStr == "" {
		return defaultLocation
	}
	var regions []model.Region
	if err := db.Where("name LIKE ? OR sname LIKE ? OR mername LIKE ?", "%"+locationStr+"%", "%"+locationStr+"%", "%"+locationStr+"%").Limit(1).Find(&regions).Error; err != nil || len(regions) == 0 {
		return locationStr
	}
	return GetRegionFullPath(regions[0].Code, db)
}
