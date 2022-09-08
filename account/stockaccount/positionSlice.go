/* 期货账户模拟，仅保留最基本核心字段和方法 */
package stockaccount

import (
	"github.com/wonderstone/QuantTools/order"
)

// 单一标的汇总持仓切片
type PositionSlice struct {
	UdTime   string
	PosTdys  []PositionDetail
	PosPrevs []PositionDetail
}

// 生成一个汇总持仓切片

func NewPosSlice() *PositionSlice {
	return &PositionSlice{}
}

// 计算PosTdyNum
func (ps *PositionSlice) CalPosTdyNum() (PosTdyNum float64) {
	for _, pd := range ps.PosTdys {
		PosTdyNum += pd.Num
	}
	return
}

// 计算PosPrevNum
func (ps *PositionSlice) CalPosPrevNum() (PosPrevNum float64) {
	for _, pd := range ps.PosPrevs {
		PosPrevNum += pd.Num
	}
	return
}

// 遍历获得所有Equity汇总
func (ps *PositionSlice) CalEquity() (Equity float64) {
	for _, pd := range ps.PosTdys {
		Equity += pd.Equity
	}
	for _, pd := range ps.PosPrevs {
		Equity += pd.Equity
	}
	return
}

// CalUnRealizedProfit

func (ps *PositionSlice) CalUnRealizedProfit() (UnRealizedProfit float64) {
	for _, pd := range ps.PosTdys {
		UnRealizedProfit += pd.CalUnRealizedProfit(pd.LastPrice)
	}
	for _, pd := range ps.PosPrevs {
		UnRealizedProfit += pd.CalUnRealizedProfit(pd.LastPrice)
	}
	return
}

func popupSlice(SO *order.StockOrder, s []PositionDetail) ([]PositionDetail, float64) {
	index := 0
	tmpsum := 0.0
	realizedprofit := 0.0
	for {
		tmpsum += s[index].Num
		if tmpsum > SO.OrderNum {
			realizedprofit += s[index].CalRealizedProfit(s[index].Num-tmpsum+SO.OrderNum, SO.OrderPrice)
			s[index].Num = tmpsum - SO.OrderNum
			break
		} else {
			realizedprofit += s[index].CalRealizedProfit(s[index].Num, SO.OrderPrice)
		}
		index++
		if index == len(s) {
			break
		}
	}
	s = s[index:]
	return s, realizedprofit
}

// order引发  添加 修改 删除 Pos进入PosTdys 计算RealizedProfit Comm
func (ps *PositionSlice) UpdateWithOrder(SO *order.StockOrder) (RealizedProfit float64, Comm float64, Equity float64) {
	// 修改
	switch SO.OrderDirection {
	case order.Buy:
		// 增加PosTdys
		ps.PosTdys = append(ps.PosTdys, NewPositionDetail(SO))
		Equity = ps.PosTdys[len(ps.PosTdys)-1].Equity
		// Equity = SO.OrderPrice * SO.OrderNum * SO.StockContractProp.ContractSize

	case order.Sell:
		if SO.OrderNum > ps.CalPosPrevNum() {
			panic("卖出超量  检查策略模块")
		}
		ps.PosPrevs, RealizedProfit = popupSlice(SO, ps.PosPrevs)
		// 这么算一遍更简单
		Equity = -SO.OrderPrice * SO.OrderNum * SO.SCP.ContractSize

	}
	// 计算comm
	Comm = CalComm(SO)
	// 更新一遍ps数据 可以避免Equity计算错误,内部包含1.updatetime 字段更新
	ps.UpdateWithUMI(SO.OrderTime, SO.OrderPrice)
	return RealizedProfit, Comm, Equity
}

// 遍历update一个UMI的同时 重新计算Equtiy
// 不输入updateMI的原因在于后续UpdateWithOrder
func (ps *PositionSlice) UpdateWithUMI(time string, value float64) {
	// 1. 时间字段更新
	ps.UdTime = time
	// 2. PosTdys内部pd全部更新
	for index := range ps.PosTdys {
		ps.PosTdys[index].UpdateLastPrice(time, value)
	}
	// 3. PosPrevs内部pd全部更新
	for index := range ps.PosPrevs {
		ps.PosPrevs[index].UpdateLastPrice(time, value)
	}
}

// CloseMarket引发  PosTdys转入PosPrevs，PosTdys清空
func (ps *PositionSlice) UpdateWithCM() {
	// 实际上已经由最后一组tick更新  而且一旦不以收盘价为最终数据  这个时间点意义也有限

	// PosTdys转入PosPrevs，PosTdys清空
	ps.PosPrevs = append(ps.PosPrevs, ps.PosTdys...)
	ps.PosTdys = nil

}
