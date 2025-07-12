package main

import (
	"cart/biz/utils"
	pb "cart/kitex_gen/cart/passenger"
	"cart/rpc/basic/global"
	"cart/rpc/basic/model"
	"context"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PassengerServiceImpl implements the last service interface defined in the IDL.
type PassengerServiceImpl struct{}

// SendSms 发送短信验证码接口
// 生成并发送短信验证码到指定手机号
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

// RegisterPassenger 乘客注册接口
// 验证短信验证码，创建新乘客账户
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

// LoginPassenger 乘客登录接口
// 验证手机号和短信验证码，返回乘客ID
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

// HomePage 首页接口
// 获取首页数据，包括欢迎信息、当前位置、附近车辆、服务信息等
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

	// 获取用户的最近目的地
	recentDestinations := utils.GetRecentDestinations(req.PassengerId, global.DB)

	// 处理当前位置
	currentLocation := appConfig.DefaultLocation
	if req.Location != nil {
		currentLocation = utils.GetCurrentLocationFromRegion(*req.Location, global.DB, appConfig.DefaultLocation)
	}

	// 获取附近车辆信息
	nearbyCars := utils.GetNearbyCars(currentLocation, global.DB)

	// 获取服务信息
	services := utils.GetServiceInfo()

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

// CallACar 叫车接口
// 乘客发起叫车请求，创建订单
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

	// 地址标准化
	standardizedStart := utils.GetCurrentLocationFromRegion(req.StartingPlace, global.DB, global.AppConf.AppConfig.DefaultLocation)
	standardizedDest := utils.GetCurrentLocationFromRegion(req.Destination, global.DB, global.AppConf.AppConfig.DefaultLocation)

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

// GetPassengerInfo 获取乘客信息接口
// 根据乘客ID查询并返回乘客的详细信息
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

// UpdatePassengerInfo 更新乘客信息接口
// 更新乘客的基本信息，如姓名、昵称、头像等
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

// CreateOrder 创建订单接口
// 创建新的出行订单，支持多种订单类型
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

	// 构造MongoDB订单结构体
	mongoOrder := model.MongoOrder{
		OrderId:     orderCode,
		PassengerId: int64(req.PassengerId),
		StartAddr:   req.StartAddr,
		EndAddr:     req.EndAddr,
		OrderStatus: "待接单",
		OrderType:   req.OrderType,
		PayStatus:   "未支付",
		Amount:      0,
		StartTime:   time.Now(),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}

	// 写入MongoDB
	collection := global.MongoDB.Database("cart").Collection("orders")
	_, err := collection.InsertOne(ctx, mongoOrder)
	if err != nil {
		return &pb.CreateOrderResp{
			Code:    503,
			Message: "创建订单失败(MongoDB)",
		}, nil
	}

	return &pb.CreateOrderResp{
		Code:      200,
		Message:   "订单创建成功(MongoDB)",
		OrderCode: &orderCode,
	}, nil
}

