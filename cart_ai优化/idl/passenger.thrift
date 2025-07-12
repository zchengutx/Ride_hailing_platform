namespace go cart.passenger

// 短信发送请求
struct SendSmsReq{
    1: string mobile       // 手机号
    2: string sendSmsCode  // 验证码
}
// 短信发送响应
struct SendSmsResp{
    1: i16 code           // 状态码
    2: string message     // 返回消息
}
// 乘客注册请求
struct RegisterPassengerReq{
    1: string mobile       // 手机号
    2: string sendSmsCode  // 验证码
}
// 乘客注册响应
struct RegisterPassengerResp{
    1: i16 code           // 状态码
    2: string message     // 返回消息
}
// 乘客登录请求
struct LoginPassengerReq{
    1: string mobile       // 手机号
    2: string sendSmsCode  // 验证码
}
// 乘客登录响应
struct LoginPassengerResp{
    1: i16 code           // 状态码
    2: string message     // 返回消息
    3: i16 passengerId    // 乘客ID
}

// 主页请求
struct HomePageReq{
    1: i16 passengerId
    2: optional string location // 当前位置
}

// 主页响应数据
struct HomePageData{
    1: string welcome_message   // 欢迎信息
    2: string current_location  // 当前位置
    3: list<string> recent_destinations // 最近目的地
    4: list<CarInfo> nearby_cars // 附近车辆
    5: string weather_info      // 天气信息
    6: list<ServiceInfo> services // 服务信息
}

// 车辆信息
struct CarInfo{
    1: i16 car_id
    2: string car_type         // 车型
    3: string license_plate    // 车牌号
    4: double distance         // 距离（km）
    5: i16 estimated_time      // 预计到达时间（分钟）
    6: string driver_name      // 司机姓名
    7: double rating           // 评分
}

// 服务信息
struct ServiceInfo{
    1: string service_name     // 服务名称
    2: string service_desc     // 服务描述
    3: string service_icon     // 服务图标
    4: string service_url      // 服务链接
}

// 主页响应
struct HomePageResp{
    1: i16 code
    2: string message
    3: optional HomePageData data
}

// 叫车请求
struct CallACarReq{
    1: i16 passengerId      // 乘客ID
    2: string startingPlace // 出发地
    3: string destination   // 目的地
}
// 叫车响应
struct CallACarResp{
    1: i16 code            // 状态码
    2: string message      // 返回消息
}

// 个人信息管理
struct PassengerInfo{
    1: i16 id
    2: string name
    3: string nickName
    4: i64 fileId              // 头像文件ID
    5: string mobile
    6: i64 age
    7: string sex
    8: string mileage
}

struct GetPassengerInfoReq{
    1: i16 passengerId
}

struct GetPassengerInfoResp{
    1: i16 code
    2: string message
    3: optional PassengerInfo data
}

struct UpdatePassengerInfoReq{
    1: i16 passengerId
    2: optional string name
    3: optional string nickName
    4: optional i64 fileId
    5: optional i64 age
    6: optional string sex
}

struct UpdatePassengerInfoResp{
    1: i16 code
    2: string message
}

// 订单管理
struct OrderInfo{
    1: i64 id
    2: string orderCode
    3: double amount
    4: string orderStatus      // 订单状态
    5: string startAddr        // 起始地
    6: string endEnd           // 目的地
    7: i64 driver              // 司机ID
    8: string startTime        // 开始时间
    9: string endTime          // 结束时间
    10: string payStatus       // 支付状态
    11: string payType         // 支付方式
    12: string orderType       // 订单类型
}

struct GetOrderListReq{
    1: i16 passengerId
    2: optional i32 page       // 页码，默认1
    3: optional i32 pageSize   // 每页数量，默认10
    4: optional string status  // 订单状态筛选
}

struct GetOrderListResp{
    1: i16 code
    2: string message
    3: optional list<OrderInfo> data
    4: optional i32 total
}

struct GetOrderDetailReq{
    1: i16 passengerId
    2: i64 orderId
}

struct GetOrderDetailResp{
    1: i16 code
    2: string message
    3: optional OrderInfo data
}

struct CancelOrderReq{
    1: i16 passengerId
    2: i64 orderId
    3: string reason           // 取消原因
    4: optional string remark  // 取消备注
}

struct CancelOrderResp{
    1: i16 code
    2: string message
}

struct CreateOrderReq{
    1: i16 passengerId
    2: string startAddr        // 起始地址
    3: double startLng         // 起始经度
    4: double startLat         // 起始纬度
    5: string endAddr          // 目的地址
    6: double endLng           // 目的经度
    7: double endLat           // 目的纬度
    8: string orderType        // 订单类型：快车/预约
    9: optional string appointTime // 预约时间
}

struct CreateOrderResp{
    1: i16 code
    2: string message
    3: optional string orderCode
}

// 收藏地址管理
struct FavoriteLocation{
    1: i64 id
    2: string locationType     // 地点类型(home/company/custom)
    3: string locationName     // 地点名称
    4: string address          // 详细地址
    5: double lng              // 经度
    6: double lat              // 纬度
    7: string province         // 省份
    8: string city             // 城市
    9: string district         // 区县
    10: i32 usageCount         // 使用次数
    11: bool isDefault         // 是否默认地址
}

struct GetFavoriteLocationsReq{
    1: i16 passengerId
    2: optional string locationType
}

