package main

import (
	pb "cart/kitex_gen/cart/passenger"
	"cart/rpc/basic/global"
	"cart/rpc/basic/model"
	"context"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

// PassengerServiceImpl implements the last service interface defined in the IDL.
type PassengerServiceImpl struct{}

// SendSms implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) SendSms(ctx context.Context, req *pb.SendSmsReq) (*pb.SendSmsResp, error) {
	ctx = context.Background()

	// 从配置文件获取验证码配置
	smsConfig := &global.AppConf.SMSConfig

	// 生成验证码
	codeRange := smsConfig.MaxCodeValue - smsConfig.MinCodeValue
	code := rand.Intn(codeRange) + smsConfig.MinCodeValue

	// 设置验证码过期时间
	expireTime := time.Duration(smsConfig.CodeExpireTime) * time.Minute

	if err := global.Rdb.Set(ctx, "sendSms"+req.Mobile, code, expireTime).Err(); err != nil {
		return &pb.SendSmsResp{
			Code:    503,
			Message: "服务器异常",
		}, nil
	}
	return &pb.SendSmsResp{
		Code:    200,
		Message: "发送成功",
	}, nil
}

// RegisterPassenger implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) RegisterPassenger(ctx context.Context, req *pb.RegisterPassengerReq) (*pb.RegisterPassengerResp, error) {
	ctx = context.Background()
	result, err := global.Rdb.Get(ctx, "sendSms"+req.Mobile).Result()
	if err != nil {
		return &pb.RegisterPassengerResp{
			Code:    601,
			Message: "验证码获取失败",
		}, nil
	}
	var passenger model.LxhPassenger
	if err = global.DB.Debug().Where("mobile = ?", req.Mobile).Limit(1).Find(&passenger).Error; err != nil {
		return &pb.RegisterPassengerResp{
			Code:    503,
			Message: "服务器异常",
		}, nil
	}
	if passenger.Id != 0 {
		return &pb.RegisterPassengerResp{
			Code:    603,
			Message: "用户已存在",
		}, nil
	}
	if result != req.SendSmsCode {
		return &pb.RegisterPassengerResp{
			Code:    400,
			Message: "验证码错误",
		}, nil
	}

	// 从配置文件获取昵称前缀
	appConfig := &global.AppConf.AppConfig
	nikeName := appConfig.NicknamePrefix + req.Mobile

	lxhPassenger := model.LxhPassenger{
		NickName: nikeName,
		Mobile:   req.Mobile,
	}
	if err = global.DB.Create(&lxhPassenger).Error; err != nil {
		return &pb.RegisterPassengerResp{
			Code:    503,
			Message: "服务器异常,请重试",
		}, nil
	}
	return &pb.RegisterPassengerResp{
		Code:    200,
		Message: "注册成功",
	}, nil
}

// LoginPassenger implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) LoginPassenger(ctx context.Context, req *pb.LoginPassengerReq) (*pb.LoginPassengerResp, error) {
	ctx = context.Background()
	result, err := global.Rdb.Get(ctx, "sendSms"+req.Mobile).Result()
	if err != nil {
		return &pb.LoginPassengerResp{
			Code:    601,
			Message: "验证码获取失败",
		}, nil
	}
	var passenger model.LxhPassenger
	if err = global.DB.Debug().Where("mobile = ?", req.Mobile).Limit(1).Find(&passenger).Error; err != nil {
		return &pb.LoginPassengerResp{
			Code:    503,
			Message: "服务器异常",
		}, nil
	}
	if passenger.Id == 0 {
		return &pb.LoginPassengerResp{
			Code:    401,
			Message: "用户未注册",
		}, nil
	}
	if result != req.SendSmsCode {
		return &pb.LoginPassengerResp{
			Code:    400,
			Message: "验证码错误",
		}, nil
	}
	return &pb.LoginPassengerResp{
		Code:        200,
		Message:     "登录成功",
		PassengerId: int16(passenger.Id),
	}, nil
}

