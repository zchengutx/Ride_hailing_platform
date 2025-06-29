package api

import (
	"github.com/cloudwego/kitex/client"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"kitexTwo/client/basic/config"
	"kitexTwo/client/handler/request"
	pb2 "kitexTwo/kitex_gen/kitexTwo/user"
	pb "kitexTwo/kitex_gen/kitexTwo/user/user"
	"kitexTwo/utlis"
	"log"
	"time"
)

func CreateUser(c *gin.Context) {
	var req request.CreateUserReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "<查询失败>",
			"data": err.Error(),
		})
		return
	}
	cli, err := pb.NewClient("user", client.WithHostPorts("127.0.0.1:50051"))
	if err != nil {
		log.Fatal(err)
	}

	user, _ := cli.CreateUser(c, &pb2.CreateUserReq{
		Title:       req.Mobile,
		SendSmsCode: req.SendSmsCode,
	})
	c.JSON(200, user)
}

func LoginUser(c *gin.Context) {
	var req request.LoginUserReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "<查询失败>",
			"data": err.Error(),
		})
		return
	}
	cli, err := pb.NewClient("user", client.WithHostPorts("127.0.0.1:50051"))
	if err != nil {
		log.Fatal(err)
	}
	user, _ := cli.LoginUser(c, &pb2.LoginUserReq{
		Title:       req.Mobile,
		SendSmsCode: req.SendSmsCode,
	})
	if user.Code != 200 {
		c.JSON(200, user)
		return
	}
	token, _ := utlis.NewJWT(config.JWT_SELECT_KEY).CreateToken(utlis.CustomClaims{
		ID: uint(user.Uid),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 86400,
		},
	})
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": map[string]string{
			"token": token,
		},
	})
}
