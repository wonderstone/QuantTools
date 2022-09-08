/* QMT tick: {'600000.SH': {'timetag': '20220726 14:30:20', 'lastPrice': 7.38, 'open': 7.36, 'high': 7.390000000000001, 'low': 7.34, 'lastClose': 7.35, 'amount': 115969100.0, 'volume': 157398, 'pvolume': 0, 'stockStatus': 0, 'openInt': 13, 'settlementPrice': 0.0, 'lastSettlementPrice': 0.0, 'askPrice': [0.0, 0.0, 0.0, 0.0, 0.0], 'bidPrice': [0.0, 0.0, 0.0, 0.0, 0.0], 'askVol': [0, 0, 0, 0, 0], 'bidVol': [0, 0, 0, 0, 0]}}*/
// 期货数据暂时没有得到，但QMT宣称一致
// 由于其没有盘口数据，考虑到应用和运行效率折中，保留一组Bid Ask

package marketdata

import (
	"github.com/wonderstone/QuantTools/data"
)

// Declaring futuresTick struct with key fields
// CTP的TradingDay与Actionday字段需要全部保留，用以区分三个不同市场对tick数据区分的不一致
// 回测简化时需要考虑统一标准，按照时间排序
type FuturesTick struct {
	// 共13项
	UpdateTimeStamp string
	InstID          string
	LastPrice       float64
	OpenPrice       float64 //建议考虑删除 日开盘价
	HighestPrice    float64 //建议考虑删除 日最高价
	LowestPrice     float64 //建议考虑删除 日最低价
	Volume          float64 //成交量
	Amount          float64 //成交额 ctp里面用的是turnover
	OpenInterest    float64 //持仓数
	BidPrice1       float64 //建议删除
	BidVolume1      float64 //建议删除
	AskPrice1       float64 //建议删除
	AskVolume1      float64 //建议删除
}

/*********************/
/* 未完结需要进一步完善 */
/*********************/
// Declaring stockTick struct with key fields
// needs more attention with level 2 data
type StockTick struct {
	// 共12项
	UpdateTimeStamp string
	InstID          string
	LastPrice       float64
	OpenPrice       float64
	HighestPrice    float64
	LowestPrice     float64
	Volume          float64
	Amount          float64 //成交额 没错  它竟然使用了turnover而不是amount
	BidPrice1       float64 //建议删除
	BidVolume1      float64 //建议删除
	AskPrice1       float64 //建议删除
	AskVolume1      float64 //建议删除
}

// NewFuturesTick returns a new FuturesTick
func NewFuturesTick(UpdateTime, InstID string, LastPrice, OpenPrice, HighestPrice, LowestPrice, Volume, Amount, OpenInterest, BidPrice1, BidVolume1, AskPrice1, AskVolume1 float64) (ft FuturesTick) {
	ft.UpdateTimeStamp = UpdateTime
	ft.InstID = InstID
	ft.LastPrice = LastPrice
	ft.OpenPrice = OpenPrice
	ft.HighestPrice = HighestPrice
	ft.LowestPrice = LowestPrice
	ft.Volume = Volume
	ft.Amount = Amount
	ft.OpenInterest = OpenInterest
	ft.BidPrice1 = BidPrice1
	ft.BidVolume1 = BidVolume1
	ft.AskPrice1 = AskPrice1
	ft.AskVolume1 = AskVolume1
	return
}

// NewStockTick returns a new StockTick
func NewStockTick(UpdateTime, InstID string, LastPrice, OpenPrice, HighestPrice, LowestPrice, Volume, Amount, BidPrice1, BidVolume1, AskPrice1, AskVolume1 float64) (st StockTick) {
	st.UpdateTimeStamp = UpdateTime
	st.InstID = InstID
	st.LastPrice = LastPrice
	st.OpenPrice = OpenPrice
	st.HighestPrice = HighestPrice
	st.LowestPrice = LowestPrice
	st.Volume = Volume
	st.Amount = Amount
	st.BidPrice1 = BidPrice1
	st.BidVolume1 = BidVolume1
	st.AskPrice1 = AskPrice1
	st.AskVolume1 = AskVolume1
	return
}

func (ft *FuturesTick) GetUpdateInfo() (ui data.UpdateMI) {
	ui.UpdateTimeStamp = ft.UpdateTimeStamp
	ui.InstID = ft.InstID
	ui.Value = ft.LastPrice

	return ui
}
func (st *StockTick) GetUpdateInfo() (ui data.UpdateMI) {
	ui.UpdateTimeStamp = st.UpdateTimeStamp
	ui.InstID = st.InstID
	ui.Value = st.LastPrice

	return ui
}
