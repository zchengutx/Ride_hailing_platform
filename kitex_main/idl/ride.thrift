namespace go Ride_hailing_platform.ride


struct CreateTripsReq{
    1:string data
}

struct CreateTripsResp{
    1:i64 code
}

struct UpdatedTripsReq{
    1:i64 tripsId
}

struct UpdatedTripsResp{
    i64 code = 1;
}

struct CreatedOrderReq{
    1:i64 Amount
    2:i64 PassengerId
    3:string StartAddr
    4:string EndEnd
}

struct CreatedOrderResp{
    i64 code = 1;
}

service Trips{
    CreateTripsResp CreateTrips (1:CreateTripsReq req)
    UpdatedTripsResp UpdatedTrips (1:UpdatedTripsReq req)
    CreatedOrderResp CreatedOrder (1:CreatedOrderReq req)
}

