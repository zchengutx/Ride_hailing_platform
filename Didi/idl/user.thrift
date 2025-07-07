namespace go Didi.user
struct SendSmsReq{ //短信
    1: string mobile
    2: string sendSmsCode
}
struct SendSmsResp{
    1: i64 code
    2: string message
}
struct LoginUserReq{ //注册登录一体化
    1: string mobile
    2: string sendSmsCode
}
struct LoginUserResp{
    1: i64 code
    2: string message
    3: i64 UId
}
struct RealNameReq{ //实名认证
    1: i64 UId
    2: string userName
    3: string sex
    4: i64 age
    5: string idCard
}
struct RealNameResp{
    1: i64 code
    2: string message
}
service UserServer{
    SendSmsResp SendSms(1: SendSmsReq req)
    LoginUserResp LoginUser(1: LoginUserReq req)
    RealNameResp RealName(1: RealNameReq req)
}