package order

import (
	cp "github.com/wonderstone/QuantTools/contractproperty"
)

type StockOrder struct {
	InstID         string
	IsEligible     bool
	IsExecuted     bool
	OrderTime      string
	OrderPrice     float64
	OrderNum       float64
	OrderDirection string
	*cp.SCP        // Promoted fields
}

// 此处用panic机制有待权衡 增加运算 过于粗暴  V.S 便于发现问题
func NewStockOrder(instID string, iseligible bool, isexecuted bool, ordertime string, orderprice float64, ordernum float64, orderdir string, pscp *cp.SCP) StockOrder {
	if orderprice <= 0 {
		panic("下单价格小于0 检查策略模块或相关数据")
	}
	if ordernum <= 0 {
		panic("下单数量小于0 检查策略模块")
	}
	return StockOrder{
		InstID:         instID,
		IsEligible:     iseligible,
		IsExecuted:     isexecuted,
		OrderTime:      ordertime,
		OrderPrice:     orderprice,
		OrderNum:       ordernum,
		OrderDirection: orderdir,
		SCP:            pscp,
	}
}

// cal equity
func (SO *StockOrder) CalEquity() (Equity float64) {
	Equity = SO.OrderPrice * SO.OrderNum * SO.ContractSize * 1.01
	return Equity
}
