package utils

import (
	"cart/biz/dal/global"
	driver "cart/kitex_gen/cart/driver/driverservice"
	mapservice "cart/kitex_gen/cart/mapservice/mapservice"
	passenger "cart/kitex_gen/cart/passenger/passengerservice"
	wechat "cart/kitex_gen/cart/wechat/wechatservice"
	"fmt"
	"sync"

	"github.com/cloudwego/kitex/client"
)

type ClientManager struct {
	passengerClient passenger.Client
	passengerOnce   sync.Once
	driverClient    driver.Client
	driverOnce      sync.Once
	mapClient       mapservice.Client
	mapOnce         sync.Once
	wechatClient    wechat.Client
	wechatOnce      sync.Once
}

var (
	clientManager = &ClientManager{}
)

func (c *ClientManager) GetPassengerClient() passenger.Client {
	c.passengerOnce.Do(func() {
		passengerConf := &global.AppConf.PassengerClient
		c.passengerClient, _ = passenger.NewClient("passenger", client.WithHostPorts(fmt.Sprintf("%v:%v", passengerConf.Host, passengerConf.Port)))
	})
	return c.passengerClient
}

func GetDefaultPassengerClient() passenger.Client {
	return clientManager.GetPassengerClient()
}

func (c *ClientManager) GetDriverClient() driver.Client {
	c.driverOnce.Do(func() {
		driverConf := &global.AppConf.DriverClient
		c.driverClient, _ = driver.NewClient("driver", client.WithHostPorts(fmt.Sprintf("%v:%v", driverConf.Host, driverConf.Port)))
	})
	return c.driverClient
}

func GetDefaultDriverClient() driver.Client {
	return clientManager.GetDriverClient()
}

func (c *ClientManager) GetMapClient() mapservice.Client {
	c.mapOnce.Do(func() {
		mapConf := &global.AppConf.MapClient
		c.mapClient, _ = mapservice.NewClient("map", client.WithHostPorts(fmt.Sprintf("%v:%v", mapConf.Host, mapConf.Port)))
	})
	return c.mapClient
}

func GetDefaultMapClient() mapservice.Client {
	return clientManager.GetMapClient()
}

func (c *ClientManager) GetWeChatClient() wechat.Client {
	c.wechatOnce.Do(func() {
		wechatConf := &global.AppConf.WeChatClient
		c.wechatClient, _ = wechat.NewClient("wechat", client.WithHostPorts(fmt.Sprintf("%v:%v", wechatConf.Host, wechatConf.Port)))
	})
	return c.wechatClient
}

func GetDefaultWeChatClient() wechat.Client {
	return clientManager.GetWeChatClient()
}
