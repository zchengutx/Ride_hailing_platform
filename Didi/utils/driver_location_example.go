package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// 司机位置管理和订单推送的完整使用示例

// SimulateDriversOnline 模拟司机上线（供测试使用）
func SimulateDriversOnline() {
	geoManager := NewRedisGeoManager()

	// 模拟一些司机上线并更新位置
	drivers := []DriverLocation{
		{
			DriverID:  "driver_001",
			Latitude:  39.9042, // 天安门附近
			Longitude: 116.4074,
			Status:    STATUS_ONLINE,
			CarType:   "快车",
			UpdatedAt: time.Now().Unix(),
		},
		{
			DriverID:  "driver_002",
			Latitude:  39.9100, // 稍微北一点
			Longitude: 116.4050,
			Status:    STATUS_ONLINE,
			CarType:   "专车",
			UpdatedAt: time.Now().Unix(),
		},
		{
			DriverID:  "driver_003",
			Latitude:  39.9000, // 稍微南一点
			Longitude: 116.4100,
			Status:    STATUS_ONLINE,
			CarType:   "快车",
			UpdatedAt: time.Now().Unix(),
		},
		{
			DriverID:  "driver_004",
			Latitude:  39.8950, // 更远一点，但仍在3公里内
			Longitude: 116.4150,
			Status:    STATUS_BUSY, // 忙碌状态，不会收到订单
			CarType:   "快车",
			UpdatedAt: time.Now().Unix(),
		},
	}

	// 批量更新司机位置
	for _, driver := range drivers {
		if err := geoManager.UpdateDriverLocation(driver); err != nil {
			log.Printf("更新司机 %s 位置失败: %v", driver.DriverID, err)
		} else {
			log.Printf("司机 %s 上线成功，位置: %.6f,%.6f",
				driver.DriverID, driver.Latitude, driver.Longitude)
		}
	}
}

// TestOrderNotification 测试订单推送功能
func TestOrderNotification() {
	// 1. 模拟司机上线
	fmt.Println("=== 模拟司机上线 ===")
	SimulateDriversOnline()

	// 2. 模拟用户下单
	fmt.Println("\n=== 模拟用户下单 ===")
	testOrder := OrderMessage{
		OrderID:       "TEST_" + GenerateOrderId(),
		UserID:        "12345",
		StartLocation: "天安门广场",
		EndLocation:   "王府井大街",
		CartType:      "快车",
		CreatedAt:     time.Now().Unix(),
		Latitude:      39.9042,
		Longitude:     116.4074,
	}

	// 3. 查找周边司机
	geoManager := NewRedisGeoManager()
	nearbyDrivers, err := geoManager.GetNearbyDrivers(
		testOrder.Latitude,
		testOrder.Longitude,
		3.0,
		testOrder.CartType,
	)

	if err != nil {
		log.Printf("查找周边司机失败: %v", err)
		return
	}

	fmt.Printf("找到 %d 个周边司机:\n", len(nearbyDrivers))
	for _, driver := range nearbyDrivers {
		fmt.Printf("- 司机 %s，距离: %.0f 米\n", driver.DriverID, driver.Distance)
		testOrder.DriverIDs = append(testOrder.DriverIDs, driver.DriverID)
	}

	// 4. 推送订单消息
	fmt.Println("\n=== 推送订单消息 ===")
	rabbitMQManager := GetRabbitMQManager()
	if rabbitMQManager == nil {
		fmt.Println("RabbitMQ未初始化，无法推送消息")
		return
	}

	if err := rabbitMQManager.PublishOrderToNearbyDrivers(testOrder); err != nil {
		log.Printf("推送订单失败: %v", err)
	} else {
		fmt.Printf("订单 %s 推送成功！\n", testOrder.OrderID)
	}
}

