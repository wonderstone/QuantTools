package order

import (
	cp "github.com/wonderstone/QuantTools/contractproperty"
)

// 存在过度设计嫌疑，回退到string判断。
// type FuturesOrderTYP int

// const (
// 	Open FuturesOrderTYP = iota
// 	CloseToday
// 	ClosePrevious
// )
// 存在过度设计嫌疑，回退到string判断。
// type OrderDir int // 与stock类型共用

// const (
// 	Buy OrderDir = iota // 注意  OrderDir 参考Python-CTP规范，直接靠近CTP端
// 	Sell
// )

type FuturesOrder struct {
	InstID         string
	IsExecuted     bool
	OrderTime      string
	OrderPrice     float64
	OrderNum       float64
	OrderDirection string
	OrderType      string
	*cp.FCP        // Promoted fields
}

// func NewFuturesOrder(instID string, isexecuted bool, ordertime string, orderprice float64, ordernum float64, orderdir OrderDir, ordertype FuturesOrderTYP, pfcp *cp.FCP) FuturesOrder {
func NewFuturesOrder(instID string, isexecuted bool, ordertime string, orderprice float64, ordernum float64, orderdir string, ordertype string, pfcp *cp.FCP) FuturesOrder {
	if orderprice <= 0 {
		panic("下单价格小于0 检查策略模块或相关数据")
	}
	if ordernum <= 0 {
		panic("下单数量小于0 检查策略模块")
	}
	return FuturesOrder{
		InstID:         instID,
		IsExecuted:     isexecuted,
		OrderTime:      ordertime,
		OrderPrice:     orderprice,
		OrderNum:       ordernum,
		OrderDirection: orderdir,
		OrderType:      ordertype,
		FCP:            pfcp,
	}
}

// futuresaccount need this method to check fundavail
func (FO *FuturesOrder) CalMargin() (Margin float64) {
	switch FO.OrderDirection {
	case "Buy":
		Margin = FO.OrderPrice * FO.OrderNum * FO.ContractSize * (FO.MarginLong + FO.MarginBroker) / 100
	case "Sell":
		Margin = FO.OrderPrice * FO.OrderNum * FO.ContractSize * (FO.MarginShort + FO.MarginBroker) / 100
	default:
		panic("OrderDirection Error")
	}
	return Margin
}
