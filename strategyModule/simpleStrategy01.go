package strategyModule

// // market timing strategy
// // one expression for one strategy, therefore one market for one strategy
// import (
// 	"fmt"
// 	"math"

// 	"account/virtualaccount"

// 	cp "contractproperty"

// 	"dataprocessor"

// 	"order"

// 	"github.com/rs/zerolog/log"

// 	"configer"
// 	"time"
// 	"errors"
// )

// type SimpleStrategy struct {
// 	SInstNames []string // 股票标的名称
// 	SIndiNames []string // 股票参与GEP指标名称，注意其数量不大于BarDE内信息数量，且strategy内可见BarDE的数据
// 	SNum       float64  // 股票标的数量
// 	SHoldNum   []float64 //每个股票标的对应的初始持仓数 和SInstNames中的名称一一对应
// 	TradeLimit int //每天限制的交易次数 初步定义为4
// 	TradeCounter int //已进行的交易次数，每次交易+1，不超过TradeLimit
// 	StartTime  string //起始时间 开始根据策略进行买卖操作 读csv数据的第一列
// 	StopTime   string //终止时间 不再根据策略进行买卖操作
// 	FInstNames []string // 期货标的名称
// 	FIndiNames []string // 期货参与GEP指标名称， there should be a rollover field in the futures IndiDataNames slice
// 	FNum       float64  // 期货标的数量
// }

// func NewSimpleStrategy(SInstNms, SIndiNms, FInstNms, FIndiNms []string,Starttime,Stoptime string,Snum, Fnum float64,Sholdnum []float64,Tradelimit,Tradecounter int) SimpleStrategy {
// 	return SimpleStrategy{
// 		SInstNames: SInstNms,
// 		SIndiNames: SIndiNms,
// 		SNum:       Snum,
// 		SHoldNum:   Sholdnum,
// 		TradeLimit: Tradelimit,
// 		TradeCounter: Tradecounter,
// 		StartTime: Starttime,
// 		StopTime: Stoptime,
// 		FInstNames: FInstNms,
// 		FIndiNames: FIndiNms,
// 		FNum:       Fnum,
// 	}
// }

// // this function is nessary for the framework
// //作用是从配置文件中读取所需要的参量，返回到SimpleStrategy里。
// func NewSimpleStrategyFromConfig(sec string, dir string,Starttime,Stoptime,string) SimpleStrategy {
// 	c := configer.New(dir + "BackTest.yaml") //构建配置文件
// 	err := c.Load()
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = c.Unmarshal()  //JSON解码函数
// 	if err != nil {
// 		panic(err)
// 	}
// 	tmpMap := c.GetStringMap(sec) //读取配置
// 	//读取配置文件里的股票标的名称、股票参与GEP指标名称、期货标的名称、期货参与GEP指标名称，放入相应数组
// 	var sinstrnames []string
// 	for _, v := range tmpMap["sinstrnames"].([]interface{}) {
// 		sinstrnames = append(sinstrnames, v.(string))
// 	}
// 	var sindinames []string
// 	for _, v := range tmpMap["sindinames"].([]interface{}) {
// 		sindinames = append(sindinames, v.(string))
// 	}

// 	var finstrnames []string
// 	for _, v := range tmpMap["finstrnames"].([]interface{}) {
// 		finstrnames = append(finstrnames, v.(string))
// 	}
// 	var findinames []string
// 	for _, v := range tmpMap["findinames"].([]interface{}) {
// 		findinames = append(findinames, v.(string))
// 	}

// 	c = configer.New(dir + "Strategy.yaml")
// 	err = c.Load()
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = c.Unmarshal()
// 	if err != nil {
// 		panic(err)
// 	}

// 	SNum := c.GetFloat64("default.snum")
// 	FNum := c.GetFloat64("default.fnum")