// getRecentDestinations 获取用户最近目的地（从数据库读取真实的地理信息）
func (s *PassengerServiceImpl) getRecentDestinations(passengerId int16) []string {
	// TODO: 从数据库查询用户的最近订单目的地
	// 这里先从region表中获取一些热门城市作为模拟数据
	var regions []model.Region

	// 获取一些热门城市（这里随机选择一些区域）
	if err := global.DB.Where("level = 3").Limit(10).Find(&regions).Error; err != nil {
		// 如果查询失败，返回默认数据
		return []string{
			"北京市朝阳区国贸中心",
			"上海市浦东新区陆家嘴",
			"广州市天河区珠江新城",
			"深圳市南山区科技园",
		}
	}

	var destinations []string
	for i, region := range regions {
		if i >= 4 { // 只取前4个
			break
		}
		// 获取区域的完整路径
		fullPath := s.getRegionFullPath(region.Code)
		destinations = append(destinations, fullPath)
	}

	if len(destinations) == 0 {
		// 如果没有数据，返回默认值
		return []string{
			"北京市朝阳区国贸中心",
			"上海市浦东新区陆家嘴",
			"广州市天河区珠江新城",
			"深圳市南山区科技园",
		}
	}

	return destinations
}

// getRegionFullPath 获取区域的完整路径
func (s *PassengerServiceImpl) getRegionFullPath(regionCode int64) string {
	var currentRegion model.Region
	if err := global.DB.Where("code = ?", regionCode).First(&currentRegion).Error; err != nil {
		return "未知地区"
	}

	pathParts := []string{currentRegion.Name}

	// 递归获取上级区域
	if currentRegion.Pcode != 0 {
		var parentRegion model.Region
		if err := global.DB.Where("code = ?", currentRegion.Pcode).First(&parentRegion).Error; err == nil {
			// 继续递归获取上级
			if parentRegion.Pcode != 0 {
				var grandParentRegion model.Region
				if err := global.DB.Where("code = ?", parentRegion.Pcode).First(&grandParentRegion).Error; err == nil {
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

// getNearbyCars 获取附近车辆（从数据库读取在线司机信息）
func (s *PassengerServiceImpl) getNearbyCars(location string) []*pb.CarInfo {
	// 查询状态为online的司机
	var drivers []model.LxhDriver
	if err := global.DB.Where("status = ?", "online").Limit(10).Find(&drivers).Error; err != nil {
		// 如果查询失败，返回空数组
		return []*pb.CarInfo{}
	}

	// 如果没有在线司机，返回空数组
	if len(drivers) == 0 {
		return []*pb.CarInfo{}
	}

	var nearbyCarInfos []*pb.CarInfo

	// 车型列表，用于随机分配
	carTypes := []string{"经济型", "舒适型", "豪华型", "商务型"}

	// 车牌前缀，用于生成真实车牌
	platePrefix := []string{"京A", "京B", "京C", "京D", "京E", "沪A", "沪B", "粤A", "粤B"}

	for i, driver := range drivers {
		// 动态生成距离（0.3-3.0公里）
		distance := 0.3 + rand.Float64()*2.7

		// 根据距离计算预计到达时间（距离*2分钟，最少2分钟）
		estimatedTime := int32(math.Max(2, distance*2))

		// 随机选择车型
		carType := carTypes[rand.Intn(len(carTypes))]

		// 生成真实车牌号：前缀+5位数字
		platePrefix := platePrefix[rand.Intn(len(platePrefix))]
		plateNumber := fmt.Sprintf("%s%05d", platePrefix, rand.Intn(100000))

		// 生成评分（4.0-5.0之间）
		rating := 4.0 + rand.Float64()

		// 处理司机姓名，如果没有昵称则使用姓名
		driverName := driver.NickName
		if driverName == "" {
			driverName = driver.Name
		}
		if driverName == "" {
			driverName = fmt.Sprintf("司机%d", driver.Id)
		}

		carInfo := &pb.CarInfo{
			CarId:         int16(driver.Id),
			CarType:       carType,
			LicensePlate:  plateNumber,
			Distance:      distance,
			EstimatedTime: int16(estimatedTime),
			DriverName:    driverName,
			Rating:        rating,
		}

		nearbyCarInfos = append(nearbyCarInfos, carInfo)

		// 最多返回5辆车
		if i >= 4 {
			break
		}
	}

	// 按距离排序，最近的在前面
	sort.Slice(nearbyCarInfos, func(i, j int) bool {
		return nearbyCarInfos[i].Distance < nearbyCarInfos[j].Distance
	})

	return nearbyCarInfos
}

// getServiceInfo 获取服务信息（从数据库读取，这里暂时使用模拟数据）
func (s *PassengerServiceImpl) getServiceInfo() []*pb.ServiceInfo {
	// TODO: 从数据库或配置文件读取服务信息
	return []*pb.ServiceInfo{
		{
			ServiceName: "快车",
			ServiceDesc: "经济实惠，快速到达",
			ServiceIcon: "icon_kuaiche",
			ServiceUrl:  "/service/kuaiche",
		},
		{
			ServiceName: "专车",
			ServiceDesc: "舒适体验，专业服务",
			ServiceIcon: "icon_zhuanche",
			ServiceUrl:  "/service/zhuanche",
		},
		{
			ServiceName: "豪华车",
			ServiceDesc: "高端体验，尊贵享受",
			ServiceIcon: "icon_haohua",
			ServiceUrl:  "/service/haohua",
		},
		{
			ServiceName: "代驾",
			ServiceDesc: "安全代驾，放心回家",
			ServiceIcon: "icon_daijia",
			ServiceUrl:  "/service/daijia",
		},
	}
}

// getCurrentLocationFromRegion 根据用户提供的位置字符串，从Region数据库中获取匹配的地区信息
func (s *PassengerServiceImpl) getCurrentLocationFromRegion(locationStr string) string {
	if locationStr == "" {
		// 从配置获取默认位置
		return global.AppConf.AppConfig.DefaultLocation
	}

	// 在region表中搜索匹配的地区
	var regions []model.Region
	if err := global.DB.Where("name LIKE ? OR sname LIKE ? OR mername LIKE ?",
		"%"+locationStr+"%", "%"+locationStr+"%", "%"+locationStr+"%").
		Limit(1).Find(&regions).Error; err != nil || len(regions) == 0 {
		// 如果找不到匹配的地区，返回原始字符串
		return locationStr
	}

	// 返回匹配地区的完整路径
	return s.getRegionFullPath(regions[0].Code)
}

// HomePage implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) HomePage(ctx context.Context, req *pb.HomePageReq) (*pb.HomePageResp, error) {
	// 验证用户是否存在
	var passenger model.LxhPassenger
	if err := global.DB.Debug().Where("id = ?", req.PassengerId).First(&passenger).Error; err != nil {
		return &pb.HomePageResp{
			Code:    401,
			Message: "用户不存在",
		}, nil
	}

	// 从配置文件获取应用配置
	appConfig := &global.AppConf.AppConfig
	weatherConfig := &global.AppConf.WeatherConfig

	// 获取用户的最近目的地（从region数据库获取真实地理信息）
	recentDestinations := s.getRecentDestinations(req.PassengerId)

	// 处理当前位置（利用region数据库进行地理位置匹配）
	currentLocation := appConfig.DefaultLocation
	if req.Location != nil {
		currentLocation = s.getCurrentLocationFromRegion(*req.Location)
	}

	// 获取附近车辆信息
	nearbyCars := s.getNearbyCars(currentLocation)

	// 获取服务信息
	services := s.getServiceInfo()

	// 获取当前时间用于欢迎信息
	currentTime := time.Now()
	var greeting string
	hour := currentTime.Hour()
	if hour < 12 {
		greeting = "早上好"
	} else if hour < 18 {
		greeting = "下午好"
	} else {
		greeting = "晚上好"
	}

	// 使用配置的应用名称
	welcomeMessage := fmt.Sprintf("%s，%s！%s", greeting, passenger.NickName, appConfig.WelcomeMessage)

	return &pb.HomePageResp{
		Code:    200,
		Message: "获取成功",
		Data: &pb.HomePageData{
			WelcomeMessage:     welcomeMessage,
			CurrentLocation:    currentLocation,
			RecentDestinations: recentDestinations,
			NearbyCars:         nearbyCars,
			WeatherInfo:        weatherConfig.DefaultWeather,
			Services:           services,
		},
	}, nil
}

// CallACar implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) CallACar(ctx context.Context, req *pb.CallACarReq) (*pb.CallACarResp, error) {
	// 验证用户是否存在
	var passenger model.LxhPassenger
	if err := global.DB.Debug().Where("id = ?", req.PassengerId).First(&passenger).Error; err != nil {
		return &pb.CallACarResp{
			Code:    401,
			Message: "用户不存在",
		}, nil
	}

	// 验证起点和终点是否为空
	if req.StartingPlace == "" || req.Destination == "" {
		return &pb.CallACarResp{
			Code:    400,
			Message: "起点和终点不能为空",
		}, nil
	}

	// 地址标准化：尝试从region数据库中匹配和标准化地址
	standardizedStart := s.getCurrentLocationFromRegion(req.StartingPlace)
	standardizedDest := s.getCurrentLocationFromRegion(req.Destination)

	// 创建订单记录
	orderData := model.LxhOrder{
		PassengerId: int64(req.PassengerId),
		StartAddr:   standardizedStart, // 使用标准化后的地址
		EndEnd:      standardizedDest,  // 使用标准化后的地址
		OrderStatus: "待接单",
		StartTime:   time.Now(),
	}

	if err := global.DB.Create(&orderData).Error; err != nil {
		return &pb.CallACarResp{
			Code:    503,
			Message: "服务器异常，叫车失败",
		}, nil
	}

	// 返回成功信息，包含标准化后的地址信息
	return &pb.CallACarResp{
		Code:    200,
		Message: fmt.Sprintf("叫车成功！订单号：%d，起点：%s，终点：%s，正在为您匹配司机，请耐心等待。", orderData.Id, standardizedStart, standardizedDest),
	}, nil
}

// GetPassengerInfo implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) GetPassengerInfo(ctx context.Context, req *pb.GetPassengerInfoReq) (*pb.GetPassengerInfoResp, error) {
	var passengerModel model.LxhPassenger
	if err := global.DB.Where("id = ?", req.PassengerId).First(&passengerModel).Error; err != nil {
		return &pb.GetPassengerInfoResp{
			Code:    404,
			Message: "用户不存在",
		}, nil
	}

	return &pb.GetPassengerInfoResp{
		Code:    200,
		Message: "获取成功",
		Data: &pb.PassengerInfo{
			Id:       int16(passengerModel.Id),
			Name:     passengerModel.Name,
			NickName: passengerModel.NickName,
			FileId:   passengerModel.FileId,
			Mobile:   passengerModel.Mobile,
			Age:      passengerModel.Age,
			Sex:      passengerModel.Sex,
			Mileage:  passengerModel.Mileage,
		},
	}, nil
}

// UpdatePassengerInfo implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) UpdatePassengerInfo(ctx context.Context, req *pb.UpdatePassengerInfoReq) (*pb.UpdatePassengerInfoResp, error) {
	var passengerModel model.LxhPassenger
	if err := global.DB.Where("id = ?", req.PassengerId).First(&passengerModel).Error; err != nil {
		return &pb.UpdatePassengerInfoResp{
			Code:    404,
			Message: "用户不存在",
		}, nil
	}

	updateData := make(map[string]interface{})
	if req.Name != nil {
		updateData["name"] = *req.Name
	}
	if req.NickName != nil {
		updateData["nick_name"] = *req.NickName
	}
	if req.FileId != nil {
		updateData["file_id"] = *req.FileId
	}
	if req.Age != nil {
		updateData["age"] = *req.Age
	}
	if req.Sex != nil {
		updateData["sex"] = *req.Sex
	}

	if len(updateData) > 0 {
		if err := global.DB.Model(&passengerModel).Updates(updateData).Error; err != nil {
			return &pb.UpdatePassengerInfoResp{
				Code:    503,
				Message: "更新失败",
			}, nil
		}
	}

	return &pb.UpdatePassengerInfoResp{
		Code:    200,
		Message: "更新成功",
	}, nil
}

// CreateOrder implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) CreateOrder(ctx context.Context, req *pb.CreateOrderReq) (*pb.CreateOrderResp, error) {
	// 验证用户是否存在
	var passengerModel model.LxhPassenger
	if err := global.DB.Where("id = ?", req.PassengerId).First(&passengerModel).Error; err != nil {
		return &pb.CreateOrderResp{
			Code:    404,
			Message: "用户不存在",
		}, nil
	}

	// 生成订单号
	orderCode := fmt.Sprintf("O%d%d", time.Now().Unix(), req.PassengerId)

	// 创建订单
	order := model.LxhOrder{
		OrderCode:   orderCode,
		PassengerId: int64(req.PassengerId),
		StartAddr:   req.StartAddr,
		EndEnd:      req.EndAddr,
		OrderStatus: "待接单",
		OrderType:   req.OrderType,
		PayStatus:   "未支付",
		StartTime:   time.Now(),
	}

	if err := global.DB.Create(&order).Error; err != nil {
		return &pb.CreateOrderResp{
			Code:    503,
			Message: "创建订单失败",
		}, nil
	}

	// 创建路线记录
	routeRecord := model.LxhRouteRecord{
		OrderId:      order.Id,
		PassengerId:  int64(req.PassengerId),
		StartAddress: req.StartAddr,
		StartLng:     req.StartLng,
		StartLat:     req.StartLat,
		EndAddress:   req.EndAddr,
		EndLng:       req.EndLng,
		EndLat:       req.EndLat,
		RouteStatus:  "planning",
		StartTime:    time.Now(),
	}

	if err := global.DB.Create(&routeRecord).Error; err != nil {
		// 路线记录创建失败不影响订单创建
		fmt.Printf("创建路线记录失败: %v", err)
	}

	return &pb.CreateOrderResp{
		Code:      200,
		Message:   "订单创建成功",
		OrderCode: &orderCode,
	}, nil
}

