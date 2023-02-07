package strategyModule

// market timing strategy
// one expression for one strategy, therefore one market for one strategy
import (
	"math"
	// "strings"

	"github.com/rs/zerolog/log"
	"github.com/wonderstone/QuantTools/account/virtualaccount"
	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/dataprocessor"
	"github.com/wonderstone/QuantTools/order"

	"github.com/wonderstone/QuantTools/configer"
)

type DMTStrategy struct {
	SInstNames     []string           // 股票标的名称
	SIndiNames     []string           // 股票参与GEP指标名称，注意其数量不大于BarDE内信息数量，且strategy内可见BarDE的数据
	SNum           float64            // 股票标的数量
	FInstNames     []string           // 期货标的名称
	FIndiNames     []string           // 期货参与GEP指标名称， there should be a rollover field in the futures IndiDataNames slice
	FNum           float64            // 期货标的数量
	STimeCritic    string             // 时间关键字，用于判断是否需要进行交易
	FTimeCritic    string             // 时间关键字，用于判断是否需要进行交易
	lastBuyValue   map[string]float64 // 公式上一次买入信号数值
	lastSellValue  map[string]float64 // 公式上一次卖出信号数值
	stimeCondition bool               // 时间条件是否满足
	ifsdone        bool               // 是否已经完成了一次交易
	ftimeCondition bool               // 时间条件是否满足
	iffdone        bool               // 是否已经完成了一次交易
}

func NewDMTStrategy(SInstNms, SIndiNms, FInstNms, FIndiNms []string, Snum, Fnum float64, STimeCritic, FTimeCritic string) DMTStrategy {
	return DMTStrategy{
		SInstNames:     SInstNms,
		SIndiNames:     SIndiNms,
		SNum:           Snum,
		FInstNames:     FInstNms,
		FIndiNames:     FIndiNms,
		FNum:           Fnum,
		STimeCritic:    STimeCritic,
		FTimeCritic:    FTimeCritic,
		lastBuyValue:   make(map[string]float64),
		lastSellValue:  make(map[string]float64),
		stimeCondition: false,
		ifsdone:        false,
		ftimeCondition: false,
		iffdone:        false,
	}
}

// this function is nessary for the framework
func NewDMTStrategyFromConfig(dir string, BTConfile string, sec string, StgConfile string) DMTStrategy {
	// c := configer.New(dir + "BackTest.yaml")
	c := configer.New(dir + BTConfile)
	err := c.Load()
	if err != nil {
		panic(err)
	}
	err = c.Unmarshal()
	if err != nil {
		panic(err)
	}
	tmpMap := c.GetStringMap(sec)

	var sinstrnames []string
	for _, v := range tmpMap["sinstrnames"].([]interface{}) {
		sinstrnames = append(sinstrnames, v.(string))
	}
	var sindinames []string
	for _, v := range tmpMap["sindinames"].([]interface{}) {
		sindinames = append(sindinames, v.(string))
	}

	var finstrnames []string
	for _, v := range tmpMap["finstrnames"].([]interface{}) {
		finstrnames = append(finstrnames, v.(string))
	}
	var findinames []string
	for _, v := range tmpMap["findinames"].([]interface{}) {
		findinames = append(findinames, v.(string))
	}

	// c = configer.New(dir + "Strategy.yaml")
	c = configer.New(dir + StgConfile)

	err = c.Load()
	if err != nil {
		panic(err)
	}
	err = c.Unmarshal()
	if err != nil {
		panic(err)
	}

	SNum := c.GetFloat64("default.snum")
	FNum := c.GetFloat64("default.fnum")
	STimeCritic := c.GetString("default.stimecritic")
	FTimeCritic := c.GetString("default.ftimecritic")
	return NewDMTStrategy(sinstrnames, sindinames, finstrnames, findinames, SNum, FNum, STimeCritic, FTimeCritic)
}

// getTimeValue get the time value which is a stirng fromt a time string
// func getTimeValue(timeString string) string {
// 	tmp := strings.Fields(timeString)
// 	return tmp[1]
// }

func getTimeValue(timeString string) string {
	// get the time value
	timeValue := timeString[11:16]
	return timeValue
}

// change the dmt.stimeCondition and dmt.ftimeCondition

