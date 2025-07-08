package utils

import (
	"Didi/biz/dal/global"
	driver "Didi/kitex_gen/Didi/driver/driverserver"
	user "Didi/kitex_gen/Didi/user/userserver"
	"fmt"
	"sync"

	"github.com/cloudwego/kitex/client"
)

// ClientManager 客户端管理器，用于统一管理RPC客户端实例
// 采用单例模式，确保每个客户端只创建一次，提高性能和资源利用率
type ClientManager struct {
	userClient   user.Client   // 用户服务客户端实例
	driverClient driver.Client // 司机服务客户端实例
	userOnce     sync.Once     // 用户客户端初始化控制器，保证只初始化一次
	driverOnce   sync.Once     // 司机客户端初始化控制器，保证只初始化一次
}

var (
	// clientManager 全局客户端管理器实例
	clientManager = &ClientManager{}
)

// GetUserClient 获取用户客户端实例（单例模式）
// 使用sync.Once确保客户端只初始化一次，提高性能并避免资源浪费
// 返回用户服务的RPC客户端实例
func (cm *ClientManager) GetUserClient() user.Client {
	cm.userOnce.Do(func() {
		// 从全局配置中获取用户服务的连接配置
		userConf := &global.AppConf.UserClient
		// 创建用户服务客户端，配置主机和端口
		cm.userClient, _ = user.NewClient("user",
			client.WithHostPorts(fmt.Sprintf("%v:%v", userConf.Host, userConf.Port)))
	})
	return cm.userClient
}

// GetDriverClient 获取司机客户端实例（单例模式）
// 使用sync.Once确保客户端只初始化一次，提高性能并避免资源浪费
// 返回司机服务的RPC客户端实例
func (cm *ClientManager) GetDriverClient() driver.Client {
	cm.driverOnce.Do(func() {
		// 从全局配置中获取司机服务的连接配置
		driverConf := &global.AppConf.DriverClient
		// 创建司机服务客户端，配置主机和端口
		cm.driverClient, _ = driver.NewClient("driver",
			client.WithHostPorts(fmt.Sprintf("%v:%v", driverConf.Host, driverConf.Port)))
	})
	return cm.driverClient
}

// GetDefaultUserClient 获取默认用户客户端
// 这是一个便捷方法，直接返回全局客户端管理器中的用户客户端实例
// 在业务层可以直接调用此方法获取用户服务客户端
func GetDefaultUserClient() user.Client {
	return clientManager.GetUserClient()
}

// GetDefaultDriverClient 获取默认司机客户端
// 这是一个便捷方法，直接返回全局客户端管理器中的司机客户端实例
// 在业务层可以直接调用此方法获取司机服务客户端
func GetDefaultDriverClient() driver.Client {
	return clientManager.GetDriverClient()
}
