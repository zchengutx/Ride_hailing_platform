namespace go passengers

struct Passengers{
    1: i32 Id
    2: string Name
    3: string NickName
    4: i32    FileId
    5: string Mobile
    6: i32    Age
    7: string Sex
    8: string Mileage
}

//注册请求
struct SendSmsReq {
    1: string Mobile
    2: string Source
}

//注册响应
struct SendSmsResp{
    1: string Message
    2: i64 Code
}

//注册请求
struct RegisterReq {
    1: string Mobile
    2: string SmsCode
    3: Passengers passengers
}
//注册响应
struct RegisterResp{
    1: string Message
    2: i64 Code
    3: Passengers passengers
}

//登录请求
struct LoginReq {
    1: string Mobile
    2: string SmsCode
}

//登录响应
struct LoginResp{
    1: string Message
    2: i64 Code
    3: Passengers passengers
    4: string Token
}

// 实时定位请求
struct LocationReq {
    1: string PassengerID
}

// 实时定位响应
struct LocationResp {
    1: string Message
    2: i16 Code
    3: double Lat
    4: double Lng
}

// 行程请求
struct TripReq {
    1: string PassengerID
}

// 行程响应
struct TripResp {
    1: string Message
    2: i16 Code
    3: list<string> HistoryTrips
    4: list<string> RecommendDestinations
}

// 计价请求
struct PricingReq {
    1: string StartLocation
    2: string EndLocation
    3: string CarType
}

// 计价响应
struct PricingResp {
    1: string Message
    2: i16 Code
    3: double Price
}

// 下单请求
struct OrderReq {
    1: string PassengerID
    2: string StartLocation
    3: string EndLocation
    4: string CarType
    5: i8    OrderType // 0: 预约 1: 正常打车 2: 代人叫车
}

// 下单响应
struct OrderResp {
    1: string Message
    2: i16 Code
    3: string OrderID
}

// 路线规划请求
struct RoutePlanReq {
    1: double StartLat
    2: double StartLng
    3: double EndLat
    4: double EndLng
    5: string CarType
}
// 路线规划响应
struct RoutePlanResp {
    1: string Message
    2: i16 Code
    3: i32 Distance
    4: i32 Duration
    5: double Price
}

//微信登录请求
struct WechatLoginReq {
    1: i32 OpenID
    2: i32 UnionID
    3: string Nickname
    4: string Avatar
}

//微信登录响应
struct WechatLoginResp{
    1: string Message
    2: i64 Code
    3: Passengers passengers
    4: string Token
}

//绑定微信请求
struct BindWechatReq {
    1: i32 UserID
    2: i32 OpenID
    3: i32 UnionID
    4: string Nickname
    5: string Avatar
}

//绑定微信响应
struct BindWechatResp{
    1: string Message
    2: i64 Code
    3: Passengers passengers
}

//解绑微信请求
struct UnbindWechatReq {
    1: i32 UserID
}

//解绑微信响应
struct UnbindWechatResp{
    1: string Message
    2: i64 Code
}

//获取用户信息请求
struct GetUserInfoReq {
    1: i32 UserID
}

//获取用户信息响应
struct GetUserInfoResp{
    1: string Message
    2: i64 Code
    3: Passengers passengers
}

//更新用户信息请求
struct UpdateUserInfoReq {
    1: i32 UserID
    2: Passengers passengers
}

//更新用户信息响应
struct UpdateUserInfoResp{
    1: string Message
    2: i64 Code
    3: Passengers passengers
}

service PassengersServer {
    //短信验证码
    SendSmsResp SendSms (1:SendSmsReq req)
    //注册功能
    RegisterResp Register (1:RegisterReq req)
    //手机号登录
    LoginResp Login (1:LoginReq req)
    //实时定位
    LocationResp GetLocation(1: LocationReq req)
    //行程
    TripResp GetTripInfo(1: TripReq req)
    //计价
    PricingResp GetPricing(1: PricingReq req)
    //下单功能
    OrderResp PlaceOrder(1: OrderReq req)
    //路线规划
    RoutePlanResp RoutePlan(1: RoutePlanReq req)
    //微信登录
    WechatLoginResp WechatLogin (1:WechatLoginReq req)
    //绑定微信
    BindWechatResp BindWechat (1:BindWechatReq req)
    //解绑微信
    UnbindWechatResp UnbindWechat (1:UnbindWechatReq req)
    //获取用户信息
    GetUserInfoResp GetUserInfo (1:GetUserInfoReq req)
    //更新用户信息
    UpdateUserInfoResp UpdateUserInfo (1:UpdateUserInfoReq req)
}