// GetOrderList implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) GetOrderList(ctx context.Context, req *pb.GetOrderListReq) (*pb.GetOrderListResp, error) {
	page := int32(1)
	pageSize := int32(10)

	if req.Page != nil {
		page = *req.Page
	}
	if req.PageSize != nil {
		pageSize = *req.PageSize
	}

	offset := (page - 1) * pageSize

	query := global.DB.Where("passenger_id = ?", req.PassengerId)
	if req.Status != nil {
		query = query.Where("order_status = ?", *req.Status)
	}

	var orders []model.LxhOrder
	var total int64

	query.Count(&total)
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Order("start_time DESC").Find(&orders).Error; err != nil {
		return &pb.GetOrderListResp{
			Code:    503,
			Message: "查询失败",
		}, nil
	}

	var orderInfos []*pb.OrderInfo
	for _, order := range orders {
		orderInfos = append(orderInfos, &pb.OrderInfo{
			Id:          order.Id,
			OrderCode:   order.OrderCode,
			Amount:      order.Amount,
			OrderStatus: order.OrderStatus,
			StartAddr:   order.StartAddr,
			EndEnd:      order.EndEnd,
			Driver:      order.Driver,
			StartTime:   order.StartTime.Format("2006-01-02 15:04:05"),
			EndTime:     order.EndTime.Format("2006-01-02 15:04:05"),
			PayStatus:   order.PayStatus,
			PayType:     order.PayType,
			OrderType:   order.OrderType,
		})
	}

	totalInt32 := int32(total)
	return &pb.GetOrderListResp{
		Code:    200,
		Message: "查询成功",
		Data:    orderInfos,
		Total:   &totalInt32,
	}, nil
}

