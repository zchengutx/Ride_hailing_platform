package main

import (
	"context"
	"kitex_main/config"
	home "kitex_main/kitex_gen/Ride_hailing_platform/home"
	"strconv"
)

// HomeImpl implements the last service interface defined in the IDL.
type HomeImpl struct{}

// CreateHistoricalSearch implements the HomeImpl interface.
// 创建历史记录
func (s *HomeImpl) CreateHistoricalSearch(ctx context.Context, req *home.CreateHistoricalSearchReq) (resp *home.CreateHistoricalSearchResp, err error) {
	// TODO: Your code here...

	itoa := strconv.Itoa(int(req.UserId))

	err = config.RDB.ZIncrBy(config.Ctx, itoa+req.IpAddrName, 1, req.AddrName).Err()
	if err != nil {
		return &home.CreateHistoricalSearchResp{Code: 400}, err
	}

	return &home.CreateHistoricalSearchResp{Code: 200}, nil
}

// HistoricalSearchList implements the HomeImpl interface.
// 历史记录展示
func (s *HomeImpl) HistoricalSearchList(ctx context.Context, req *home.HistoricalSearchListReq) (resp *home.HistoricalSearchListResp, err error) {
	// TODO: Your code here...

	itoa := strconv.Itoa(int(req.UserId))

	val := config.RDB.ZRevRange(config.Ctx, itoa+req.IpAddrName, 0, 9).Val()

	var list []*home.HistoricalSearchList

	for _, v := range val {
		data := home.HistoricalSearchList{
			AddrName: v,
		}
		list = append(list, &data)
	}

	return &home.HistoricalSearchListResp{
		Code: 200,
		List: list,
	}, nil
}
