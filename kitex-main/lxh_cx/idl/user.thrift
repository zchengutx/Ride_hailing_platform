namespace go lxh_cx.user

include "base.thrift"

struct SendSmsReq {
    1: string mobile
    2: string source
}

struct SendSmsResp {
    1: base.BaseResp baseResp
}

struct RegisterReq {
    1: string mobile
    2: string sendCode
}

struct RegisterResp{
    1: base.BaseResp baseResp
}

struct LoginReq {
    1: string mobile
    2: string sendCode
}

struct LoginResp{
    1: base.BaseResp baseResp
    2: i64 id
}

struct UserInfoListReq {
    1: i64 id
}

struct UserInfoListResp{
    1: base.BaseResp baseResp
    2: UserInfoResp userInfo
}
struct UserInfoResp{
    1:string nickName
    2:string sex
    3:string mileage
}

struct BindMobileReq {
    1: i64 id
    2: string mobile
}

struct BindMobileResp{
    1: base.BaseResp baseResp
}

struct UpdateUserInfoReq {
    1: i64 id
    2: string mobile
    3: string nickName
    4: string sex
}

struct UpdateUserInfoResp{
    1: base.BaseResp baseResp
}

struct UpdateCancelOrderReq {
    1: i32 orderId
    2: i32 confirmPerson
    3: string confirmReason
    4: i32 id
    5: string confirmRemark
}

struct UpdateCancelOrderResp{
    1: base.BaseResp baseResp
}

struct FindRouteRecordReq {
    1: i32 userId
}

struct FindRouteRecordResp{
    1: base.BaseResp baseResp
    2: list <FindRouteRecord> route
}
struct FindRouteRecord {
    1: string address
}

struct SetRouteUserReq {
    1: i32 userId
    2: string startAddr
    3: string endEnd
}

struct SetRouteUserResp{
    1: base.BaseResp baseResp
}

service UserService {
    SendSmsResp SendSms(1: SendSmsReq req)
    RegisterResp Register(1: RegisterReq req)
    LoginResp Login(1: LoginReq req)
    UserInfoListResp UserInfoList(1: UserInfoListReq req)
    BindMobileResp BindMobile(1: BindMobileReq req)
    UpdateUserInfoResp UpdateUserInfo(1: UpdateUserInfoReq req)
    UpdateCancelOrderResp UpdateCancelOrder(1: UpdateCancelOrderReq req)
    FindRouteRecordResp FindRouteRecord(1: FindRouteRecordReq req)
    SetRouteUserResp SetRouteUser(1: SetRouteUserReq req)
}