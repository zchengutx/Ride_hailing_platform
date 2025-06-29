namespace go kitexTwo.user
include "base.thrift"

struct CreateUserReq {
    1: required string title
    2: required string SendSmsCode
}
struct CreateUserResp{
    1: string message
    2: i64 code
}
struct LoginUserReq{
    1: required string title
    2: required string SendSmsCode
}
struct LoginUserResp{
    1: string message
    2: i64 code
    3: i64 Uid
}
service User {
    CreateUserResp CreateUser(1: CreateUserReq req)
    LoginUserResp LoginUser(1: LoginUserReq req)
}