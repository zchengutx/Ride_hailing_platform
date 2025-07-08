package utils

import (
	"fmt"
	"time"
)

// Redis哈希存储使用示例
func RedisHashExamples() {
	// 创建管理器实例
	hashManager := NewRedisHashManager()

	// 示例1: 用户会话信息存储
	ExampleUserSession(hashManager)

	// 示例2: 短信验证码管理
	ExampleSmsCodeManagement(hashManager)

	// 示例3: 用户缓存信息
	ExampleUserCache(hashManager)
}

// 示例1: 用户会话信息存储
func ExampleUserSession(manager *RedisHashManager) {
	fmt.Println("=== 用户会话信息存储示例 ===")

	userId := "12345"
	sessionKey := "user_session:" + userId

	// 存储用户会话信息
	sessionData := map[string]interface{}{
		"user_id":    userId,
		"username":   "张三",
		"login_time": time.Now().Unix(),
		"ip_address": "192.168.1.100",
		"device":     "mobile",
	}

	// 批量设置会话信息，30分钟过期
	err := manager.SetHashFields(sessionKey, sessionData, time.Minute*30)
	if err != nil {
		fmt.Printf("设置会话失败: %v\n", err)
		return
	}

	// 获取单个字段
	username, _ := manager.GetHashField(sessionKey, "username")
	fmt.Printf("用户名: %s\n", username)

	// 获取多个字段
	fields, _ := manager.GetHashFields(sessionKey, "username", "login_time", "device")
	fmt.Printf("用户信息: %v\n", fields)

	// 获取所有会话信息
	allFields, _ := manager.GetAllHashFields(sessionKey)
	fmt.Printf("完整会话信息: %v\n", allFields)
}

// 示例2: 短信验证码管理（已在handler中使用）
func ExampleSmsCodeManagement(manager *RedisHashManager) {
	fmt.Println("=== 短信验证码管理示例 ===")

	mobile := "13800138000"

	// 设置验证码
	code := 123456
	err := manager.SetSmsCode(mobile, code, time.Minute*5)
	if err != nil {
		fmt.Printf("设置验证码失败: %v\n", err)
		return
	}

	// 获取验证码
	storedCode, _ := manager.GetSmsCode(mobile)
	fmt.Printf("存储的验证码: %s\n", storedCode)

	// 增加尝试次数
	attempts, _ := manager.IncrementSmsAttempts(mobile)
	fmt.Printf("当前尝试次数: %d\n", attempts)

	// 获取尝试次数
	attemptsStr, _ := manager.GetSmsAttempts(mobile)
	fmt.Printf("尝试次数: %s\n", attemptsStr)
}

// 示例3: 用户缓存信息
func ExampleUserCache(manager *RedisHashManager) {
	fmt.Println("=== 用户缓存信息示例 ===")

	userId := "user:98765"

	// 设置用户基本信息
	userInfo := map[string]interface{}{
		"id":       98765,
		"nickname": "滴滴用户",
		"phone":    "13900139000",
		"level":    "VIP",
		"balance":  "199.50",
		"points":   1500,
	}

	// 缓存用户信息，1小时过期
	err := manager.SetHashFields(userId, userInfo, time.Hour)
	if err != nil {
		fmt.Printf("缓存用户信息失败: %v\n", err)
		return
	}

	// 增加用户积分
	newPoints, _ := manager.IncrementHashField(userId, "points", 100)
	fmt.Printf("新的积分: %d\n", newPoints)

	// 检查某个字段是否存在
	exists, _ := manager.HashFieldExists(userId, "balance")
	fmt.Printf("余额字段是否存在: %v\n", exists)

	// 获取哈希表字段数量
	length, _ := manager.HashLength(userId)
	fmt.Printf("用户信息字段数量: %d\n", length)

	// 删除特定字段
	err = manager.DeleteHashField(userId, "points")
	if err != nil {
		fmt.Printf("删除积分字段失败: %v\n", err)
	}
}

// 高级使用示例: 订单状态管理
func ExampleOrderStatusManagement(manager *RedisHashManager) {
	fmt.Println("=== 订单状态管理示例 ===")

	orderId := "order:202312150001"

	// 订单创建
	orderData := map[string]interface{}{
		"order_id":        "202312150001",
		"user_id":         "12345",
		"driver_id":       "",
		"status":          "waiting", // waiting, assigned, picked_up, completed, cancelled
		"created_time":    time.Now().Unix(),
		"start_location":  "北京市朝阳区",
		"end_location":    "北京市海淀区",
		"estimated_price": "25.80",
		"actual_price":    "",
	}

	// 存储订单信息，24小时过期
	err := manager.SetHashFields(orderId, orderData, time.Hour*24)
	if err != nil {
		fmt.Printf("创建订单失败: %v\n", err)
		return
	}

	// 更新订单状态 - 分配司机
	manager.SetHashField(orderId, "driver_id", "driver_67890", 0)
	manager.SetHashField(orderId, "status", "assigned", 0)
	manager.SetHashField(orderId, "assigned_time", time.Now().Unix(), 0)

	// 获取当前订单状态
	status, _ := manager.GetHashField(orderId, "status")
	fmt.Printf("当前订单状态: %s\n", status)

	// 批量获取订单关键信息
	keyFields, _ := manager.GetHashFields(orderId, "status", "driver_id", "estimated_price")
	fmt.Printf("订单关键信息: %v\n", keyFields)
}

// 使用redis哈希存储的优势说明
func RedisHashAdvantages() {
	fmt.Println(`
=== Redis哈希存储的优势 ===

1. 数据结构优势:
   - 一个key可以存储多个字段，减少key的数量
   - 支持原子性操作，如HINCRBY递增字段值
   - 内存使用更高效，特别是小的哈希表

2. 操作便利性:
   - 可以单独获取/设置某个字段，无需获取整个数据
   - 支持批量操作，如HMGET、HMSET
   - 可以检查字段是否存在

3. 业务场景适用:
   - 用户信息缓存（用户名、等级、积分等）
   - 会话管理（登录状态、权限等）
   - 配置管理（系统配置、用户偏好等）
   - 计数器管理（点赞数、访问量等）

4. 与简单键值存储的对比:
   原来: SET user:123:name "张三", SET user:123:level "VIP"  # 多个key
   现在: HSET user:123 name "张三" level "VIP"                # 一个hash key

5. 性能优势:
   - 减少网络往返次数
   - 降低内存碎片
   - 提高缓存命中率
`)
}