// GetOrderList 获取订单列表接口
// 查询乘客的订单历史，支持分页和状态筛选
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

	// 1. 查MySQL已完成订单
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

	// 2. 查MongoDB未完成订单
	mongoFilter := bson.M{"passenger_id": int64(req.PassengerId)}
	if req.Status != nil {
		mongoFilter["order_status"] = *req.Status
	} else {
		mongoFilter["order_status"] = bson.M{"$in": []string{"待接单", "已接单", "进行中"}}
	}
	collection := global.MongoDB.Database("cart").Collection("orders")
	findOpts := options.Find().SetSort(bson.D{{"create_time", -1}})
	cur, err := collection.Find(ctx, mongoFilter, findOpts)
	if err != nil {
		return &pb.GetOrderListResp{
			Code:    503,
			Message: "MongoDB查询失败",
		}, nil
	}
	defer cur.Close(ctx)
	var mongoOrders []model.MongoOrder
	for cur.Next(ctx) {
		var mo model.MongoOrder
		if err := cur.Decode(&mo); err == nil {
			mongoOrders = append(mongoOrders, mo)
		}
	}

	// 3. 合并结果
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
	for _, mo := range mongoOrders {
		orderInfos = append(orderInfos, &pb.OrderInfo{
			OrderCode:   mo.OrderId,
			Amount:      mo.Amount,
			OrderStatus: mo.OrderStatus,
			StartAddr:   mo.StartAddr,
			EndEnd:      mo.EndAddr,
			Driver:      mo.DriverId,
			StartTime:   mo.StartTime.Format("2006-01-02 15:04:05"),
			EndTime:     mo.EndTime.Format("2006-01-02 15:04:05"),
			PayStatus:   mo.PayStatus,
			OrderType:   mo.OrderType,
		})
	}
	// 4. 排序（按StartTime倒序）
	sort.Slice(orderInfos, func(i, j int) bool {
		return orderInfos[i].StartTime > orderInfos[j].StartTime
	})
	// 5. 分页
	start := int(offset)
	end := start + int(pageSize)
	if start > len(orderInfos) {
		start = len(orderInfos)
	}
	if end > len(orderInfos) {
		end = len(orderInfos)
	}
	orderInfosPage := orderInfos[start:end]
	totalInt32 := int32(len(orderInfos))
	return &pb.GetOrderListResp{
		Code:    200,
		Message: "查询成功",
		Data:    orderInfosPage,
		Total:   &totalInt32,
	}, nil
}

// GetOrderDetail 获取订单详情接口
// 查询指定订单的详细信息
func (s *PassengerServiceImpl) GetOrderDetail(ctx context.Context, req *pb.GetOrderDetailReq) (*pb.GetOrderDetailResp, error) {
	// 1. 先查MySQL
	var order model.LxhOrder
	if err := global.DB.Where("order_code = ? AND passenger_id = ?", req.OrderId, req.PassengerId).First(&order).Error; err == nil {
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
	// 2. 查MongoDB
	collection := global.MongoDB.Database("cart").Collection("orders")
	var mo model.MongoOrder
	err := collection.FindOne(ctx, bson.M{"order_id": req.OrderId, "passenger_id": int64(req.PassengerId)}).Decode(&mo)
	if err != nil {
		return &pb.GetOrderDetailResp{
			Code:    404,
			Message: "订单不存在",
		}, nil
	}
	return &pb.GetOrderDetailResp{
		Code:    200,
		Message: "查询成功",
		Data: &pb.OrderInfo{
			OrderCode:   mo.OrderId,
			Amount:      mo.Amount,
			OrderStatus: mo.OrderStatus,
			StartAddr:   mo.StartAddr,
			EndEnd:      mo.EndAddr,
			Driver:      mo.DriverId,
			StartTime:   mo.StartTime.Format("2006-01-02 15:04:05"),
			EndTime:     mo.EndTime.Format("2006-01-02 15:04:05"),
			PayStatus:   mo.PayStatus,
			OrderType:   mo.OrderType,
		},
	}, nil
}

// CancelOrder 取消订单接口
// 乘客取消订单，需要提供取消原因
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

// EvaluateOrder 评价订单接口
// 乘客对已完成订单进行评价，包括评分和评论
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

// GetFavoriteLocations 获取收藏地点接口
// 查询乘客收藏的常用地点列表
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

// AddFavoriteLocation 添加收藏地点接口
// 添加新的收藏地点，支持设置默认地址
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

// DeleteFavoriteLocation 删除收藏地点接口
// 删除指定的收藏地点
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

// GetHotLocations 获取热门地点接口
// 查询指定城市的热门地点列表
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

// BindWechat 绑定微信接口
// 将乘客账户与微信账户进行绑定
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

// UnbindWechat 解绑微信接口
// 解除乘客账户与微信账户的绑定关系
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

// GetWechatBindStatus 获取微信绑定状态接口
// 查询乘客账户的微信绑定状态和相关信息
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

// GetRouteRecords 获取路线记录接口
// 查询乘客的历史出行路线记录
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

// SearchAddress 搜索地址接口
// 根据关键词搜索地址，支持模糊匹配
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