struct GetFavoriteLocationsResp{
    1: i16 code
    2: string message
    3: optional list<FavoriteLocation> data
}

struct AddFavoriteLocationReq{
    1: i16 passengerId
    2: string locationType     // home/company/custom
    3: string locationName
    4: string address
    5: double lng
    6: double lat
    7: string province
    8: string city
    9: string district
    10: optional bool isDefault
}

struct AddFavoriteLocationResp{
    1: i16 code
    2: string message
    3: optional i64 locationId
}

struct DeleteFavoriteLocationReq{
    1: i16 passengerId
    2: i64 locationId
}

struct DeleteFavoriteLocationResp{
    1: i16 code
    2: string message
}

// 热门地点推荐
struct HotLocation{
    1: i64 id
    2: string locationName
    3: string address
    4: double lng
    5: double lat
    6: string city
    7: string category         // 地点分类
    8: double hotScore         // 热度得分
}

struct GetHotLocationsReq{
    1: string city
    2: optional string category
    3: optional i32 limit      // 返回数量限制，默认10
}

struct GetHotLocationsResp{
    1: i16 code
    2: string message
    3: optional list<HotLocation> data
}

// 微信绑定管理
struct BindWechatReq{
    1: i16 passengerId
    2: string code             // 微信授权码
}

struct BindWechatResp{
    1: i16 code
    2: string message
}

struct UnbindWechatReq{
    1: i16 passengerId
}

struct UnbindWechatResp{
    1: i16 code
    2: string message
}

struct GetWechatBindStatusReq{
    1: i16 passengerId
}

struct WechatBindInfo{
    1: bool isBind             // 是否已绑定
    2: optional string nickname // 微信昵称
    3: optional string headimgurl // 微信头像
    4: optional string bindTime // 绑定时间
}

struct GetWechatBindStatusResp{
    1: i16 code
    2: string message
    3: optional WechatBindInfo data
}

// 路线记录查询
struct RouteRecord{
    1: i64 id
    2: i64 orderId
    3: string startAddress
    4: double startLng
    5: double startLat
    6: string endAddress
    7: double endLng
    8: double endLat
    9: double distance         // 距离(公里)
    10: i32 estimatedTime      // 预估时间(分钟)
    11: i32 actualTime         // 实际用时(分钟)
    12: string routeStatus     // 路径状态
    13: string startTime       // 开始时间
    14: string endTime         // 结束时间
}

struct GetRouteRecordsReq{
    1: i16 passengerId
    2: optional i32 page
    3: optional i32 pageSize
}

struct GetRouteRecordsResp{
    1: i16 code
    2: string message
    3: optional list<RouteRecord> data
    4: optional i32 total
}

// 地址搜索建议
struct SearchAddressReq{
    1: string keyword          // 搜索关键词
    2: optional string city    // 城市限制
    3: optional double lng     // 当前位置经度
    4: optional double lat     // 当前位置纬度
    5: optional i32 limit      // 返回数量，默认10
}

struct AddressSuggestion{
    1: string name             // 地址名称
    2: string address          // 详细地址
    3: double lng              // 经度
    4: double lat              // 纬度
    5: string district         // 区域
    6: double distance         // 距离当前位置(米)
}

struct SearchAddressResp{
    1: i16 code
    2: string message
    3: optional list<AddressSuggestion> data
}

// 订单评价
struct EvaluateOrderReq{
    1: i16 passengerId
    2: i64 orderId
    3: i32 rating              // 评分 1-5
    4: optional string comment // 评价内容
    5: optional list<string> tags // 评价标签
}

struct EvaluateOrderResp{
    1: i16 code
    2: string message
}

service PassengerService{
    SendSmsResp SendSms(1: SendSmsReq req)
    RegisterPassengerResp RegisterPassenger(1: RegisterPassengerReq req)
    LoginPassengerResp LoginPassenger(1: LoginPassengerReq req)
    HomePageResp HomePage(1: HomePageReq req)
    CallACarResp CallACar(1: CallACarReq req)
    GetPassengerInfoResp GetPassengerInfo(1: GetPassengerInfoReq req)
    UpdatePassengerInfoResp UpdatePassengerInfo(1: UpdatePassengerInfoReq req)
    CreateOrderResp CreateOrder(1: CreateOrderReq req)
    GetOrderListResp GetOrderList(1: GetOrderListReq req)
    GetOrderDetailResp GetOrderDetail(1: GetOrderDetailReq req)
    CancelOrderResp CancelOrder(1: CancelOrderReq req)
    EvaluateOrderResp EvaluateOrder(1: EvaluateOrderReq req)
    GetFavoriteLocationsResp GetFavoriteLocations(1: GetFavoriteLocationsReq req)
    AddFavoriteLocationResp AddFavoriteLocation(1: AddFavoriteLocationReq req)
    DeleteFavoriteLocationResp DeleteFavoriteLocation(1: DeleteFavoriteLocationReq req)
    GetHotLocationsResp GetHotLocations(1: GetHotLocationsReq req)
    
    // 微信绑定管理
    BindWechatResp BindWechat(1: BindWechatReq req)
    UnbindWechatResp UnbindWechat(1: UnbindWechatReq req)
    GetWechatBindStatusResp GetWechatBindStatus(1: GetWechatBindStatusReq req)
    GetRouteRecordsResp GetRouteRecords(1: GetRouteRecordsReq req)
    SearchAddressResp SearchAddress(1: SearchAddressReq req)
}