package utils

import (
	"Didi/rpc/common/global"
	"context"
	"time"
)

// RedisHashManager Redis哈希存储管理器
type RedisHashManager struct {
	ctx context.Context
}

// NewRedisHashManager 创建Redis哈希管理器实例
func NewRedisHashManager() *RedisHashManager {
	return &RedisHashManager{
		ctx: context.Background(),
	}
}

// SetHashField 设置哈希字段
// key: 哈希表的键名
// field: 字段名
// value: 字段值
// expiration: 过期时间
func (r *RedisHashManager) SetHashField(key, field string, value interface{}, expiration time.Duration) error {
	// 设置哈希字段
	if err := global.Rdb.HSet(r.ctx, key, field, value).Err(); err != nil {
		return err
	}

	// 设置过期时间
	if expiration > 0 {
		return global.Rdb.Expire(r.ctx, key, expiration).Err()
	}

	return nil
}

// GetHashField 获取哈希字段值
func (r *RedisHashManager) GetHashField(key, field string) (string, error) {
	return global.Rdb.HGet(r.ctx, key, field).Result()
}

// SetHashFields 批量设置哈希字段
func (r *RedisHashManager) SetHashFields(key string, fields map[string]interface{}, expiration time.Duration) error {
	// 批量设置哈希字段
	if err := global.Rdb.HMSet(r.ctx, key, fields).Err(); err != nil {
		return err
	}

	// 设置过期时间
	if expiration > 0 {
		return global.Rdb.Expire(r.ctx, key, expiration).Err()
	}

	return nil
}

// GetHashFields 批量获取哈希字段
func (r *RedisHashManager) GetHashFields(key string, fields ...string) ([]interface{}, error) {
	return global.Rdb.HMGet(r.ctx, key, fields...).Result()
}

// GetAllHashFields 获取哈希表中所有字段和值
func (r *RedisHashManager) GetAllHashFields(key string) (map[string]string, error) {
	return global.Rdb.HGetAll(r.ctx, key).Result()
}

// DeleteHashField 删除哈希字段
func (r *RedisHashManager) DeleteHashField(key string, fields ...string) error {
	return global.Rdb.HDel(r.ctx, key, fields...).Err()
}

// DeleteHash 删除整个哈希表
func (r *RedisHashManager) DeleteHash(key string) error {
	return global.Rdb.Del(r.ctx, key).Err()
}

// HashFieldExists 检查哈希字段是否存在
func (r *RedisHashManager) HashFieldExists(key, field string) (bool, error) {
	return global.Rdb.HExists(r.ctx, key, field).Result()
}

// HashLength 获取哈希表字段数量
func (r *RedisHashManager) HashLength(key string) (int64, error) {
	return global.Rdb.HLen(r.ctx, key).Result()
}

// IncrementHashField 对哈希字段进行数值增加
func (r *RedisHashManager) IncrementHashField(key, field string, increment int64) (int64, error) {
	return global.Rdb.HIncrBy(r.ctx, key, field, increment).Result()
}

// SMS验证码专用方法

// SetSmsCode 设置短信验证码（使用哈希存储）
func (r *RedisHashManager) SetSmsCode(mobile string, code int, expiration time.Duration) error {
	hashKey := "sms_codes:" + mobile
	fields := map[string]interface{}{
		"code":       code,
		"created_at": time.Now().Unix(),
		"attempts":   0, // 验证尝试次数
	}
	return r.SetHashFields(hashKey, fields, expiration)
}

// GetSmsCode 获取短信验证码
func (r *RedisHashManager) GetSmsCode(mobile string) (string, error) {
	hashKey := "sms_codes:" + mobile
	return r.GetHashField(hashKey, "code")
}

// IncrementSmsAttempts 增加短信验证码尝试次数
func (r *RedisHashManager) IncrementSmsAttempts(mobile string) (int64, error) {
	hashKey := "sms_codes:" + mobile
	return r.IncrementHashField(hashKey, "attempts", 1)
}

// GetSmsAttempts 获取短信验证码尝试次数
func (r *RedisHashManager) GetSmsAttempts(mobile string) (string, error) {
	hashKey := "sms_codes:" + mobile
	return r.GetHashField(hashKey, "attempts")
}

// DeleteSmsCode 删除短信验证码
func (r *RedisHashManager) DeleteSmsCode(mobile string) error {
	hashKey := "sms_codes:" + mobile
	return r.DeleteHash(hashKey)
}

// 订单管理专用方法

// SetOrderInfo 设置订单信息（使用哈希存储）
func (r *RedisHashManager) SetOrderInfo(orderId string, orderData map[string]interface{}, expiration time.Duration) error {
	hashKey := "order:" + orderId
	// 添加创建时间和默认状态
	orderData["created_at"] = time.Now().Unix()
	if _, exists := orderData["status"]; !exists {
		orderData["status"] = "waiting" // 默认状态：等待接单
	}
	return r.SetHashFields(hashKey, orderData, expiration)
}

// GetOrderInfo 获取订单信息
func (r *RedisHashManager) GetOrderInfo(orderId string) (map[string]string, error) {
	hashKey := "order:" + orderId
	return r.GetAllHashFields(hashKey)
}

// UpdateOrderStatus 更新订单状态
func (r *RedisHashManager) UpdateOrderStatus(orderId, status string) error {
	hashKey := "order:" + orderId
	return r.SetHashField(hashKey, "status", status, 0)
}

// UpdateOrderDriver 更新订单司机信息
func (r *RedisHashManager) UpdateOrderDriver(orderId, driverId string) error {
	hashKey := "order:" + orderId
	fields := map[string]interface{}{
		"driver_id":     driverId,
		"assigned_time": time.Now().Unix(),
		"status":        "assigned",
	}
	return r.SetHashFields(hashKey, fields, 0)
}

// DeleteOrderInfo 删除订单信息
func (r *RedisHashManager) DeleteOrderInfo(orderId string) error {
	hashKey := "order:" + orderId
	return r.DeleteHash(hashKey)
}

// GetOrderField 获取订单特定字段
func (r *RedisHashManager) GetOrderField(orderId, field string) (string, error) {
	hashKey := "order:" + orderId
	return r.GetHashField(hashKey, field)
}
