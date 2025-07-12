namespace go Ride_hailing_platform.user


struct SendSmsReq{
    1:string mobile
    2:string source
}

struct SendSmsResp {
    1:i64 code
    2:string msg
}

struct LoginReq{
    1:string mobile
    2:string sendSmsCode
}

struct LoginResp {
    1:i64 code
    2:string msg
    3:i64 token
}

struct CallBackReq{
    1:string openId
    2:string nickName
}

struct CallBackResp {
    1:i64 code
    2:string msg
    3:i64 token
}

struct InfoUserReq{
    1:i64 userId
}

struct InfoUserResp {
    1:i64 code
    2:InfoUser Info
}

struct InfoUser{
    1:string userName
    2:string mobile
    3:string nike_name
    4:string sex
    5:string mileage
    6:string age
}

struct DriverLoginReq {
    1:string mobile
    2:string sendSmsCode
}

struct DriverLoginResp{
    1:i64 code
}

struct DriverRealNameReq {
    1:string Name
    2:string IdCard
}

struct DriverRealNameResp{
    1:i64 code
}


service UserServer{
    SendSmsResp SendSms(1:SendSmsReq req)
    LoginResp Login(1:LoginReq req)
    CallBackResp CallBack(1:CallBackReq req)
    InfoUserResp InfoUser(1:InfoUserReq req)
    //司机端
    DriverLoginResp DriverLogin(1:DriverLoginReq req)  //登录
    DriverRealNameResp DriverRealName(1:DriverRealNameReq req)  //实名


}