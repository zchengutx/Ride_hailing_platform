namespace go cart.driver

// 司机申请注册请求
struct DriverPetitionReq {
    1: string name                    // 司机姓名
    2: string mobile                  // 联系电话
    3: string nick_name               // 昵称
    4: i64 car_age                    // 车龄
    5: string id_card_file_id         // 身份证文件ID
    6: string driver_license_file_id  // 驾照文件ID
    7: string driving_license_file_id // 行驶证文件ID
    8: string avatar_file_id          // 司机头像文件ID
}

// 司机申请注册响应
struct DriverPetitionResp {
    1: i16 code
    2: string message
    3: optional i64 application_id    // 申请ID
}

// 查询审核状态请求
struct CheckStatusReq {
    1: i64 driver_id                  // 司机ID
}

// 审核状态响应
struct CheckStatusResp {
    1: i16 code
    2: string message
    3: optional string check_status   // 审核状态：pending-待审核, approved-通过, rejected-拒绝
    4: optional string remark         // 审核备注
}

// 司机登录请求
struct DriverLoginReq {
    1: string mobile                  // 手机号
    2: string sms_code                // 短信验证码
}

// 司机登录响应
struct DriverLoginResp {
    1: i16 code
    2: string message
    3: optional i64 driver_id         // 司机ID
    4: optional string token          // 登录token
    5: optional DriverInfo driver_info // 司机信息
}

// 司机基本信息
struct DriverInfo {
    1: i64 id                         // 司机ID
    2: string name                    // 姓名
    3: string nick_name               // 昵称
    4: string mobile                  // 联系电话
    5: double all_account             // 总收益
    6: i64 car_age                    // 车龄
    7: string status                  // 接单状态：online-在线, offline-离线
    8: i64 file_id                    // 头像文件ID
}

// 获取司机信息请求
struct GetDriverInfoReq {
    1: i64 driver_id                  // 司机ID
}

// 获取司机信息响应
struct GetDriverInfoResp {
    1: i16 code
    2: string message
    3: optional DriverInfo driver_info
}

// 更新司机信息请求
struct UpdateDriverInfoReq {
    1: i64 driver_id                  // 司机ID
    2: optional string nick_name      // 昵称
    3: optional i64 file_id           // 头像文件ID
}

// 更新司机信息响应
struct UpdateDriverInfoResp {
    1: i16 code
    2: string message
}

// 司机上线/下线请求
struct ChangeStatusReq {
    1: i64 driver_id                  // 司机ID
    2: string status                  // 状态：online-上线, offline-下线
    3: optional string longitude      // 经度（上线时必填）
    4: optional string latitude       // 纬度（上线时必填）
}

// 司机上线/下线响应
struct ChangeStatusResp {
    1: i16 code
    2: string message
}

// 订单基本信息
struct OrderInfo {
    1: i64 id                         // 订单ID
    2: string order_code              // 订单编号
    3: double amount                  // 订单金额
    4: string order_status            // 订单状态
    5: i64 passenger_id               // 乘客ID
    6: string passenger_name          // 乘客姓名
    7: string passenger_mobile        // 乘客电话
    8: string start_addr              // 起始地址
    9: string end_addr                // 目的地址
    10: string start_time             // 开始时间
    11: string end_time               // 结束时间
    12: string order_type             // 订单类型：express-快车, booking-预约
}

// 获取待接订单请求
struct GetPendingOrdersReq {
    1: i64 driver_id                  // 司机ID
    2: string longitude               // 司机当前经度
    3: string latitude                // 司机当前纬度
    4: optional i32 radius            // 搜索半径（公里），默认5公里
}

// 获取待接订单响应
struct GetPendingOrdersResp {
    1: i16 code
    2: string message
    3: optional list<OrderInfo> orders // 可接订单列表
}

// 接单请求
struct AcceptOrderReq {
    1: i64 driver_id                  // 司机ID
    2: i64 order_id                   // 订单ID
}

// 接单响应
struct AcceptOrderResp {
    1: i16 code
    2: string message
    3: optional OrderInfo order_info  // 订单详情
}

// 开始行程请求
struct StartTripReq {
    1: i64 driver_id                  // 司机ID
    2: i64 order_id                   // 订单ID
    3: string longitude               // 当前经度
    4: string latitude                // 当前纬度
}

// 开始行程响应
struct StartTripResp {
    1: i16 code
    2: string message
}

// 完成订单请求
struct CompleteOrderReq {
    1: i64 driver_id                  // 司机ID
    2: i64 order_id                   // 订单ID
    3: string longitude               // 当前经度
    4: string latitude                // 当前纬度
    5: double actual_amount           // 实际费用
}

// 完成订单响应
struct CompleteOrderResp {
    1: i16 code
    2: string message
    3: optional double driver_income  // 司机收入
}

// 取消订单请求
struct CancelOrderReq {
    1: i64 driver_id                  // 司机ID
    2: i64 order_id                   // 订单ID
    3: string cancel_reason           // 取消原因
    4: optional string cancel_remark  // 取消备注
}

// 取消订单响应
struct CancelOrderResp {
    1: i16 code
    2: string message
}

// 收益明细
struct IncomeDetail {
    1: string date                    // 日期
    2: i32 order_count               // 订单数量
    3: double total_income           // 总收入
    4: double platform_fee           // 平台费用
    5: double net_income             // 净收入
}

// 查询收益请求
struct GetIncomeReq {
    1: i64 driver_id                  // 司机ID
    2: string start_date              // 开始日期 YYYY-MM-DD
    3: string end_date                // 结束日期 YYYY-MM-DD
}

// 查询收益响应
struct GetIncomeResp {
    1: i16 code
    2: string message
    3: optional list<IncomeDetail> income_details // 收益明细列表
    4: optional double total_income               // 总收益
}

// 更新位置请求
struct UpdateLocationReq {
    1: i64 driver_id                  // 司机ID
    2: string longitude               // 经度
    3: string latitude                // 纬度
    4: optional double speed          // 行驶速度
    5: optional double direction      // 行驶方向
}

// 更新位置响应
struct UpdateLocationResp {
    1: i16 code
    2: string message
}
service DriverService {
    // 司机注册申请
    DriverPetitionResp DriverPetition(1: DriverPetitionReq req)
    
    // 查询审核状态
    CheckStatusResp CheckStatus(1: CheckStatusReq req)
    
    // 司机登录
    DriverLoginResp DriverLogin(1: DriverLoginReq req)
    
    // 获取司机信息
    GetDriverInfoResp GetDriverInfo(1: GetDriverInfoReq req)
    
    // 更新司机信息
    UpdateDriverInfoResp UpdateDriverInfo(1: UpdateDriverInfoReq req)
    
    // 司机上线/下线
    ChangeStatusResp ChangeStatus(1: ChangeStatusReq req)
    
    // 获取待接订单
    GetPendingOrdersResp GetPendingOrders(1: GetPendingOrdersReq req)
    
    // 接单
    AcceptOrderResp AcceptOrder(1: AcceptOrderReq req)
    
    // 开始行程
    StartTripResp StartTrip(1: StartTripReq req)
    
    // 完成订单
    CompleteOrderResp CompleteOrder(1: CompleteOrderReq req)
    
    // 取消订单
    CancelOrderResp CancelOrder(1: CancelOrderReq req)
    
    // 查询收益
    GetIncomeResp GetIncome(1: GetIncomeReq req)
    
    // 更新位置
    UpdateLocationResp UpdateLocation(1: UpdateLocationReq req)
}