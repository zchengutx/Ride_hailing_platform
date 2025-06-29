namespace go lxh_cx.driver

include "base.thrift"

struct DriverRegisterReq {
    1: string idCardFileId
    2: string driverLicenseFileId
    3: string drivingLicenseFileId
    4: string avatorFileId
    5: i32 DriverId
}

struct DriverRegisterResp {
    1: base.BaseResp baseResp
}

struct DriverAddReq {
    1: string idCardFileId
    2: string driverLicenseFileId
    3: string drivingLicenseFileId
    4: string avatorFileId
    5: i32 DriverId
    6: i32 passengerId
}

struct DriverAddResp {
    1: base.BaseResp baseResp
}

struct DriverOverOrderReq {
    1: i32 orderId
}

struct DriverOverOrderResp {
    1: base.BaseResp baseResp
}

struct DriverCancelOrderReq {
    1: i32 driverId
    2: i32 confirmPerson
    3: string confirmReason
    4: string confirmRemark
    5: i32 orderId
}

struct DriverCancelOrderResp {
    1: base.BaseResp baseResp
}

struct DriverInfoListReq {
    1: i32 driverId
}

struct DriverInfoListResp {
    1: base.BaseResp baseResp
    2: DriverInfo driverInfo
}
struct DriverInfo {
    1: string name
    2: string nickName
    3: double allAcount
    4: string carAge
    5: string handlerFile
}

service DriverService {
    DriverRegisterResp DriverRegister(1: DriverRegisterReq req)
    DriverAddResp DriverAdd(1: DriverAddReq req)
    DriverOverOrderResp DriverOverOrder(1: DriverOverOrderReq req)
    DriverCancelOrderResp DriverCancelOrder(1: DriverCancelOrderReq req)
    DriverInfoListResp DriverInfoList(1: DriverInfoListReq req)
}