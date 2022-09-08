/* 期货账户模拟，仅保留最基本核心字段和方法 */
package futuresaccount

import (
	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/order"
)

type DirType int

const (
	Long DirType = iota
	Short
)

// Declaring futuresaccount struct with key fields
// 单一标的单一下单的持仓对应
type PositionDetail struct {
	UdTime    string
	InstID    string
	Dir       DirType
	BasePrice float64 //MTM会被修改
	LastPrice float64
	Num       float64
	Margin    float64
	FCP       *cp.FCP
}

// PositionDetail的产生和删除均基于order.FuturesOrder 但是删除操作需要交给PositionSlice
// 依据futuresorder产生PositionDetail
func NewPositionDetail(fo *order.FuturesOrder) PositionDetail {
	pd := PositionDetail{
		UdTime:    fo.OrderTime,
		InstID:    fo.InstID,
		Dir:       DirType(fo.OrderDirection),
		BasePrice: fo.OrderPrice, //MTM会被修改
		LastPrice: fo.OrderPrice,
		Num:       fo.OrderNum,
		FCP:       fo.FCP,
	}
	pd.CalMargin(pd.LastPrice)
	return pd
}

// 计算保证金 因行情数据引发
func (pd *PositionDetail) CalMargin(updatevalue float64) {
	switch pd.Dir {
	case Long:
		pd.Margin = updatevalue * pd.Num * pd.FCP.ContractSize * (pd.FCP.MarginLong + pd.FCP.MarginBroker) / 100
	case Short:
		pd.Margin = updatevalue * pd.Num * pd.FCP.ContractSize * (pd.FCP.MarginShort + pd.FCP.MarginBroker) / 100
	}
}

// MTM更新基准价 因行情数据引发
func (pd *PositionDetail) UpdateBasePrice(time string, value float64) {
	pd.BasePrice = value
	pd.UdTime = time
}

// UMI更新最新价 因行情数据引发 基于类型为data.UpdateMI的MTMInfo指针参数
func (pd *PositionDetail) UpdateLastPrice(time string, value float64) {
	pd.LastPrice = value
	pd.UdTime = time
	pd.CalMargin(pd.LastPrice)
}

// 计算手续费 因order引发 基于order或者position的num变化值
func CalComm(FO *order.FuturesOrder) (comm float64) {
	if FO.IsCommRateType {
		switch FO.OrderType {
		case order.Open:
			comm = FO.OrderPrice * FO.OrderNum * FO.ContractSize * (FO.CommOpen + FO.CommBroker)
		case order.CloseToday:
			comm = FO.OrderPrice * FO.OrderNum * FO.ContractSize * (FO.CommCloseToday + FO.CommBroker)
		case order.ClosePrevious:
			comm = FO.OrderPrice * FO.OrderNum * FO.ContractSize * (FO.CommClosePrevious + FO.CommBroker)
		}
	} else {
		switch FO.OrderType {
		case order.Open:
			comm = FO.OrderNum * (FO.CommOpen + FO.CommBroker)
		case order.CloseToday:
			comm = FO.OrderNum * (FO.CommCloseToday + FO.CommBroker)
		case order.ClosePrevious:
			comm = FO.OrderNum * (FO.CommClosePrevious + FO.CommBroker)
		}
	}
	return
}

// 计算盈亏实际上是一个核心公式  但前后伴随着不同的操作
// 计算盯市盈亏 由MTM机制的行情数据引发 settlementprice作为updatevalue，并修改baseprice，profit计入account.cash
// 计算持仓盈亏 由市值刷新引发的浮动盈亏 FuturesContract.ContractSize 与 UpdateMI.Value充当updatevalue
func (pd *PositionDetail) CalUnRealizedProfit(updatevalue float64) (Profit float64) {
	switch pd.Dir {
	case Long:
		Profit = pd.Num * pd.FCP.ContractSize * (updatevalue - pd.BasePrice)
	case Short:
		Profit = pd.Num * pd.FCP.ContractSize * (pd.BasePrice - updatevalue)
	}
	return
}

// 计算平仓盈亏 加入num是因为单子可能会拆 对应不同的持仓 由order引发 RealizedProfit计入MktVal   updatevalue使用order.price
// 持仓数改变问题应该在上一层slice级解决，因为下降到0时会被slice剔除，这个操作在pd级别做不了
func (pd *PositionDetail) CalRealizedProfit(num float64, updatevalue float64) (RealizedProfit float64) {
	switch pd.Dir {
	case Long:
		RealizedProfit = num * pd.FCP.ContractSize * (updatevalue - pd.BasePrice)
	case Short:
		RealizedProfit = num * pd.FCP.ContractSize * (pd.BasePrice - updatevalue)
	}
	return
}