// GetOrderDetail implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) GetOrderDetail(ctx context.Context, req *pb.GetOrderDetailReq) (*pb.GetOrderDetailResp, error) {
	var order model.LxhOrder
	if err := global.DB.Where("id = ? AND passenger_id = ?", req.OrderId, req.PassengerId).First(&order).Error; err != nil {
		return &pb.GetOrderDetailResp{
			Code:    404,
			Message: "订单不存在",
		}, nil
	}

	return &pb.GetOrderDetailResp{
		Code:    200,
		Message: "查询成功",
		Data: &pb.OrderInfo{
			Id:          order.Id,
			OrderCode:   order.OrderCode,
			Amount:      order.Amount,
			OrderStatus: order.OrderStatus,
			StartAddr:   order.StartAddr,
			EndEnd:      order.EndEnd,
			Driver:      order.Driver,
			StartTime:   order.StartTime.Format("2006-01-02 15:04:05"),
			EndTime:     order.EndTime.Format("2006-01-02 15:04:05"),
			PayStatus:   order.PayStatus,
			PayType:     order.PayType,
			OrderType:   order.OrderType,
		},
	}, nil
}

// CancelOrder implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) CancelOrder(ctx context.Context, req *pb.CancelOrderReq) (*pb.CancelOrderResp, error) {
	var order model.LxhOrder
	if err := global.DB.Where("id = ? AND passenger_id = ?", req.OrderId, req.PassengerId).First(&order).Error; err != nil {
		return &pb.CancelOrderResp{
			Code:    404,
			Message: "订单不存在",
		}, nil
	}

	if order.OrderStatus != "待接单" && order.OrderStatus != "已接单" {
		return &pb.CancelOrderResp{
			Code:    400,
			Message: "订单状态不允许取消",
		}, nil
	}

	updateData := map[string]interface{}{
		"order_status":   "已取消",
		"confirm_by":     "乘客",
		"confirm_person": req.PassengerId,
		"confirm_reason": req.Reason,
		"end_time":       time.Now(),
	}

	if req.Remark != nil {
		updateData["confirm_remark"] = *req.Remark
	}

	if err := global.DB.Model(&order).Updates(updateData).Error; err != nil {
		return &pb.CancelOrderResp{
			Code:    503,
			Message: "取消订单失败",
		}, nil
	}

	return &pb.CancelOrderResp{
		Code:    200,
		Message: "订单已取消",
	}, nil
}