// 	Tradelimit:=4 //限制交易次数最多为4
// 	Tradecounter:=0 //初始次数为0
// 	var Sholdnum []float64
// 	for _, v := range tmpMap["sinstrnames"].([]interface{}) {
// 		Sholdnum = append(Sholdnum, v.(string))
// 	}
// 	Starttime := "2017/10/9 9:50"
// 	Stoptime := "2017/10/9 10:50"

// 	return SimpleStrategy{
// 		SInstNames: sinstrnames,
// 		SIndiNames: sindinames,
// 		SNum:       SNum,
// 		SHoldNum:   Sholdnum,
// 		TradeLimit: Tradelimit,
// 		TradeCounter: Tradecounter,
// 		StartTime: Starttime,
// 		StopTime: Stoptime,
// 		FInstNames: finstrnames,
// 		FIndiNames: findinames,
// 		FNum:       FNum,
// 	}
// }

// // simple strategy: only market data and indicators
// func (ss SimpleStrategy) ActOnData(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap, Eval func([]float64) []float64) (orderRes OrderResult) {
// 	// 判断股票标的切片SInstrNames是否为空，如果为空，则不操作股票数据循环
// 	if len(ss.SInstNames) != 0 {

// 		// 依据标的循环Data得到数据
// 		//instID标的代码 bc:将股票数据和期货数据进行了组合 bc.stockdata:取出其中股票数据
// 		for instID, SBDE := range bc.Stockdata {
// 			// 判断是否数据为NaN，如果为NaN，则跳过
// 			//IndiDataMap:所有信息 key :indicator and bar element name , value: indicator and bar element value
// 			if !ContainNaN(SBDE.IndiDataMap) {
// 				tmpSCP := cp.SimpleNewSCPFromMap(CPMap, instID) //这个cp不太懂？
// 				// 展示常规垃圾结果请将下列注释恢复并注释上面的一行
// 				// tmpSCP := CPMap.StockPropMap[instID]
// 				//把股票参与GEP指标名称放入GEPSlice
// 				var GEPSlice = make([]float64, len(ss.SIndiNames))
// 				for i := 0; i < len(ss.SIndiNames); i++ {
// 					GEPSlice[i] = SBDE.IndiDataMap[ss.SIndiNames[i]]
// 				}
// 				// 针对slice评估得到数值
// 				// do sth about GEPSlice
// 				tmps := Eval(GEPSlice)

// 				// 针对数值进行策略映射：什么时间条件、标的、操作方向、多少、价格
// 				// check if the tmps[0] is NaN
// 				if math.IsNaN(tmps[0]) {
// 					fmt.Println(tmps[0])
// 				}
// 				// check if the tmps[0] is inf
// 				if math.IsInf(tmps[0], 0) {
// 					fmt.Println(tmps[0])
// 				}

// 				if tmps[0] >= 0 {
// 					if _, ok := vAcct.SAcct.PosMap[instID]; ok {
// 						// 注意，下面这行是判断如果有持仓就不买了。为配合样例先注释了吧，今后记得改回来。尽管你一定会忘
// 						// if vAcct.SAcct.PosMap[instID].CalEquity() == 0 {
// 						// 	orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, datetime, SBDE.IndiDataMap["close"], ss.SNum, order.Buy, &tmpSCP))
// 						// }
// 						orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["close"], ss.SNum, "Buy", &tmpSCP))

// 					} else {
// 						// I know! it's for you to do sth more meaningful
// 						orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["close"], ss.SNum, "Buy", &tmpSCP))

// 					}
// 					// this part is for test only
// 					log.Info().Str("Account UUID", vAcct.SAcct.UUID).Str("TimeStamp", datetime).
// 						Float64("Eval", tmps[0]).Str("InstID", instID).
// 						Msg("Strategy buy")

// 				}
// 				if tmps[0] < 0 {
// 					//check if target is in the vAcct.SAcct, if yes, sell them all if not, do nothing
// 					if _, ok := vAcct.SAcct.PosMap[instID]; ok {
// 						if vAcct.SAcct.PosMap[instID].CalPosPrevNum() > 0 {
// 							orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["close"], vAcct.SAcct.PosMap[instID].CalPosPrevNum(), "Sell", &tmpSCP))
// 						}
// 					}

