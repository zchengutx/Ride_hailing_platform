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
service driverServer{
   CallACarResp CallACar(1: CallACarReq req)
}