// simple strategy: only market data and indicators
func (dmt *DMTStrategy) ActOnData(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap, Eval func([]float64) []float64) (orderRes OrderResult) {
	// 判断股票标的切片SInstrNames是否为空，如果为空，则不操作股票数据循环
	if len(dmt.SInstNames) != 0 {

		// 依据标的循环Data得到数据
		for instID, SBDE := range bc.Stockdata {
			// 判断是否数据为NaN，如果为NaN，则跳过
			if !ContainNaN(SBDE.IndiDataMap) {
				tmpSCP := cp.SimpleNewSCPFromMap(CPMap, instID)
				// 展示常规垃圾结果请将下列注释恢复并注释上面的一行
				// tmpSCP := CPMap.StockPropMap[instID]
				var GEPSlice = make([]float64, len(dmt.SIndiNames))
				for i := 0; i < len(dmt.SIndiNames); i++ {
					GEPSlice[i] = SBDE.IndiDataMap[dmt.SIndiNames[i]]
				}
				// 针对slice评估得到数值
				// do sth about GEPSlice
				tmps := Eval(GEPSlice)
				// DCE: debug info
				if debug {
					// 针对数值进行策略映射：什么时间条件、标的、操作方向、多少、价格
					// check if the tmps[0] is NaN
					if math.IsNaN(tmps[0]) {
						log.Info().Float64("Eval", tmps[0]).Msg("Math Related tmp info NaN")
					}
					// check if the tmps[0] is inf
					if math.IsInf(tmps[0], 0) {
						log.Info().Float64("Eval", tmps[0]).Msg("Math Related tmp info inf")
					}
				}
				if tmps[0] >= 0 {
					if _, ok := vAcct.SAcct.PosMap[instID]; ok {
						// 注意，下面这行是判断如果有持仓就不买了。为配合样例先注释了吧，今后记得改回来。尽管你一定会忘
						// if vAcct.SAcct.PosMap[instID].CalEquity() == 0 {
						// 	orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, datetime, SBDE.IndiDataMap["close"], ss.SNum, order.Buy, &tmpSCP))
						// }
						orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["close"], dmt.SNum, "Buy", &tmpSCP))

					} else {
						// I know! it's for you to do sth more meaningful
						orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["close"], dmt.SNum, "Buy", &tmpSCP))

					}
					// DCE: debug info
					if debug {
						// this part is for test only
						log.Info().Str("Account UUID", vAcct.SAcct.UUID).Str("TimeStamp", datetime).
							Float64("Eval", tmps[0]).Str("InstID", instID).
							Msg("Strategy buy")
					}
				}
				if tmps[0] < 0 {
					//check if target is in the vAcct.SAcct, if yes, sell them all if not, do nothing
					if _, ok := vAcct.SAcct.PosMap[instID]; ok {
						if vAcct.SAcct.PosMap[instID].CalPosPrevNum() > 0 {
							orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["close"], vAcct.SAcct.PosMap[instID].CalPosPrevNum(), "Sell", &tmpSCP))
						}
					}
					// DCE: debug info
					if debug {
						// this part is for test only
						log.Info().Str("Account UUID", vAcct.SAcct.UUID).Str("TimeStamp", datetime).
							Float64("Eval", tmps[0]).Str("InstID", instID).
							Msg("Strategy sell")
					}
				}
			}
		}
	}

	// 判断期货标的切片FInstrNames是否为空，如果为空，则不操作期货数据循环

	if len(dmt.FInstNames) != 0 {
		for instrID, FBDE := range bc.Futuresdata {
			// 判断是否数据为NaN，如果为NaN，则跳过
			if !ContainNaN(FBDE.IndiDataMap) {
				tmpFCP := CPMap.FuturesPropMap[instrID]
				var GEPSlice = make([]float64, len(dmt.FIndiNames))
				for i := 0; i < len(dmt.FIndiNames); i++ {
					GEPSlice[i] = FBDE.IndiDataMap[dmt.FIndiNames[i]]
				}
				tmpf := Eval(GEPSlice)

				// 期货操作逻辑实现当日开仓 当日平仓  时间超过FActTime 平仓  不再开仓
				// 如果存在isrolloverday字段，则表示移仓换月日，规则可能不同 是否单独处理的自由度交给用户

				// if val, ok := FBDE.DataMap["isrolloverday"];ok {
				//
				// }
				tdynuml, tdynums := vAcct.FAcct.PosMap[instrID].CalPosTdyNum()
				prevnuml, prevnums := vAcct.FAcct.PosMap[instrID].CalPosPrevNum()
				if tmpf[0] >= 0 {
					// 清仓今日昨日空头持仓
					if tdynums > 0 {
						orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], tdynums, "Buy", "CloseToday", &tmpFCP))
					}
					if prevnums > 0 {
						orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], prevnums, "Buy", "ClosePrevious", &tmpFCP))
					}
					// 这里应该有一个期货账户资金使用率的控制  一个比例
					// moneyoneachf := vAcct.FAcct.MktVal / float64(len(ss.FInstNames))
					// tmpnum := math.Floor(moneyoneachf * ss.FuturesFundUseRate / FBDE.FBarData.Close / CPMap.FuturesPropMap[instrID].ContractSize / (CPMap.FuturesPropMap[instrID].MarginLong + CPMap.FuturesPropMap[instrID].MarginBroker))

					orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], dmt.FNum, "Buy", "Open", &tmpFCP))

				}

				if tmpf[0] < 0 {
					if tdynuml > 0 {
						orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], tdynuml, "Sell", "CloseToday", &tmpFCP))
					}
					if prevnuml > 0 {
						orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], prevnuml, "Sell", "ClosePrevious", &tmpFCP))
					}
					// 这里应该有一个期货账户资金使用率的控制  一个比例
					// moneyoneachf := vAcct.FAcct.MktVal / float64(len(ss.FInstNames))
					// tmpnum := math.Floor(moneyoneachf * ss.FuturesFundUseRate / FBDE.FBarData.Close / CPMap.FuturesPropMap[instrID].ContractSize / (CPMap.FuturesPropMap[instrID].MarginShort + CPMap.FuturesPropMap[instrID].MarginBroker))

					orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], dmt.FNum, "Sell", "Open", &tmpFCP))
				}
			}
		}
	}
	return orderRes
}