// 					// this part is for test only
// 					log.Info().Str("Account UUID", vAcct.SAcct.UUID).Str("TimeStamp", datetime).
// 						Float64("Eval", tmps[0]).Str("InstID", instID).
// 						Msg("Strategy sell")
// 				}
// 			}
// 		}
// 	}

// 	// 判断期货标的切片FInstrNames是否为空，如果为空，则不操作期货数据循环

// 	if len(ss.FInstNames) != 0 {
// 		for instrID, FBDE := range bc.Futuresdata {
// 			// 判断是否数据为NaN，如果为NaN，则跳过
// 			if !ContainNaN(FBDE.IndiDataMap) {
// 				tmpFCP := CPMap.FuturesPropMap[instrID]
// 				var GEPSlice = make([]float64, len(ss.FIndiNames))
// 				for i := 0; i < len(ss.FIndiNames); i++ {
// 					GEPSlice[i] = FBDE.IndiDataMap[ss.FIndiNames[i]]
// 				}
// 				tmpf := Eval(GEPSlice)

// 				// 期货操作逻辑实现当日开仓 当日平仓  时间超过FActTime 平仓  不再开仓
// 				// 如果存在isrolloverday字段，则表示移仓换月日，规则可能不同 是否单独处理的自由度交给用户

// 				// if val, ok := FBDE.DataMap["isrolloverday"];ok {
// 				//
// 				// }
// 				tdynuml, tdynums := vAcct.FAcct.PosMap[instrID].CalPosTdyNum() //判断今日有无持仓
// 				prevnuml, prevnums := vAcct.FAcct.PosMap[instrID].CalPosPrevNum() //判断一下之前有没有持仓
// 				if tmpf[0] >= 0 {
// 					// 清仓今日昨日空头持仓
// 					if tdynums > 0 {
// 						//买入平今仓 指的就是平掉今天的仓位，指的是期货日内交易
// 						orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], tdynums, "Buy", "CloseToday", &tmpFCP))
// 					}
// 					if prevnums > 0 {
// 						//买入平昨仓 指的就是非日内交易，平掉以前的历史仓位
// 						orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], prevnums, "Buy", "ClosePrevious", &tmpFCP))
// 					}
// 					// 这里应该有一个期货账户资金使用率的控制  一个比例
// 					// moneyoneachf := vAcct.FAcct.MktVal / float64(len(ss.FInstNames))
// 					// tmpnum := math.Floor(moneyoneachf * ss.FuturesFundUseRate / FBDE.FBarData.Close / CPMap.FuturesPropMap[instrID].ContractSize / (CPMap.FuturesPropMap[instrID].MarginLong + CPMap.FuturesPropMap[instrID].MarginBroker))
// 					//都没有就开仓买
// 					orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], ss.FNum, "Buy", "Open", &tmpFCP))

// 				}

// 				if tmpf[0] < 0 {
// 					if tdynuml > 0 {
// 						//卖出平今仓
// 						orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], tdynuml, "Sell", "CloseToday", &tmpFCP))
// 					}
// 					if prevnuml > 0 {
// 						//卖出平昨仓
// 						orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], prevnuml, "Sell", "ClosePrevious", &tmpFCP))
// 					}
// 					// 这里应该有一个期货账户资金使用率的控制  一个比例
// 					// moneyoneachf := vAcct.FAcct.MktVal / float64(len(ss.FInstNames))
// 					// tmpnum := math.Floor(moneyoneachf * ss.FuturesFundUseRate / FBDE.FBarData.Close / CPMap.FuturesPropMap[instrID].ContractSize / (CPMap.FuturesPropMap[instrID].MarginShort + CPMap.FuturesPropMap[instrID].MarginBroker))

