package strategyModule

// market timing strategy
// one expression for one strategy, therefore one market for one strategy
import (
	"math"

	"github.com/rs/zerolog/log"
	"github.com/wonderstone/QuantTools/account/virtualaccount"
	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/dataprocessor"
	"github.com/wonderstone/QuantTools/order"

	"github.com/wonderstone/QuantTools/configer"
)

type SimpleStrategy struct {
	SInstNames []string // 股票标的名称
	SIndiNames []string // 股票参与GEP指标名称，注意其数量不大于BarDE内信息数量，且strategy内可见BarDE的数据
	SNum       float64  // 股票标的数量
	FInstNames []string // 期货标的名称
	FIndiNames []string // 期货参与GEP指标名称， there should be a rollover field in the futures IndiDataNames slice
	FNum       float64  // 期货标的数量
}

func NewSimpleStrategy(SInstNms, SIndiNms, FInstNms, FIndiNms []string, Snum, Fnum float64) SimpleStrategy {
	return SimpleStrategy{
		SInstNames: SInstNms,
		SIndiNames: SIndiNms,
		SNum:       Snum,
		FInstNames: FInstNms,
		FIndiNames: FIndiNms,
		FNum:       Fnum,
	}
}

// this function is nessary for the framework
func NewSimpleStrategyFromConfig(dir string, file string, sec string, StgConfile string) SimpleStrategy {
	// c := configer.New(dir + "BackTest.yaml")
	c := configer.New(dir + file)
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

	SNum := c.GetFloat64(sec + ".snum")
	FNum := c.GetFloat64(sec + ".fnum")

	return SimpleStrategy{
		SInstNames: sinstrnames,
		SIndiNames: sindinames,
		SNum:       SNum,
		FInstNames: finstrnames,
		FIndiNames: findinames,
		FNum:       FNum,
	}
}

// simple strategy: only market data and indicators
func (ss SimpleStrategy) ActOnData(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap, Eval func([]float64) []float64) (orderRes OrderResult) {
	// 判断股票标的切片SInstrNames是否为空，如果为空，则不操作股票数据循环
	if len(ss.SInstNames) != 0 {

		// 依据标的循环Data得到数据
		for instID, SBDE := range bc.Stockdata {
			// 判断是否数据为NaN，如果为NaN，则跳过
			if !ContainNaN(SBDE.IndiDataMap) {
				tmpSCP := cp.SimpleNewSCPFromMap(CPMap, instID)
				// 展示常规垃圾结果请将下列注释恢复并注释上面的一行
				// tmpSCP := CPMap.StockPropMap[instID]
				var GEPSlice = make([]float64, len(ss.SIndiNames))
				for i := 0; i < len(ss.SIndiNames); i++ {
					GEPSlice[i] = SBDE.IndiDataMap[ss.SIndiNames[i]]
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
						orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["close"], ss.SNum, "Buy", &tmpSCP))

					} else {
						// I know! it's for you to do sth more meaningful
						orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["close"], ss.SNum, "Buy", &tmpSCP))

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

	if len(ss.FInstNames) != 0 {
		for instrID, FBDE := range bc.Futuresdata {
			// 判断是否数据为NaN，如果为NaN，则跳过
			if !ContainNaN(FBDE.IndiDataMap) {
				tmpFCP := CPMap.FuturesPropMap[instrID]
				var GEPSlice = make([]float64, len(ss.FIndiNames))
				for i := 0; i < len(ss.FIndiNames); i++ {
					GEPSlice[i] = FBDE.IndiDataMap[ss.FIndiNames[i]]
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

					orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], ss.FNum, "Buy", "Open", &tmpFCP))

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

					orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], ss.FNum, "Sell", "Open", &tmpFCP))
				}
			}
		}
	}
	return orderRes
}

