package order

import (
	cp "github.com/wonderstone/QuantTools/contractproperty"
)

type FuturesOrder struct {
	InstID         string
	IsEligible     bool
	IsExecuted     bool
	OrderTime      string
	OrderPrice     float64
	OrderNum       float64
	OrderDirection string
	OrderType      string
	*cp.FCP        // Promoted fields
}

// func NewFuturesOrder(instID string, isexecuted bool, ordertime string, orderprice float64, ordernum float64, orderdir OrderDir, ordertype FuturesOrderTYP, pfcp *cp.FCP) FuturesOrder {
func NewFuturesOrder(instID string, iseligible bool, isexecuted bool, ordertime string, orderprice float64, ordernum float64, orderdir string, ordertype string, pfcp *cp.FCP) FuturesOrder {
	if orderprice <= 0 {
		panic("下单价格小于0 检查策略模块或相关数据")
	}
	if ordernum <= 0 {
		panic("下单数量小于0 检查策略模块")
	}
	return FuturesOrder{
		InstID:         instID,
		IsEligible:     iseligible,
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