// EvaluateOrder implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) EvaluateOrder(ctx context.Context, req *pb.EvaluateOrderReq) (*pb.EvaluateOrderResp, error) {
	var order model.LxhOrder
	if err := global.DB.Where("id = ? AND passenger_id = ?", req.OrderId, req.PassengerId).First(&order).Error; err != nil {
		return &pb.EvaluateOrderResp{
			Code:    404,
			Message: "订单不存在",
		}, nil
	}

	if order.OrderStatus != "已完成" {
		return &pb.EvaluateOrderResp{
			Code:    400,
			Message: "只能评价已完成的订单",
		}, nil
	}

	// 这里可以扩展评价表存储详细评价信息
	// 暂时将评价信息存储在订单的备注字段
	evaluateInfo := fmt.Sprintf("评分:%d", req.Rating)
	if req.Comment != nil {
		evaluateInfo += fmt.Sprintf(",评价:%s", *req.Comment)
	}
	if req.Tags != nil && len(req.Tags) > 0 {
		evaluateInfo += fmt.Sprintf(",标签:%v", req.Tags)
	}

	if err := global.DB.Model(&order).Update("confirm_remark", evaluateInfo).Error; err != nil {
		return &pb.EvaluateOrderResp{
			Code:    503,
			Message: "评价失败",
		}, nil
	}

	return &pb.EvaluateOrderResp{
		Code:    200,
		Message: "评价成功",
	}, nil
}

