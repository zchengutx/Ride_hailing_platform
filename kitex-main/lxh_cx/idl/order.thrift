namespace go lxh_cx.order

include "base.thrift"

struct PassengerAddOrderReq {
    1: i32 passengerId
    2: string startAddr
    3: string endEnd
}

struct PassengerAddOrderResp {
    1: base.BaseResp baseResp
}

struct RouteAddReq {
    1: string dataString
}

struct RouteAddResp {
    1: base.BaseResp baseResp
}

service OrderService {
    PassengerAddOrderResp PassengerAddOrder(1: PassengerAddOrderReq req)
    RouteAddResp RouteAdd(1: RouteAddReq req)
}