func (dmt *DMTStrategy) ActOnDataMAN(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap) (orderRes OrderResult) {
	// 时间准则判定 stock
	if getTimeValue(datetime) == dmt.STimeCritic && !dmt.ifsdone {
		dmt.stimeCondition = true
	} else if dmt.ifsdone {
		dmt.stimeCondition = false
		dmt.ifsdone = false
	}
	// 时间准则判定 futures
	if getTimeValue(datetime) == dmt.FTimeCritic && !dmt.iffdone {
		dmt.ftimeCondition = true
	} else if dmt.iffdone {
		dmt.ftimeCondition = false
		dmt.iffdone = false
	}
	// 判断股票标的切片SInstrNames是否为空 并且 时间准则为真，如果为空，则不操作股票数据循环
	if len(dmt.SInstNames) != 0 && dmt.stimeCondition {

		// 依据标的循环Data得到数据
		for instID, SBDE := range bc.Stockdata {
			// 判断是否数据为NaN，如果为NaN，则跳过
			if !ContainNaN(SBDE.IndiDataMap) {
				tmpSCP := cp.SimpleNewSCPFromMap(CPMap, instID)
				buyval := SBDE.IndiDataMap["Close"] - SBDE.IndiDataMap["Amount"]/SBDE.IndiDataMap["Volume"]
				lstbv, bok := dmt.lastBuyValue[instID]
				sellval := SBDE.IndiDataMap["Close"] - SBDE.IndiDataMap["Open"]
				lstsv, sok := dmt.lastSellValue[instID]
				if bok && buyval >= lstbv {
					if _, ok := vAcct.SAcct.PosMap[instID]; ok {
						// 注意，下面这行是判断如果有持仓就不买了。为配合样例先注释了吧，今后记得改回来。尽管你一定会忘
						// if vAcct.SAcct.PosMap[instID].CalEquity() == 0 {
						// 	orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, datetime, SBDE.IndiDataMap["close"], ss.SNum, order.Buy, &tmpSCP))
						// }
						orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["Close"], dmt.SNum, "Buy", &tmpSCP))

					} else {
						orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["Close"], dmt.SNum, "Buy", &tmpSCP))
					}
					dmt.ifsdone = true
					// DCE: debug info
					if debug {
						// this part is for test only
						log.Info().Str("Account UUID", vAcct.SAcct.UUID).Str("TimeStamp", datetime).
							Float64("Close", SBDE.IndiDataMap["Close"]).Float64("Open", SBDE.IndiDataMap["Open"]).Str("InstID", instID).
							Msg("Strategy buy")
					}
				}
				// update the last buy value
				dmt.lastBuyValue[instID] = buyval
				if sok && sellval <= lstsv {
					//check if target is in the vAcct.SAcct, if yes, sell them all if not, do nothing
					if _, ok := vAcct.SAcct.PosMap[instID]; ok {
						if vAcct.SAcct.PosMap[instID].CalPosPrevNum() > 0 {
							orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["Close"], vAcct.SAcct.PosMap[instID].CalPosPrevNum(), "Sell", &tmpSCP))
						}
					}
					dmt.ifsdone = true
					// DCE: debug info
					if debug {
						// this part is for test only
						log.Info().Str("Account UUID", vAcct.SAcct.UUID).Str("TimeStamp", datetime).
							Float64("Close", SBDE.IndiDataMap["Close"]).Float64("Open", SBDE.IndiDataMap["Open"]).Str("InstID", instID).
							Msg("Strategy sell")
					}
				}
				// update the last sell value
				dmt.lastSellValue[instID] = sellval

			}

		}

	}

	// 判断期货标的切片FInstrNames是否为空，如果为空，则不操作期货数据循环

	if len(dmt.FInstNames) != 0 && dmt.ftimeCondition {
		for instrID, FBDE := range bc.Futuresdata {
			// 判断是否数据为NaN，如果为NaN，则跳过
			if !ContainNaN(FBDE.IndiDataMap) {
				tmpFCP := CPMap.FuturesPropMap[instrID]
				buyval := FBDE.IndiDataMap["close"] - FBDE.IndiDataMap["open"]
				lstbv, bok := dmt.lastBuyValue[instrID]
				sellval := FBDE.IndiDataMap["close"] - FBDE.IndiDataMap["open"]
				lstsv, sok := dmt.lastSellValue[instrID]

				// var GEPSlice = make([]float64, len(ss.FIndiNames))
				// for i := 0; i < len(ss.FIndiNames); i++ {
				// 	GEPSlice[i] = FBDE.IndiDataMap[ss.FIndiNames[i]]
				// }
				// tmpf := Eval(GEPSlice)

				// 期货操作逻辑实现当日开仓 当日平仓  时间超过FActTime 平仓  不再开仓
				// 如果存在isrolloverday字段，则表示移仓换月日，规则可能不同 是否单独处理的自由度交给用户

				// if val, ok := FBDE.DataMap["isrolloverday"];ok {
				//
				// }
				tdynuml, tdynums := vAcct.FAcct.PosMap[instrID].CalPosTdyNum()
				prevnuml, prevnums := vAcct.FAcct.PosMap[instrID].CalPosPrevNum()
				if bok && buyval >= lstbv {
					// 清仓今日昨日空头持仓
					if tdynums > 0 {
						orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], tdynums, "Buy", "CloseToday", &tmpFCP))
					}
					if prevnums > 0 {
						orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], prevnums, "Buy", "ClosePrevious", &tmpFCP))
					}
					// 这里应该有一个期货账户资金使用率的控制  一个比例
					// moneyoneachf := vAcct.FAcct.MktVal / float64(len(ss.FInstNames))
					// tmpnum := math.Floor(moneyoneachf * ss.FuturesFundUseRate / FBDE.FBarData.Close / CPMap.FuturesPropMap[instrID].ContractSize / (CPMap.FuturesPropMap[instrID].MarginLong + CPMap.FuturesPropMap[instrID].MarginBroker))

					orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], dmt.FNum, "Buy", "Open", &tmpFCP))
					dmt.iffdone = true
				}
				// update the last buy value
				dmt.lastBuyValue[instrID] = buyval

				if sok && sellval <= lstsv {
					if tdynuml > 0 {
						orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], tdynuml, "Sell", "CloseToday", &tmpFCP))
					}
					if prevnuml > 0 {
						orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], prevnuml, "Sell", "ClosePrevious", &tmpFCP))
					}
					// 这里应该有一个期货账户资金使用率的控制  一个比例
					// moneyoneachf := vAcct.FAcct.MktVal / float64(len(ss.FInstNames))
					// tmpnum := math.Floor(moneyoneachf * ss.FuturesFundUseRate / FBDE.FBarData.Close / CPMap.FuturesPropMap[instrID].ContractSize / (CPMap.FuturesPropMap[instrID].MarginShort + CPMap.FuturesPropMap[instrID].MarginBroker))

					orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], dmt.FNum, "Sell", "Open", &tmpFCP))
					dmt.iffdone = true
				}
				// update the last sell value
				dmt.lastSellValue[instrID] = sellval
			}
		}
	}
	return orderRes
}