// GetFavoriteLocations implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) GetFavoriteLocations(ctx context.Context, req *pb.GetFavoriteLocationsReq) (*pb.GetFavoriteLocationsResp, error) {
	query := global.DB.Where("passenger_id = ?", req.PassengerId)
	if req.LocationType != nil {
		query = query.Where("location_type = ?", *req.LocationType)
	}

	var locations []model.LxhFavoriteLocation
	if err := query.Order("sort_order ASC, created_at DESC").Find(&locations).Error; err != nil {
		return &pb.GetFavoriteLocationsResp{
			Code:    503,
			Message: "查询失败",
		}, nil
	}

	var favoriteLocations []*pb.FavoriteLocation
	for _, loc := range locations {
		favoriteLocations = append(favoriteLocations, &pb.FavoriteLocation{
			Id:           loc.Id,
			LocationType: loc.LocationType,
			LocationName: loc.LocationName,
			Address:      loc.Address,
			Lng:          loc.Lng,
			Lat:          loc.Lat,
			Province:     loc.Province,
			City:         loc.City,
			District:     loc.District,
			UsageCount:   int32(loc.UsageCount),
			IsDefault:    loc.IsDefault == 1,
		})
	}

	return &pb.GetFavoriteLocationsResp{
		Code:    200,
		Message: "查询成功",
		Data:    favoriteLocations,
	}, nil
}

// AddFavoriteLocation implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) AddFavoriteLocation(ctx context.Context, req *pb.AddFavoriteLocationReq) (*pb.AddFavoriteLocationResp, error) {
	// 检查是否已存在相同地址
	var existingLoc model.LxhFavoriteLocation
	if err := global.DB.Where("passenger_id = ? AND address = ?", req.PassengerId, req.Address).First(&existingLoc).Error; err == nil {
		return &pb.AddFavoriteLocationResp{
			Code:    400,
			Message: "该地址已添加到收藏",
		}, nil
	}

	isDefault := 0
	if req.IsDefault != nil && *req.IsDefault {
		isDefault = 1
		// 如果设为默认地址，需要先取消其他默认地址
		global.DB.Model(&model.LxhFavoriteLocation{}).
			Where("passenger_id = ? AND location_type = ?", req.PassengerId, req.LocationType).
			Update("is_default", 0)
	}

	location := model.LxhFavoriteLocation{
		PassengerId:  int64(req.PassengerId),
		LocationType: req.LocationType,
		LocationName: req.LocationName,
		Address:      req.Address,
		Lng:          req.Lng,
		Lat:          req.Lat,
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		IsDefault:    isDefault,
		SortOrder:    0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := global.DB.Create(&location).Error; err != nil {
		return &pb.AddFavoriteLocationResp{
			Code:    503,
			Message: "添加失败",
		}, nil
	}

	return &pb.AddFavoriteLocationResp{
		Code:       200,
		Message:    "添加成功",
		LocationId: &location.Id,
	}, nil
}

