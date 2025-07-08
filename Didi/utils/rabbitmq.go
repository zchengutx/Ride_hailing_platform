package utils

import (
	"Didi/rpc/common/global"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// RabbitMQManager RabbitMQ管理器
type RabbitMQManager struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

var rabbitMQManager *RabbitMQManager

// InitRabbitMQ 初始化RabbitMQ连接
func InitRabbitMQ() error {
	rmqConf := &global.AppConf.RabbitMQ

	// 构建连接URL
	url := fmt.Sprintf("amqp://%s:%s@%s:%d%s",
		rmqConf.User,
		rmqConf.Password,
		rmqConf.Host,
		rmqConf.Port,
		rmqConf.Vhost)

	// 创建连接
	conn, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("RabbitMQ连接失败: %v", err)
	}

	// 创建通道
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("RabbitMQ通道创建失败: %v", err)
	}

	rabbitMQManager = &RabbitMQManager{
		conn:    conn,
		channel: channel,
	}

	// 声明所需的交换机和队列
	if err := rabbitMQManager.setupQueues(); err != nil {
		return fmt.Errorf("RabbitMQ队列设置失败: %v", err)
	}

	log.Println("RabbitMQ连接成功")
	return nil
}

// GetRabbitMQManager 获取RabbitMQ管理器实例
func GetRabbitMQManager() *RabbitMQManager {
	return rabbitMQManager
}

// setupQueues 设置交换机和队列
func (r *RabbitMQManager) setupQueues() error {
	// 声明订单推送交换机（扇出模式，广播给所有司机）
	err := r.channel.ExchangeDeclare(
		"order.fanout", // 交换机名称
		"fanout",       // 交换机类型
		true,           // 持久化
		false,          // 自动删除
		false,          // 内部使用
		false,          // 等待服务器响应
		nil,            // 额外参数
	)
	if err != nil {
		return fmt.Errorf("声明交换机失败: %v", err)
	}

	// 声明司机通知队列
	_, err = r.channel.QueueDeclare(
		"driver.notifications", // 队列名称
		true,                   // 持久化
		false,                  // 自动删除
		false,                  // 独占
		false,                  // 等待服务器响应
		nil,                    // 额外参数
	)
	if err != nil {
		return fmt.Errorf("声明队列失败: %v", err)
	}

	return nil
}

// OrderMessage 订单消息结构
type OrderMessage struct {
	OrderID       string   `json:"order_id"`
	UserID        string   `json:"user_id"`
	StartLocation string   `json:"start_location"`
	EndLocation   string   `json:"end_location"`
	CartType      string   `json:"cart_type"`
	CreatedAt     int64    `json:"created_at"`
	Latitude      float64  `json:"latitude"`   // 起始地纬度
	Longitude     float64  `json:"longitude"`  // 起始地经度
	DriverIDs     []string `json:"driver_ids"` // 周边司机ID列表
}

// PublishOrderToNearbyDrivers 发布订单消息给周边司机
func (r *RabbitMQManager) PublishOrderToNearbyDrivers(orderMsg OrderMessage) error {
	// 将消息序列化为JSON
	msgBody, err := json.Marshal(orderMsg)
	if err != nil {
		return fmt.Errorf("消息序列化失败: %v", err)
	}

	// 发布消息到交换机
	err = r.channel.Publish(
		"order.fanout", // 交换机名称
		"",             // 路由键（fanout类型忽略）
		false,          // 立即发送
		false,          // 不等待确认
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         msgBody,
			DeliveryMode: amqp.Persistent, // 消息持久化
			Headers: amqp.Table{
				"message_type": "new_order",
			},
		},
	)

	if err != nil {
		return fmt.Errorf("消息发布失败: %v", err)
	}

	log.Printf("订单消息发布成功，订单ID: %s, 目标司机: %v", orderMsg.OrderID, orderMsg.DriverIDs)
	return nil
}

// PublishDirectMessage 向特定司机发送直接消息
func (r *RabbitMQManager) PublishDirectMessage(driverID string, message interface{}) error {
	msgBody, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("消息序列化失败: %v", err)
	}

	// 声明司机专用队列
	queueName := fmt.Sprintf("driver.%s", driverID)
	_, err = r.channel.QueueDeclare(
		queueName,
		true,  // 持久化
		false, // 不自动删除
		false, // 不独占
		false, // 不等待
		nil,
	)
	if err != nil {
		return fmt.Errorf("声明司机队列失败: %v", err)
	}

	// 发送消息
	err = r.channel.Publish(
		"",        // 默认交换机
		queueName, // 路由键为队列名
		false,     // 立即发送
		false,     // 不等待确认
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         msgBody,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		return fmt.Errorf("消息发布失败: %v", err)
	}

	log.Printf("向司机 %s 发送消息成功", driverID)
	return nil
}

// ConsumeDriverMessages 消费司机消息（供司机端使用）
func (r *RabbitMQManager) ConsumeDriverMessages(driverID string, handler func([]byte) error) error {
	queueName := fmt.Sprintf("driver.%s", driverID)

	// 声明队列
	_, err := r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("声明队列失败: %v", err)
	}

	// 绑定到通知交换机（接收广播消息）
	err = r.channel.QueueBind(
		queueName,      // 队列名
		"",             // 路由键
		"order.fanout", // 交换机名
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("队列绑定失败: %v", err)
	}

	// 开始消费消息
	msgs, err := r.channel.Consume(
		queueName,
		"",    // 消费者标签
		true,  // 自动确认
		false, // 不独占
		false, // 不等待
		false, // 不阻塞
		nil,
	)
	if err != nil {
		return fmt.Errorf("开始消费失败: %v", err)
	}

	// 处理消息
	go func() {
		for msg := range msgs {
			if err := handler(msg.Body); err != nil {
				log.Printf("处理消息失败: %v", err)
			}
		}
	}()

	log.Printf("司机 %s 开始监听消息", driverID)
	return nil
}

// Close 关闭连接
func (r *RabbitMQManager) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
