package model

import "time"

// 钱包表
type Wallet struct {
	Id           string    `gorm:"column:id;type:varchar(191);primaryKey;not null;" json:"id"`
	UserID       string    `gorm:"column:user_id;type:varchar(191);not null;" json:"user_id"`
	Balance      float64   `gorm:"column:balance;type:decimal(10,2);default:0.00;" json:"balance"`
	FrozenAmount float64   `gorm:"column:frozen_amount;type:decimal(10,2);default:0.00;" json:"frozen_amount"`
	TotalIncome  float64   `gorm:"column:total_income;type:decimal(10,2);default:0.00;" json:"total_income"`
	TotalExpense float64   `gorm:"column:total_expense;type:decimal(10,2);default:0.00;" json:"total_expense"`
	CreateTime   time.Time `gorm:"column:create_time;type:datetime(3);default:CURRENT_TIMESTAMP(3);" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;type:datetime(3);default:CURRENT_TIMESTAMP(3);" json:"update_time"`

	// 关联
	LxhPassenger *LxhPassenger `gorm:"foreignKey:UserID;references:Id" json:"user,omitempty"`
}

// 钱包流水表
type WalletTransaction struct {
	Id              string    `gorm:"column:id;type:varchar(191);primaryKey;not null;" json:"id"`
	WalletID        string    `gorm:"column:wallet_id;type:varchar(191);not null;" json:"wallet_id"`
	UserID          string    `gorm:"column:user_id;type:varchar(191);not null;" json:"user_id"`
	TransactionType int8      `gorm:"column:transaction_type;type:tinyint;not null;" json:"transaction_type"`
	Amount          float64   `gorm:"column:amount;type:decimal(10,2);not null;" json:"amount"`
	BalanceBefore   float64   `gorm:"column:balance_before;type:decimal(10,2);not null;" json:"balance_before"`
	BalanceAfter    float64   `gorm:"column:balance_after;type:decimal(10,2);not null;" json:"balance_after"`
	OrderID         *string   `gorm:"column:order_id;type:varchar(191);default:NULL;" json:"order_id"`
	Description     *string   `gorm:"column:description;type:varchar(200);default:NULL;" json:"description"`
	CreateTime      time.Time `gorm:"column:create_time;type:datetime(3);default:CURRENT_TIMESTAMP(3);" json:"create_time"`

	// 关联
	Wallet       *Wallet       `gorm:"foreignKey:WalletID;references:Id" json:"wallet,omitempty"`
	LxhPassenger *LxhPassenger `gorm:"foreignKey:UserID;references:Id" json:"user,omitempty"`
	LxhOrder     *LxhOrder     `gorm:"foreignKey:OrderID;references:Id" json:"order,omitempty"`
}

// 交易类型常量
const (
	TransactionTypeRecharge = 1 // 充值
	TransactionTypeExpense  = 2 // 消费
	TransactionTypeRefund   = 3 // 退款
	TransactionTypeWithdraw = 4 // 提现
	TransactionTypeReward   = 5 // 奖励
)

// 交易类型文本映射
var TransactionTypeText = map[int8]string{
	TransactionTypeRecharge: "充值",
	TransactionTypeExpense:  "消费",
	TransactionTypeRefund:   "退款",
	TransactionTypeWithdraw: "提现",
	TransactionTypeReward:   "奖励",
}

// 获取交易类型文本
func (wt *WalletTransaction) GetTypeText() string {
	if text, exists := TransactionTypeText[wt.TransactionType]; exists {
		return text
	}
	return "未知类型"
}

// 检查是否为收入类型
func (wt *WalletTransaction) IsIncome() bool {
	return wt.TransactionType == TransactionTypeRecharge ||
		wt.TransactionType == TransactionTypeRefund ||
		wt.TransactionType == TransactionTypeReward
}

// 检查是否为支出类型
func (wt *WalletTransaction) IsExpense() bool {
	return wt.TransactionType == TransactionTypeExpense ||
		wt.TransactionType == TransactionTypeWithdraw
}

func (Wallet) TableName() string {
	return "wallet"
}

func (WalletTransaction) TableName() string {
	return "wallet_transaction"
}
