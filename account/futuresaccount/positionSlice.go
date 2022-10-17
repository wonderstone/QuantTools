/* 期货账户模拟，仅保留最基本核心字段和方法 */
package futuresaccount

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

// 计算PosTdyNumL 和PosTdyNumS
func (ps *PositionSlice) CalPosTdyNum() (PosTdyNumL float64, PosTdyNumS float64) {
	if ps == nil {
		return
	}
	for _, pd := range ps.PosTdys {
		if pd.Dir == "Buy" {
			PosTdyNumL += pd.Num
		} else {
			PosTdyNumS += pd.Num
		}
	}
	return
}

// 计算PosPrevNumL 和osPrevNumS
func (ps *PositionSlice) CalPosPrevNum() (PosPrevNumL float64, PosPrevNumS float64) {
	if ps == nil {
		return
	}
	for _, pd := range ps.PosPrevs {
		if pd.Dir == "Buy" {
			PosPrevNumL += pd.Num
		} else {
			PosPrevNumS += pd.Num
		}
	}
	return
}

// 遍历获得所有Margin汇总
func (ps *PositionSlice) CalMargin() (Margin float64) {
	if ps == nil {
		return
	}
	for _, pd := range ps.PosTdys {
		Margin += pd.Margin
	}
	for _, pd := range ps.PosPrevs {
		Margin += pd.Margin
	}
	return
}

// CalUnRealizedProfit

func (ps *PositionSlice) CalUnRealizedProfit() (UnRealizedProfit float64) {
	if ps == nil {
		return
	}
	for _, pd := range ps.PosTdys {
		UnRealizedProfit += pd.CalUnRealizedProfit(pd.LastPrice)
	}
	for _, pd := range ps.PosPrevs {
		UnRealizedProfit += pd.CalUnRealizedProfit(pd.LastPrice)
	}
	return
}

func popupSlice(FO *order.FuturesOrder, s []PositionDetail) ([]PositionDetail, float64) {
	index := 0
	tmpsum := 0.0
	realizedprofit := 0.0
	for {
		tmpsum += s[index].Num
		if tmpsum > FO.OrderNum {
			realizedprofit += s[index].CalRealizedProfit(s[index].Num-tmpsum+FO.OrderNum, FO.OrderPrice)
			s[index].Num = tmpsum - FO.OrderNum
			break
		} else {
			realizedprofit += s[index].CalRealizedProfit(s[index].Num, FO.OrderPrice)
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
func (ps *PositionSlice) UpdateWithOrder(FO *order.FuturesOrder) (RealizedProfit float64, Comm float64, UnRealizedProfit float64) {
	// 修改
	switch FO.OrderType {
	case "Open":
		// 增加PosTdys
		ps.PosTdys = append(ps.PosTdys, NewPositionDetail(FO))
		// MarginChange = FO.CalMargin()
		// 没有
	case "CloseToday":
		ps.PosTdys, RealizedProfit = popupSlice(FO, ps.PosTdys)
		// MarginChange = -FO.CalMargin()
	case "ClosePrevious":
		ps.PosPrevs, RealizedProfit = popupSlice(FO, ps.PosPrevs)
		// MarginChange = -FO.CalMargin()
	default:
		panic("OrderType Error")
	}
	// 先更新一遍ps数据 可以避免Margin计算错误,内部包含1.updatetime 字段更新
	ps.UpdateWithUMI(FO.OrderTime, FO.OrderPrice)
	UnRealizedProfit = ps.CalUnRealizedProfit()
	// 计算comm
	Comm = CalComm(FO)
	return RealizedProfit, Comm, UnRealizedProfit
}

// 遍历update一个UMI的同时 重新计算Margin
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

// MTM引发  PosTdys转入PosPrevs，PosTdys清空
func (ps *PositionSlice) UpdateWithMTM(time string, value float64) (RealizedProfit float64) {
	// 1. 时间字段更新
	ps.UdTime = time

	// PosTdys转入PosPrevs，PosTdys清空
	ps.PosPrevs = append(ps.PosPrevs, ps.PosTdys...)
	ps.PosTdys = nil

	// 更新时间 计算MTM的profit, 更改baseprice
	for index := range ps.PosPrevs {
		// MTM时，所有pd基于MTMprice将unrealized转化为realized
		RealizedProfit += ps.PosPrevs[index].CalUnRealizedProfit(value)
		ps.PosPrevs[index].UpdateBasePrice(time, value)
		ps.PosPrevs[index].UpdateLastPrice(time, value)
	}
	return
}
