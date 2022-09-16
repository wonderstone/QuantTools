/* 股票账户模拟，仅保留最基本核心字段和方法 */
package stockaccount

import (
	"math" // for minimal comm is 5

	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/order"
)

// Declaring stockaccount struct with key fields
// no needs for settlement price or MTM P/L calculation
// 单一标的单一下单的持仓对应
type PositionDetail struct {
	UdTime    string
	InstID    string
	BasePrice float64 // 基准价 即买入价格 不修改
	LastPrice float64
	Num       float64
	Equity    float64
	SCP       *cp.SCP
}

// PositionDetail的产生和删除均基于order.StockOrder 但是删除操作需要交给PositionSlice
// 依据stockorder产生PositionDetail
// isExcuted 字段交给framework检查
func NewPositionDetail(so *order.StockOrder) PositionDetail {
	pd := PositionDetail{
		UdTime:    so.OrderTime,
		InstID:    so.InstID,
		BasePrice: so.OrderPrice,
		LastPrice: so.OrderPrice,
		Num:       so.OrderNum,
		SCP:       so.SCP,
	}
	pd.CalEquity(pd.LastPrice)
	return pd
}

// 计算股票市值equity 因行情数据引发
func (pd *PositionDetail) CalEquity(updatevalue float64) {
	pd.Equity = updatevalue * pd.Num * pd.SCP.ContractSize
}

// UMI更新最新价 因行情数据引发
func (pd *PositionDetail) UpdateLastPrice(time string, value float64) {
	pd.LastPrice = value
	pd.UdTime = time
	pd.CalEquity(pd.LastPrice)
}

// 计算手续费 因order引发 基于order或者position的num变化值
func CalComm(SO *order.StockOrder) float64 {
	TransFee := SO.OrderPrice * SO.OrderNum * SO.ContractSize * SO.SCP.TransferFeeRate
	Tax := 0.0
	switch SO.OrderDirection {
	case "Buy":
		Tax = 0
	case "Sell":
		Tax = SO.OrderPrice * SO.OrderNum * SO.ContractSize * SO.SCP.TaxRate
	default:
		panic("OrderDirection Error")
	}
	Comm := math.Max(SO.OrderPrice*SO.OrderNum*SO.ContractSize*SO.SCP.CommBrokerRate, 5.0)

	return TransFee + Tax + Comm
}

func (pd *PositionDetail) CalUnRealizedProfit(updatevalue float64) (Profit float64) {
	Profit = pd.Num * pd.SCP.ContractSize * (updatevalue - pd.BasePrice)
	return
}

func (pd *PositionDetail) CalRealizedProfit(num float64, updatevalue float64) (RealizedProfit float64) {
	RealizedProfit = num * pd.SCP.ContractSize * (updatevalue - pd.BasePrice)
	return
}
