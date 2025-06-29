package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"lxh_cx/config"
	"lxh_cx/kitex_gen/lxh_cx/base"
	order "lxh_cx/kitex_gen/lxh_cx/order"
	"lxh_cx/model"
	"net/http"
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct{}

// PassengerAddOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) PassengerAddOrder(ctx context.Context, req *order.PassengerAddOrderReq) (resp *order.PassengerAddOrderResp, err error) {
	//用户进行下单操作，但此时订单不是完整的订单，所以要将订单放进map中存入redis中的哈希列表中
	var user model.LxhPassenger
	err = config.Db.Where("id = ?", req.PassengerId).Find(&user).Limit(1).Error
	if err != nil {
		return &order.PassengerAddOrderResp{
			BaseResp: &base.BaseResp{
				Code: http.StatusBadRequest,
				Msg:  "user find failed",
			},
		}, err
	}

	var orderMap map[string]interface{}
	orderMap["passengerId"] = req.PassengerId
	orderMap["startAddr"] = req.StartAddr
	orderMap["endEnd"] = req.EndEnd
	orderMap["orderCode"] = uuid.NewString()
	orderMap["orderStatus"] = "1"
	orderMap["payStatus"] = "1"
	config.Rdb.HSet(ctx, "passenger"+string(req.PassengerId)+":"+"order", orderMap)

	return &order.PassengerAddOrderResp{
		BaseResp: &base.BaseResp{
			Code: http.StatusOK,
			Msg:  "order send orders success",
		},
	}, nil

}

// RouteAdd implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) RouteAdd(ctx context.Context, req *order.RouteAddReq) (resp *order.RouteAddResp, err error) {
	//路程的添加，因为路程是不规则的长度，所以将路程存入mongodb中，mysql不适合存储
	lxh := config.Mdb.Database("lxh").Collection("lxh_cx")
	Data := bson.M{
		"req": req.DataString,
	}

	one, err := lxh.InsertOne(context.Background(), Data)
	if err != nil {
		return &order.RouteAddResp{BaseResp: &base.BaseResp{
			Code: http.StatusBadRequest,
			Msg:  "route add failed",
		}}, err
	}

	fmt.Println(one.InsertedID)

	return &order.RouteAddResp{BaseResp: &base.BaseResp{
		Code: http.StatusOK,
		Msg:  "route add success",
	}}, nil

}
