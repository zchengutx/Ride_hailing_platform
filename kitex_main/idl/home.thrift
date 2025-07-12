namespace go Ride_hailing_platform.home

struct CreateHistoricalSearchReq{
    1:i64 userId
    2:string AddrName
    3:string ipAddrName
}

struct CreateHistoricalSearchResp{
    1:i64 code
}

struct HistoricalSearchListReq{
    1:i64 userId
    2:string ipAddrName
}

struct HistoricalSearchListResp{
    1:i64 code
    2:list <HistoricalSearchList> List
}

struct HistoricalSearchList{
    1:string AddrName
}

service Home{
    CreateHistoricalSearchResp CreateHistoricalSearch(1:CreateHistoricalSearchReq req)
    HistoricalSearchListResp HistoricalSearchList(1:HistoricalSearchListReq req)
    //司机端
    
}