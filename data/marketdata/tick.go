/* QMT tick: {'600000.SH': {'timetag': '20220726 14:30:20', 'lastPrice': 7.38, 'open': 7.36, 'high': 7.390000000000001, 'low': 7.34, 'lastClose': 7.35, 'amount': 115969100.0, 'volume': 157398, 'pvolume': 0, 'stockStatus': 0, 'openInt': 13, 'settlementPrice': 0.0, 'lastSettlementPrice': 0.0, 'askPrice': [0.0, 0.0, 0.0, 0.0, 0.0], 'bidPrice': [0.0, 0.0, 0.0, 0.0, 0.0], 'askVol': [0, 0, 0, 0, 0], 'bidVol': [0, 0, 0, 0, 0]}}*/
// 期货数据暂时没有得到，但QMT宣称一致
// 由于其没有盘口数据，考虑到应用和运行效率折中，保留一组Bid Ask

package marketdata

import (
	"strconv"

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
	OpenPrice       float64   //建议考虑删除 日开盘价
	HighestPrice    float64   //建议考虑删除 日最高价
	LowestPrice     float64   //建议考虑删除 日最低价
	Volume          float64   //成交量
	Amount          float64   //成交额 ctp里面用的是turnover
	OpenInterest    float64   //持仓数
	BidPrice        []float64 //建议删除
	BidVolume       []float64 //建议删除
	AskPrice        []float64 //建议删除
	AskVolume       []float64 //建议删除
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
	Amount          float64   //成交额 没错  它竟然使用了turnover而不是amount
	BidPrice        []float64 //建议删除
	BidVolume       []float64 //建议删除
	AskPrice        []float64 //建议删除
	AskVolume       []float64 //建议删除
}

// NewFuturesTick returns a new FuturesTick
func NewFuturesTick(UpdateTime, InstID string, LastPrice, OpenPrice, HighestPrice, LowestPrice, Volume, Amount, OpenInterest float64, BidPrice, BidVolume, AskPrice, AskVolume []float64) (ft FuturesTick) {
	ft.UpdateTimeStamp = UpdateTime
	ft.InstID = InstID
	ft.LastPrice = LastPrice
	ft.OpenPrice = OpenPrice
	ft.HighestPrice = HighestPrice
	ft.LowestPrice = LowestPrice
	ft.Volume = Volume
	ft.Amount = Amount
	ft.OpenInterest = OpenInterest
	ft.BidPrice = BidPrice
	ft.BidVolume = BidVolume
	ft.AskPrice = AskPrice
	ft.AskVolume = AskVolume
	return
}

// NewStockTick returns a new StockTick
func NewStockTick(UpdateTime, InstID string, LastPrice, OpenPrice, HighestPrice, LowestPrice, Volume, Amount float64, BidPrice, BidVolume, AskPrice, AskVolume []float64) (st StockTick) {
	st.UpdateTimeStamp = UpdateTime
	st.InstID = InstID
	st.LastPrice = LastPrice
	st.OpenPrice = OpenPrice
	st.HighestPrice = HighestPrice
	st.LowestPrice = LowestPrice
	st.Volume = Volume
	st.Amount = Amount
	st.BidPrice = BidPrice
	st.BidVolume = BidVolume
	st.AskPrice = AskPrice
	st.AskVolume = AskVolume
	return
}

// Get updateInfo from  tick data by tag
func (ft *FuturesTick) GetUpdateInfo(tag string) (ui data.UpdateMI) {
	// get the last character of tag
	// if it is a number, get the number
	v, _ := strconv.Atoi(tag[len(tag)-1:])
	ui.UpdateTimeStamp = ft.UpdateTimeStamp
	ui.InstID = ft.InstID
	switch tag {
	case "LastPrice":
		ui.Value = ft.LastPrice
	case "OpenPrice":
		ui.Value = ft.OpenPrice
	case "HighestPrice":
		ui.Value = ft.HighestPrice
	case "LowestPrice":
		ui.Value = ft.LowestPrice
	case "Volume":
		ui.Value = ft.Volume
	case "Amount":
		ui.Value = ft.Amount
	case "OpenInterest":
		ui.Value = ft.OpenInterest
	case "BidPrice" + strconv.Itoa(v):
		ui.Value = ft.BidPrice[v]
	case "BidVolume" + strconv.Itoa(v):
		ui.Value = ft.BidVolume[v]
	case "AskPrice" + strconv.Itoa(v):
		ui.Value = ft.AskPrice[v]
	case "AskVolume" + strconv.Itoa(v):
		ui.Value = ft.AskVolume[v]
	default:
		panic("tag not found")
	}
	return ui
}

// Get updateInfor from tick data by tag
func (st *StockTick) GetUpdateInfo(tag string) (ui data.UpdateMI) {
	v, _ := strconv.Atoi(tag[len(tag)-1:])
	ui.UpdateTimeStamp = st.UpdateTimeStamp
	ui.InstID = st.InstID
	switch tag {
	case "LastPrice":
		ui.Value = st.LastPrice
	case "OpenPrice":
		ui.Value = st.OpenPrice
	case "HighestPrice":
		ui.Value = st.HighestPrice
	case "LowestPrice":
		ui.Value = st.LowestPrice
	case "Volume":
		ui.Value = st.Volume
	case "Amount":
		ui.Value = st.Amount
	case "BidPrice" + strconv.Itoa(v):
		ui.Value = st.BidPrice[v]
	case "BidVolume" + strconv.Itoa(v):
		ui.Value = st.BidVolume[v]
	case "AskPrice" + strconv.Itoa(v):
		ui.Value = st.AskPrice[v]
	case "AskVolume" + strconv.Itoa(v):
		ui.Value = st.AskVolume[v]
	default:
		ui.Value = st.LastPrice
	}
	return ui
}
