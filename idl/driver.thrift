namespace go driver


struct LxhDriver {
    1: i32 Id
    2: string Name
    3: string NickName
    4: double AllAcount
    5: i32    CarAge
    6: string Status
    7: string Mobile
    8: i32    FileId
}
//短信验证码请求
struct SendSmsDriverReq{
    1: string Mobile
    2: string Source
}
//短信验证码响应
struct SendSmsDriverResp{
    1: string Message
    2: i64 Code
}

//司机注册请求
struct DriverRegisterReq{
    1: string Mobile
    2: string SmsCode
    3: string Name
}
//司机注册响应
struct DriverRegisterResp{
    1: string Message
    2: i64 Code
}

//获取司机信息请求
struct GetDriverInfoReq{
    1: i32 DriverID
}
//获取司机信息响应
struct GetDriverInfoResp{
    1: string Message
    2: i64 Code
    3: LxhDriver Driver
}

//更新司机状态请求
struct UpdateDriverStatusReq {
    1: i32    DriverID
    2: i16    WorkStatus
}
//更新司机状态响应
struct UpdateDriverStatusResp {
    1: string Message
    2: i16 Code
}

//获取待接订单请求
struct GetAwaitOrderReq {
    1: i32 DriverID
    2: i32    Distance  //距离
}
//获取待接订单响应
struct GetAwaitOrderResp{
    1: string Message
    2: i16 Code
}

//接单请求
struct AccessOrderReq{
    1: i32 OrderID
    2: i32 DriverID
}
//接单响应
struct AccessOrderResp{
    1: string Message
    2: i16 Code
}

//拒接请求
struct RejectionOrderReq{
    1: i32 OrderID
    2: i32 DriverID
}
//拒接响应
struct RejectionOrderResp{
    1: string Message
    2: i16    Code
}

//开启行程请求
struct StartTripReq{
    1: i32 DriverID
    2: i32 OrderID
}
//开启行程响应
struct StartTripResp{
    1: string Message
    2: i16    Code
}

//结束行程请求
struct EndTripReq{
    1: i32 DriverID
    2: i32 OrderID
}
//结束行程响应
struct EndTripResp{
    1: string Message
    2: i16    Code
}

//到达乘客位置请求
struct ArriveLocationReq{
    1: i32 DriverID
    2: i32 OrderID
}
//到达乘客位置响应
struct ArriveLocationResp{
    1: string Message
    2: i16    Code
}
service DriverServer {
    //短信验证码
    SendSmsDriverResp SendSms (1: SendSmsDriverReq req)
    //司机注册功能
    DriverRegisterResp Register (1:DriverRegisterReq req)
    //获取司机信息
    GetDriverInfoResp GetDriverInfo (1:GetDriverInfoReq req)
    //更新司机状态
    UpdateDriverStatusResp UpdateDriverStatus (1:UpdateDriverStatusReq req)
    //获取待接订单功能
    GetAwaitOrderResp GetAwaitOrder (1:GetAwaitOrderReq req)
    //接单功能
    AccessOrderResp AccessOrder (1:AccessOrderReq req)
    //拒接功能
    RejectionOrderResp RejectionOrder (1:RejectionOrderReq req)
    //到达乘客位置功能
    ArriveLocationResp ArriveLocation (1:ArriveLocationReq req)
    //开启行程功能
    StartTripResp StartTrip (1:StartTripReq req)
    //结束行程功能
    EndTripResp EndTrip (1:EndTripReq req)
    //
}