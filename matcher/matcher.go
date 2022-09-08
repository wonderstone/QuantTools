package matcher

import (
	"github.com/wonderstone/QuantTools/order"
)

// 简易撮合机制 用以抽象模拟报单 撮合成交过程
// 内部包含滑点调整、订单时间调整(假定全部成交 顾只保留一个时间字段以备检索 尽管更接近订单原本时间)、订单执行状态调整
// 撮合简易使用下一时间戳的价格，例如下一个bar的开盘价
// 均外部可见，请最大限度避免在其他代码中隐式修改，应该保证有锁的情况下调用
type SimpleMatcher struct {
	Slippage4S float64
	Slippage4F float64
}

// Slippage简易使用合约最小变动单位的整数倍 股票类型建议为0.01元的整数倍   期货直接取TickSize字段*Slippage4F调整数量
func NewSimpleMatcher(slpg4s float64, slpg4f float64) SimpleMatcher {
	return SimpleMatcher{
		Slippage4S: slpg4s,
		Slippage4F: slpg4f,
	}
}

// 模拟撮合进行不利方向的价格调整，matchprice的选择建议使用下一个bar的开盘价matchtime使用下一个bar的时间戳
func (m *SimpleMatcher) MatchFuturesOrder(FO *order.FuturesOrder, matchprice float64, matchtime string) {
	// in principle, backtest should be done under one mutex lock
	// insurance: add a mutex for matchfuturesorder
	// num4ticksize for how many ticksizes to adjust

	switch FO.OrderDirection {
	case order.Buy:
		FO.OrderPrice = matchprice + m.Slippage4F*FO.TickSize
	case order.Sell:
		FO.OrderPrice = matchprice - m.Slippage4F*FO.TickSize
	}
	FO.OrderTime = matchtime
	FO.IsExecuted = true
}

// 模拟撮合进行不利方向的价格调整，matchprice的选择建议使用下一个bar的开盘价matchtime使用下一个bar的时间戳
func (m *SimpleMatcher) MatchStockOrder(SO *order.StockOrder, matchprice float64, matchtime string) {
	// in principle, backtest should be done under one mutex lock
	// insurance: add a mutex for matchstockorder
	// m.Lock()
	switch SO.OrderDirection {
	case order.Buy:
		SO.OrderPrice = matchprice + m.Slippage4S
	case order.Sell:
		SO.OrderPrice = matchprice - m.Slippage4S
	}
	SO.OrderTime = matchtime
	SO.IsExecuted = true
	// m.Unlock()
}