// DriverMessageHandler 司机端消息处理示例
func DriverMessageHandler(driverID string) func([]byte) error {
	return func(msgBody []byte) error {
		var orderMsg OrderMessage
		if err := json.Unmarshal(msgBody, &orderMsg); err != nil {
			return fmt.Errorf("解析订单消息失败: %v", err)
		}

		// 检查这个订单是否推送给了当前司机
		isTargetDriver := false
		for _, targetDriverID := range orderMsg.DriverIDs {
			if targetDriverID == driverID {
				isTargetDriver = true
				break
			}
		}

		if isTargetDriver {
			fmt.Printf("\n[司机 %s] 收到新订单通知:\n", driverID)
			fmt.Printf("  订单号: %s\n", orderMsg.OrderID)
			fmt.Printf("  起点: %s\n", orderMsg.StartLocation)
			fmt.Printf("  终点: %s\n", orderMsg.EndLocation)
			fmt.Printf("  车型: %s\n", orderMsg.CartType)
			fmt.Printf("  距离: 计算中...\n")

			// 这里可以添加司机接单逻辑
			// 比如自动接单、显示接单界面等
		}

		return nil
	}
}

// StartDriverListening 启动司机消息监听（模拟司机端）
func StartDriverListening(driverID string) error {
	rabbitMQManager := GetRabbitMQManager()
	if rabbitMQManager == nil {
		return fmt.Errorf("RabbitMQ未初始化")
	}

	handler := DriverMessageHandler(driverID)
	return rabbitMQManager.ConsumeDriverMessages(driverID, handler)
}

// GetDriverLocationStatus 获取司机位置状态报告
func GetDriverLocationStatus() {
	geoManager := NewRedisGeoManager()

	fmt.Println("\n=== 司机位置状态报告 ===")

	// 获取在线司机数量
	onlineCount, err := geoManager.GetOnlineDriverCount()
	if err != nil {
		fmt.Printf("获取在线司机数量失败: %v\n", err)
	} else {
		fmt.Printf("当前在线司机数量: %d\n", onlineCount)
	}

	// 获取天安门周边5公里的司机
	driversInArea, err := geoManager.GetDriversInArea(39.9042, 116.4074, 5.0)
	if err != nil {
		fmt.Printf("获取区域司机失败: %v\n", err)
	} else {
		fmt.Printf("天安门周边5公里司机数量: %d\n", len(driversInArea))
		for i, driverID := range driversInArea {
			if i < 5 { // 只显示前5个
				location, err := geoManager.GetDriverLocation(driverID)
				if err == nil {
					fmt.Printf("  - %s: %.6f,%.6f (%s)\n",
						driverID, location.Latitude, location.Longitude, location.Status)
				}
			}
		}
		if len(driversInArea) > 5 {
			fmt.Printf("  ... 还有 %d 个司机\n", len(driversInArea)-5)
		}
	}
}

// CleanupTestData 清理测试数据
func CleanupTestData() {
	geoManager := NewRedisGeoManager()

	testDrivers := []string{"driver_001", "driver_002", "driver_003", "driver_004"}

	fmt.Println("\n=== 清理测试数据 ===")
	for _, driverID := range testDrivers {
		if err := geoManager.RemoveDriver(driverID); err != nil {
			fmt.Printf("移除司机 %s 失败: %v\n", driverID, err)
		} else {
			fmt.Printf("司机 %s 已下线\n", driverID)
		}
	}
}

// 完整的功能演示
func DemoCompleteWorkflow() {
	fmt.Println("🚗 滴滴出行 - 订单推送系统演示")
	fmt.Println("=====================================")

	// 1. 司机上线
	SimulateDriversOnline()
	time.Sleep(1 * time.Second)

	// 2. 查看状态
	GetDriverLocationStatus()
	time.Sleep(1 * time.Second)

	// 3. 测试订单推送
	TestOrderNotification()
	time.Sleep(2 * time.Second)

	// 4. 清理数据
	CleanupTestData()

	fmt.Println("\n✅ 演示完成！")
}