func (ss SimpleStrategy) ActOnDataMAN(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap) (orderRes OrderResult) {
	// 判断股票标的切片SInstrNames是否为空，如果为空，则不操作股票数据循环
	if len(ss.SInstNames) != 0 {

		// 依据标的循环Data得到数据
		for instID, SBDE := range bc.Stockdata {
			// 判断是否数据为NaN，如果为NaN，则跳过
			if !ContainNaN(SBDE.IndiDataMap) {
				tmpSCP := cp.SimpleNewSCPFromMap(CPMap, instID)
				// 展示常规垃圾结果请将下列注释恢复并注释上面的一行
				// tmpSCP := CPMap.StockPropMap[instID]
				// var GEPSlice = make([]float64, len(ss.SIndiNames))
				// for i := 0; i < len(ss.SIndiNames); i++ {
				// 	GEPSlice[i] = SBDE.IndiDataMap[ss.SIndiNames[i]]
				// }
				// 针对slice评估得到数值
				// do sth about GEPSlice
				// tmps := Eval(GEPSlice)

				// 针对数值进行策略映射：什么时间条件、标的、操作方向、多少、价格
				// check if the tmps[0] is NaN
				// if math.IsNaN(tmps[0]) {
				// }
				// // check if the tmps[0] is inf
				// if math.IsInf(tmps[0], 0) {
				// }

				if SBDE.IndiDataMap["close"]-SBDE.IndiDataMap["ma1"] >= 0 {
					if _, ok := vAcct.SAcct.PosMap[instID]; ok {
						// 注意，下面这行是判断如果有持仓就不买了。为配合样例先注释了吧，今后记得改回来。尽管你一定会忘
						// if vAcct.SAcct.PosMap[instID].CalEquity() == 0 {
						// 	orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, datetime, SBDE.IndiDataMap["close"], ss.SNum, order.Buy, &tmpSCP))
						// }
						orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["close"], ss.SNum, "Buy", &tmpSCP))

					} else {
						orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["close"], ss.SNum, "Buy", &tmpSCP))
					}
					// DCE: debug info
					if debug {
						// this part is for test only
						log.Info().Str("Account UUID", vAcct.SAcct.UUID).Str("TimeStamp", datetime).
							Float64("close", SBDE.IndiDataMap["close"]).Float64("ma1", SBDE.IndiDataMap["ma1"]).Str("InstID", instID).
							Msg("Strategy buy")
					}

				}
				if SBDE.IndiDataMap["close"]-SBDE.IndiDataMap["ma1"] < 0 {
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
							Float64("close", SBDE.IndiDataMap["close"]).Float64("ma1", SBDE.IndiDataMap["ma1"]).Str("InstID", instID).
							Msg("Strategy sell")
					}
				}

			}

		}

	}

	// 判断期货标的切片FInstrNames是否为空，如果为空，则不操作期货数据循环

	if len(ss.FInstNames) != 0 {
		for instrID, FBDE := range bc.Futuresdata {
			// 判断是否数据为NaN，如果为NaN，则跳过
			if !ContainNaN(FBDE.IndiDataMap) {
				tmpFCP := CPMap.FuturesPropMap[instrID]
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
				if FBDE.IndiDataMap["close"]-FBDE.IndiDataMap["open"] >= 0 {
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

					orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], ss.FNum, "Buy", "Open", &tmpFCP))
				}

				if FBDE.IndiDataMap["close"]-FBDE.IndiDataMap["open"] < 0 {
					if tdynuml > 0 {
						orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], tdynuml, "Sell", "CloseToday", &tmpFCP))
					}
					if prevnuml > 0 {
						orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], prevnuml, "Sell", "ClosePrevious", &tmpFCP))
					}
					// 这里应该有一个期货账户资金使用率的控制  一个比例
					// moneyoneachf := vAcct.FAcct.MktVal / float64(len(ss.FInstNames))
					// tmpnum := math.Floor(moneyoneachf * ss.FuturesFundUseRate / FBDE.FBarData.Close / CPMap.FuturesPropMap[instrID].ContractSize / (CPMap.FuturesPropMap[instrID].MarginShort + CPMap.FuturesPropMap[instrID].MarginBroker))

					orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], ss.FNum, "Sell", "Open", &tmpFCP))
				}
			}
		}
	}
	return orderRes
}
