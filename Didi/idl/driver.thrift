namespace go Didi.driver
struct CallACarReq{ //司机申请
    1: i64 userId
    2: string address
    3: string drivingLicenseNumber
    4: string quasiDrivingType
    5: i64 drivingAge
    6: string carNum
    7: string carType
    8: string vehicleMileage
    9: string servingTheCity
}
struct CallACarResp{
    1: i64 code
    2: string message
}
struct DriverAuditReq{ //司机审核
    1: i64 driverId
    2: string auditStatus
}
struct DriverAuditResp{
    1: i64 code
    2: string message
}
struct AddDriverReq{ //司机添加
    1: i64 driverId
}
struct AddDriverResp{
    1: i64 code
    2: string message
}
struct DriverOnlineReq{ //司机状态
    1: i64 driverId
    2: string driverStatus
}
struct DriverOnlineResp{
    1: i64 code
    2: string message
}
struct ReceivingOrderReq{ //司机接单
    1: i64 driverId
    2: string orderId
    3: double currentLatitude
    4: double currentLongitude
}
struct ReceivingOrderResp{
    1: i64 code
    2: string message
    3: string orderDetails
}

service driverServer{
   CallACarResp CallACar(1: CallACarReq req)
   DriverAuditResp DriverAudit(1: DriverAuditReq req)
   AddDriverResp AddDriver(1: AddDriverReq req)
   DriverOnlineResp DriverOnline(1: DriverOnlineReq req)
   ReceivingOrderResp ReceivingOrder(1: ReceivingOrderReq req)
}