// 					orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(instrID, false, false, datetime, FBDE.IndiDataMap["Close"], ss.FNum, "Sell", "Open", &tmpFCP))
// 				}
// 			}
// 		}
// 	}
// 	return orderRes
// }

// func (ss SimpleStrategy) ActOnDataMAN(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap) (orderRes OrderResult) {
// 	start, err1 := time.Parse("2006-01-02 15:04:05",Starttime)
// 	if err1 != nil {
// 		panic("Starttime parse error")
// 	}
// 	stop, err2 := time.Parse("2006-01-02 15:04:05",Stoptime)
// 	if err2 != nil {
// 		panic("Stoptime parse error")
// 	}
// 	time,err3 := time.Parse("2006-01-02 15:04:05",datetime)
// 	if err3 != nil {
// 		panic("Datetime parse error")
// 	}
// 	TradeN:=math.Floor(ss.SHoldNum[0]/float64(ss.TradeLimit)) //这里选择了第一支股票SHoldNum[0] TradeN为每次交易数量
// 	if ss.TradeCounter<ss.TradeLimit {
// 		if orderRes.StockOrderS{
// 			return
// 		}
// 		if err1==nil && err2==nil &&err3==nill && start.Before(time) && time.Before(stop) {
// 			// 判断股票标的切片SInstrNames是否为空，如果为空，则不操作股票数据循环
// 			if len(ss.SInstNames) != 0 {

// 				// 依据标的循环Data得到数据
// 				for instID, SBDE := range bc.Stockdata {
// 					// 判断是否数据为NaN，如果为NaN，则跳过
// 					if !ContainNaN(SBDE.IndiDataMap) {
// 						tmpSCP := cp.SimpleNewSCPFromMap(CPMap, instID)

// 						if SBDE.IndiDataMap["close"]-SBDE.IndiDataMap["ma3"] >= 0 {
// 							if _, ok := vAcct.SAcct.PosMap[instID]; ok {
// 								// 注意，下面这行是判断如果有持仓就才操作
// 								if vAcct.SAcct.PosMap[instID].CalEquity() != 0 {
// 									orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["close"], TradeN, "Buy", &tmpSCP))
// 									ss.TradeCounter+=1
// 								}
// 							}
// 							// this part is for test only
// 							log.Info().Str("Account UUID", vAcct.SAcct.UUID).Str("TimeStamp", datetime).
// 								Float64("close", SBDE.IndiDataMap["close"]).Float64("ma3", SBDE.IndiDataMap["ma3"]).Str("InstID", instID).
// 								Msg("Strategy buy")

// 						}
// 						if SBDE.IndiDataMap["close"]-SBDE.IndiDataMap["ma3"] < 0 {
// 							//check if target is in the vAcct.SAcct, if yes, sell them all if not, do nothing
// 							if _, ok := vAcct.SAcct.PosMap[instID]; ok {
// 								if vAcct.SAcct.PosMap[instID].CalPosPrevNum() > 0 {
// 									orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["close"], TradeN, "Sell", &tmpSCP))
// 									ss.TradeCounter+=1
// 								}
// 							}

// 							// this part is for test only
// 							log.Info().Str("Account UUID", vAcct.SAcct.UUID).Str("TimeStamp", datetime).
// 								Float64("close", SBDE.IndiDataMap["close"]).Float64("ma3", SBDE.IndiDataMap["ma3"]).Str("InstID", instID).
// 								Msg("Strategy sell")
// 						}

// 					}

// 				}

// 			}
// 		}
// 		if err1==nil && err2==nil &&err3==nill && stop.Before(time) {
// 			if orderRes.StockOrderS < ss.SHoldNum[0] {
// 				orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["close"], ss.SHoldNum[0]-orderRes.StockOrderS, "Buy", &tmpSCP))
// 			}
// 			if orderRes.StockOrderS > ss.SHoldNum[0] {
// 				orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["close"], orderRes.StockOrderS-ss.SHoldNum[0], "Sell", &tmpSCP))
// 			}
// 		}
// 	}
// 	return orderRes
// }
