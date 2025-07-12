package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"kitex_main/config"
	ride "kitex_main/kitex_gen/Ride_hailing_platform/ride"
	"kitex_main/model"
	"strconv"
	"time"
)

// TripsImpl implements the last service interface defined in the IDL.
type TripsImpl struct{}

// CreateTrips implements the TripsImpl interface.
// 创建行程
func (s *TripsImpl) CreateTrips(ctx context.Context, req *ride.CreateTripsReq) (resp *ride.CreateTripsResp, err error) {
	// TODO: Your code here...

	collection := config.Client.Database("lxh").Collection("lxh_car")

	Data := bson.M{
		"req": req.Data,
	}

	one, err := collection.InsertOne(config.Ctx, Data)

	fmt.Println(err)

	if err != nil {
		return &ride.CreateTripsResp{
			Code: 400,
		}, err
	}

	fmt.Println(one.InsertedID)

	return &ride.CreateTripsResp{Code: 200}, nil
}

// UpdatedTrips implements the TripsImpl interface.

func (s *TripsImpl) UpdatedTrips(ctx context.Context, req *ride.UpdatedTripsReq) (resp *ride.UpdatedTripsResp, err error) {
	// TODO: Your code here...

	var trips model.Trips

	err = trips.FindTripsId(int(req.TripsId))
	if err != nil {
		return &ride.UpdatedTripsResp{Code: 400}, err
	}

	now := time.Now().Unix()

	trips = model.Trips{
		EndTime: strconv.FormatInt(now, 10),
	}

	err = trips.UpdateTrip(int(req.TripsId))
	if err != nil {
		return &ride.UpdatedTripsResp{Code: 400}, err
	}
	return &ride.UpdatedTripsResp{Code: 200}, nil
}

// CreatedOrder implements the TripsImpl interface.
// 创建订单
func (s *TripsImpl) CreatedOrder(ctx context.Context, req *ride.CreatedOrderReq) (resp *ride.CreatedOrderResp, err error) {
	// TODO: Your code here...

	now := time.Now()

	Uuid := uuid.NewString()

	var Orders map[string]interface{}

	Orders = map[string]interface{}{
		"orderCode":   Uuid,
		"Amount":      req.Amount,
		"OrderStatus": "1",
		"PassengerId": req.PassengerId,
		"StartAddr":   req.StartAddr,
		"EndEnd":      req.EndEnd,
		"StartTime":   now,
	}

	err = config.RDB.HSet(config.Ctx, "Order", Orders).Err()
	if err != nil {
		return &ride.CreatedOrderResp{Code: 400}, err
	}

	//var order model.LxhOrder
	//
	//order = model.LxhOrder{
	//	OrderCode:   Uuid,
	//	Amount:      float64(req.Amount),
	//	OrderStatus: "1",
	//	PassengerId: int32(req.PassengerId),
	//	StartAddr:   req.StartAddr,
	//	EndEnd:      req.EndEnd,
	//	StartTime:   &now,
	//}
	//
	//err = order.CreateOrder()
	//
	//if err != nil {
	//	return &ride.CreatedOrderResp{Code: 400}, err
	//}

	return &ride.CreatedOrderResp{Code: 200}, nil
}
