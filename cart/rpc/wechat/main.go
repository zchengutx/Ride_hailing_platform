package main

import (
	wechat "cart/kitex_gen/cart/wechat/wechatservice"
	"cart/rpc/basic/global"
	_ "cart/rpc/basic/inits"
	"fmt"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
)

func main() {
	wechatConf := &global.AppConf.WeChatService
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", wechatConf.Host, wechatConf.Port))
	svr := wechat.NewServer(new(WeChatServiceImpl), server.WithServiceAddr(addr))
	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