// DeleteFavoriteLocation implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) DeleteFavoriteLocation(ctx context.Context, req *pb.DeleteFavoriteLocationReq) (*pb.DeleteFavoriteLocationResp, error) {
	result := global.DB.Where("id = ? AND passenger_id = ?", req.LocationId, req.PassengerId).Delete(&model.LxhFavoriteLocation{})
	if result.Error != nil {
		return &pb.DeleteFavoriteLocationResp{
			Code:    503,
			Message: "删除失败",
		}, nil
	}

	if result.RowsAffected == 0 {
		return &pb.DeleteFavoriteLocationResp{
			Code:    404,
			Message: "地址不存在",
		}, nil
	}

	return &pb.DeleteFavoriteLocationResp{
		Code:    200,
		Message: "删除成功",
	}, nil
}

// GetHotLocations implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) GetHotLocations(ctx context.Context, req *pb.GetHotLocationsReq) (*pb.GetHotLocationsResp, error) {
	limit := int32(10)
	if req.Limit != nil {
		limit = *req.Limit
	}

	query := global.DB.Where("city = ? AND status = 1", req.City)
	if req.Category != nil {
		query = query.Where("category = ?", *req.Category)
	}

	var hotLocations []model.LxhHotLocation
	if err := query.Order("hot_score DESC, order_count DESC").Limit(int(limit)).Find(&hotLocations).Error; err != nil {
		return &pb.GetHotLocationsResp{
			Code:    503,
			Message: "查询失败",
		}, nil
	}

	var locations []*pb.HotLocation
	for _, loc := range hotLocations {
		locations = append(locations, &pb.HotLocation{
			Id:           loc.Id,
			LocationName: loc.LocationName,
			Address:      loc.Address,
			Lng:          loc.Lng,
			Lat:          loc.Lat,
			City:         loc.City,
			Category:     loc.Category,
			HotScore:     loc.HotScore,
		})
	}

	return &pb.GetHotLocationsResp{
		Code:    200,
		Message: "查询成功",
		Data:    locations,
	}, nil
}

// BindWechat implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) BindWechat(ctx context.Context, req *pb.BindWechatReq) (*pb.BindWechatResp, error) {
	// 这里需要根据微信授权码获取用户信息
	// 暂时简化处理，实际需要调用微信API
	// TODO: 实现微信授权流程，获取openid和用户信息

	// 检查用户是否已绑定微信
	var existingBind model.LxhUserWechatBind
	if err := global.DB.Where("passenger_id = ? AND bind_status = 1", req.PassengerId).First(&existingBind).Error; err == nil {
		return &pb.BindWechatResp{
			Code:    400,
			Message: "该账号已绑定微信",
		}, nil
	}

	// 模拟openid，实际应该从微信API获取
	mockOpenid := fmt.Sprintf("mock_openid_%d_%d", req.PassengerId, time.Now().Unix())

	// 创建微信用户记录
	wechatUser := model.LxhWechatUser{
		Openid:    mockOpenid,
		Nickname:  fmt.Sprintf("微信用户_%d", req.PassengerId),
		Sex:       1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := global.DB.Create(&wechatUser).Error; err != nil {
		return &pb.BindWechatResp{
			Code:    503,
			Message: "绑定失败",
		}, nil
	}

	// 创建绑定记录
	bind := model.LxhUserWechatBind{
		PassengerId: int64(req.PassengerId),
		Openid:      mockOpenid,
		BindType:    1,
		BindStatus:  1,
		BindTime:    time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := global.DB.Create(&bind).Error; err != nil {
		return &pb.BindWechatResp{
			Code:    503,
			Message: "绑定失败",
		}, nil
	}

	return &pb.BindWechatResp{
		Code:    200,
		Message: "绑定成功",
	}, nil
}

// UnbindWechat implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) UnbindWechat(ctx context.Context, req *pb.UnbindWechatReq) (*pb.UnbindWechatResp, error) {
	result := global.DB.Model(&model.LxhUserWechatBind{}).
		Where("passenger_id = ? AND bind_status = 1", req.PassengerId).
		Updates(map[string]interface{}{
			"bind_status": 2,
			"unbind_time": time.Now(),
			"updated_at":  time.Now(),
		})

	if result.Error != nil {
		return &pb.UnbindWechatResp{
			Code:    503,
			Message: "解绑失败",
		}, nil
	}

	if result.RowsAffected == 0 {
		return &pb.UnbindWechatResp{
			Code:    400,
			Message: "未找到绑定记录",
		}, nil
	}

	return &pb.UnbindWechatResp{
		Code:    200,
		Message: "解绑成功",
	}, nil
}

// GetWechatBindStatus implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) GetWechatBindStatus(ctx context.Context, req *pb.GetWechatBindStatusReq) (*pb.GetWechatBindStatusResp, error) {
	var bind model.LxhUserWechatBind
	if err := global.DB.Where("passenger_id = ? AND bind_status = 1", req.PassengerId).First(&bind).Error; err != nil {
		return &pb.GetWechatBindStatusResp{
			Code:    200,
			Message: "查询成功",
			Data: &pb.WechatBindInfo{
				IsBind: false,
			},
		}, nil
	}

	// 获取微信用户信息
	var wechatUser model.LxhWechatUser
	var bindInfo pb.WechatBindInfo
	bindInfo.IsBind = true

	if err := global.DB.Where("openid = ?", bind.Openid).First(&wechatUser).Error; err == nil {
		bindInfo.Nickname = &wechatUser.Nickname
		bindInfo.Headimgurl = &wechatUser.Headimgurl
	}

	bindTimeStr := bind.BindTime.Format("2006-01-02 15:04:05")
	bindInfo.BindTime = &bindTimeStr

	return &pb.GetWechatBindStatusResp{
		Code:    200,
		Message: "查询成功",
		Data:    &bindInfo,
	}, nil
}

// GetRouteRecords implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) GetRouteRecords(ctx context.Context, req *pb.GetRouteRecordsReq) (*pb.GetRouteRecordsResp, error) {
	page := int32(1)
	pageSize := int32(10)

	if req.Page != nil {
		page = *req.Page
	}
	if req.PageSize != nil {
		pageSize = *req.PageSize
	}

	offset := (page - 1) * pageSize

	var routes []model.LxhRouteRecord
	var total int64

	query := global.DB.Where("passenger_id = ?", req.PassengerId)
	query.Count(&total)

	if err := query.Offset(int(offset)).Limit(int(pageSize)).Order("created_at DESC").Find(&routes).Error; err != nil {
		return &pb.GetRouteRecordsResp{
			Code:    503,
			Message: "查询失败",
		}, nil
	}

	var routeRecords []*pb.RouteRecord
	for _, route := range routes {
		routeRecords = append(routeRecords, &pb.RouteRecord{
			Id:            route.Id,
			OrderId:       route.OrderId,
			StartAddress:  route.StartAddress,
			StartLng:      route.StartLng,
			StartLat:      route.StartLat,
			EndAddress:    route.EndAddress,
			EndLng:        route.EndLng,
			EndLat:        route.EndLat,
			Distance:      route.Distance,
			EstimatedTime: int32(route.EstimatedTime),
			ActualTime:    int32(route.ActualTime),
			RouteStatus:   route.RouteStatus,
			StartTime:     route.StartTime.Format("2006-01-02 15:04:05"),
			EndTime:       route.EndTime.Format("2006-01-02 15:04:05"),
		})
	}

	totalInt32 := int32(total)
	return &pb.GetRouteRecordsResp{
		Code:    200,
		Message: "查询成功",
		Data:    routeRecords,
		Total:   &totalInt32,
	}, nil
}

// SearchAddress implements the PassengerServiceImpl interface.
func (s *PassengerServiceImpl) SearchAddress(ctx context.Context, req *pb.SearchAddressReq) (*pb.SearchAddressResp, error) {
	limit := int32(10)
	if req.Limit != nil {
		limit = *req.Limit
	}

	// 从热门地点表搜索
	var hotLocations []model.LxhHotLocation
	query := global.DB.Where("status = 1 AND (location_name LIKE ? OR address LIKE ?)",
		"%"+req.Keyword+"%", "%"+req.Keyword+"%")

	if req.City != nil {
		query = query.Where("city = ?", *req.City)
	}

	if err := query.Limit(int(limit)).Order("hot_score DESC").Find(&hotLocations).Error; err != nil {
		return &pb.SearchAddressResp{
			Code:    503,
			Message: "搜索失败",
		}, nil
	}

	var suggestions []*pb.AddressSuggestion
	for _, loc := range hotLocations {
		distance := 0.0
		if req.Lng != nil && req.Lat != nil {
			// 简单计算距离（实际应使用更精确的算法）
			dlng := loc.Lng - *req.Lng
			dlat := loc.Lat - *req.Lat
			distance = dlng*dlng + dlat*dlat // 简化距离计算
		}

		suggestions = append(suggestions, &pb.AddressSuggestion{
			Name:     loc.LocationName,
			Address:  loc.Address,
			Lng:      loc.Lng,
			Lat:      loc.Lat,
			District: loc.District,
			Distance: distance,
		})
	}

	return &pb.SearchAddressResp{
		Code:    200,
		Message: "搜索成功",
		Data:    suggestions,
	}, nil
}
