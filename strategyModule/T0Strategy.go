package strategyModule

import (
	"strings"

	"github.com/wonderstone/QuantTools/account/stockaccount"
	"github.com/wonderstone/QuantTools/account/virtualaccount"
	"github.com/wonderstone/QuantTools/configer"
	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/dataprocessor"
	"github.com/wonderstone/QuantTools/order"
)

// Stock T0Strategy

type ST0Strategy struct {
	IfBT          bool               // ~ 是否为回测/实盘状态 T0特有问题 回测第一天为资金状态 需要买入标的
	initBTstate   map[string]bool    // ~ 是否已经初始化了标的回测买入状态(内部)
	InstNames     []string           // 股票标的名称
	IndiNames     []string           // 股票参与GEP指标名称，注意其数量不大于BarDE内信息数量，且strategy内可见BarDE的数据
	HoldNums      []int              // ~ 股票标的持仓数量, 外部获取
	holdNumMap    map[string]int     // ~ 股票标的持仓数量, key is the stock name,get from SInstNames and SHoldNums(内部)
	Tlimit        int                // 每日交易次数限制
	tCounter      map[string]int     // 交易次数计数器(内部)
	tState        map[string]int     // ~ 交易状态(内部) 0:未买入 1:已买入 -1:已卖出
	InitBuyTime   string             // ~ 初始化买入时间
	StartTime     string             // ~ 交易开始时间
	StopTime      string             // ~ 交易结束时间
	ifReCoverMap  map[string]bool    // 是否恢复了股票持仓(内部)
	lastDate      string             // 上一交易日日期(内部)
	lastBuyValue  map[string]float64 // 公式上一次买入信号数值
	lastSellValue map[string]float64 // 公式上一次卖出信号数值
}

func NewST0Strategy(IfBT bool, InstNames, IndiNames []string, HoldNums []int, Tlimit int, InitBuyTime, StartTime, StopTime string) ST0Strategy {
	initBTstate := make(map[string]bool)
	holdNumMap := make(map[string]int)
	tCounter := make(map[string]int)
	tState := make(map[string]int)
	ifReCoverMap := make(map[string]bool)
	// 通过数据结构记录状态是常用方式
	for i, name := range InstNames {
		initBTstate[name] = false // # 初始都为false
		holdNumMap[name] = HoldNums[i]
		tCounter[name] = 0
		tState[name] = 0
		ifReCoverMap[name] = false
	}
	return ST0Strategy{
		IfBT:          IfBT,
		initBTstate:   initBTstate,
		InstNames:     InstNames,
		IndiNames:     IndiNames,
		HoldNums:      HoldNums,
		holdNumMap:    holdNumMap,
		Tlimit:        Tlimit,
		tCounter:      tCounter,
		tState:        tState,
		InitBuyTime:   InitBuyTime,
		StartTime:     StartTime,
		StopTime:      StopTime,
		ifReCoverMap:  ifReCoverMap,
		lastBuyValue:  make(map[string]float64),
		lastSellValue: make(map[string]float64),
	}
}

func NewST0StrategyFromConfig(dir string, BTConfile string, sec string, StgConfile string) ST0Strategy {
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
	c = configer.New(dir + StgConfile)
	err = c.Load()
	if err != nil {
		panic(err)
	}
	err = c.Unmarshal()
	if err != nil {
		panic(err)
	}
	tmp := c.GetString((sec + ".IfBT"))
	// declear a bool variable
	var IfBT bool
	// 如果lower后的值为true或t，则为true

	if strings.ToLower(tmp) == "true" || strings.ToLower(tmp) == "t" {
		IfBT = true
	} else {
		IfBT = false
	}
	return NewST0Strategy(IfBT, sinstrnames, sindinames, c.GetIntSlice(sec+".HoldNums"), c.GetInt(sec+".Tlimit"), c.GetString(sec+".InitBuyTime"), c.GetString(sec+".StartTime"), c.GetString(sec+".StopTime"))
}
func (t0 *ST0Strategy) CheckEligible(o *order.StockOrder, SA *stockaccount.StockAccount) bool {
	switch o.OrderDirection {
	case "Buy":
		if o.CalEquity() <= SA.Fundavail {
			return true
		}
	case "Sell":
		// check o.InstID is in SA.PosMap
		if _, ok := SA.PosMap[o.InstID]; ok {
			// check the previous position is enough
			if o.OrderNum <= SA.PosMap[o.InstID].CalPosPrevNum() {
				return true
			}
		}
	}
	return false
}
func (t0 *ST0Strategy) ActOnData(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap, Eval func([]float64) []float64) (orderRes OrderResult) {
	// 1. check if a new day, then change all ifReCoverMap to false
	// # 1.1 使用字符串形式比较日期，一旦VDS更改格式需要调整
	if t0.lastDate != datetime[:10] {
		for k := range t0.ifReCoverMap {
			t0.ifReCoverMap[k] = false
		}
	}

	// 2. 获取当前时间并判定是否介于StartTime与StopTime，是则进行常规操作
	// 否 且大于StopTime,则查看恢复操作状态并进行恢复持仓操作。
	for instID, SBDE := range bc.Stockdata {
		if !ContainNaN(SBDE.IndiDataMap) {
			// ~ 计算策略依托的有效数值
			// // 常规手动操作
			// % GEP 引入
			var GEPSlice = make([]float64, len(t0.IndiNames))
			for i := 0; i < len(t0.IndiNames); i++ {
				GEPSlice[i] = SBDE.IndiDataMap[t0.IndiNames[i]]
			}
			tmps := Eval(GEPSlice)
			// % GEP 引入完毕

			// % GEP 应用 注意这是genomeset的应用
			buyval := tmps[0] //SBDE.IndiDataMap["Close"] - SBDE.IndiDataMap["Amount"]/SBDE.IndiDataMap["Volume"]
			lstbv, bok := t0.lastBuyValue[instID]
			// % GEP 应用 注意这是genomeset的应用
			sellval := tmps[1] //SBDE.IndiDataMap["Close"] - SBDE.IndiDataMap["Open"]
			// % GEP 应用完毕
			lstsv, sok := t0.lastSellValue[instID]
			// # 2.1 在ifBT状态下，需要特判第一天的买入操作
			if t0.IfBT {
				// / 因为T+1,所以不用再担心第一天股票交易时间段下单问题。
				if !t0.initBTstate[instID] {
					if datetime[11:] == t0.InitBuyTime {
						if val, ok := SBDE.IndiDataMap["Close"]; ok {
							// 买入到持仓
							tmpSCP := cp.SimpleNewSCPFromMap(CPMap, instID)
							new_order := order.NewStockOrder(instID, false, false, datetime, val, float64(t0.holdNumMap[instID]), "Buy", &tmpSCP)
							if t0.CheckEligible(&new_order, &vAcct.SAcct) {
								orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
								t0.initBTstate[instID] = true
							}
						}
					}
				}
			}
			if datetime[11:] >= t0.StartTime && datetime[11:] <= t0.StopTime {
				tmpSCP := cp.SimpleNewSCPFromMap(CPMap, instID)
				TradeN := t0.holdNumMap[instID] / t0.Tlimit //这里选择了第一支股票SHoldNum[0] TradeN为每次交易数量

				if bok && sok {
					if val, ok := SBDE.IndiDataMap["Close"]; ok {
						// ~ 建立tmpOrderSlice
						tmpOrderSlice := make([]order.StockOrder, 0)
						// ~ 策略买入逻辑
						if buyval > 0 && lstbv < 0 && t0.tState[instID] <= 0 && t0.tCounter[instID] <= t0.Tlimit*2 {
							// 买入TradeN
							new_order := order.NewStockOrder(instID, false, false, datetime, val, float64(TradeN), "Buy", &tmpSCP)
							if t0.CheckEligible(&new_order, &vAcct.SAcct) {
								tmpOrderSlice = append(tmpOrderSlice, new_order)
								// tCounter计数器+1
								t0.tCounter[instID] += 1
								// tState 赋值
								t0.tState[instID] = 1
							}
						}
						// ~ 策略卖出逻辑
						if sellval > 0 && lstsv < 0 && t0.tState[instID] >= 0 && t0.tCounter[instID] <= t0.Tlimit*2 {
							// 卖出TradeN
							new_order := order.NewStockOrder(instID, false, false, datetime, val, float64(TradeN), "Sell", &tmpSCP)
							if t0.CheckEligible(&new_order, &vAcct.SAcct) {
								tmpOrderSlice = append(tmpOrderSlice, new_order)
								// tCounter计数器+1
								t0.tCounter[instID] += 1
								// tState 赋值
								t0.tState[instID] = -1
							}
						}
						// ~ net tmpOrderSlice
						if len(tmpOrderSlice) > 0 {
							nettedOrders := NetSOrders(tmpOrderSlice)
							orderRes.StockOrderS = append(orderRes.StockOrderS, nettedOrders...)
						}
					}
				}
			} else if t0.initBTstate[instID] && datetime[11:] > t0.StopTime {
				// / 恢复持仓操作
				tmpSCP := cp.SimpleNewSCPFromMap(CPMap, instID)
				if !t0.ifReCoverMap[instID] {
					if val, ok := SBDE.IndiDataMap["Close"]; ok {
						// 恢复持仓
						netNum := t0.holdNumMap[instID] - int(vAcct.SAcct.PosMap[instID].CalPosPrevNum()+vAcct.SAcct.PosMap[instID].CalPosTdyNum())
						if netNum > 0 {
							new_order := order.NewStockOrder(instID, false, false, datetime, val, float64(netNum), "Buy", &tmpSCP)
							if t0.CheckEligible(&new_order, &vAcct.SAcct) {
								orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
							}
						} else if netNum < 0 {
							new_order := order.NewStockOrder(instID, false, false, datetime, val, -float64(netNum), "Sell", &tmpSCP)
							if t0.CheckEligible(&new_order, &vAcct.SAcct) {
								orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
							}
						}
						t0.ifReCoverMap[instID] = true
					}
				}
			}
			// ~lastBuyValue 赋值 lastSellValue 赋值
			t0.lastBuyValue[instID] = buyval
			t0.lastSellValue[instID] = sellval
		}
	}
	t0.lastDate = datetime[:10]
	return
}

func (t0 *ST0Strategy) ActOnDataMAN(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap) (orderRes OrderResult) {
	// 1. check if a new day, then change all ifReCoverMap to false
	// # 1.1 使用字符串形式比较日期，一旦VDS更改格式需要调整
	if t0.lastDate != datetime[:10] {
		for k := range t0.ifReCoverMap {
			t0.ifReCoverMap[k] = false
		}
	}
	// 2. 获取当前时间并判定是否介于StartTime与StopTime，是则进行常规操作
	// 否 且大于StopTime,则查看恢复操作状态并进行恢复持仓操作。
	for instID, SBDE := range bc.Stockdata {
		if !ContainNaN(SBDE.IndiDataMap) {
			// ~ 计算策略依托的有效数值
			buyval := SBDE.IndiDataMap["Close"] - SBDE.IndiDataMap["Amount"]/SBDE.IndiDataMap["Volume"]
			lstbv, bok := t0.lastBuyValue[instID]
			sellval := SBDE.IndiDataMap["Close"] - SBDE.IndiDataMap["MA5"]
			lstsv, sok := t0.lastSellValue[instID]
			// # 2.1 在ifBT状态下，需要特判第一天的买入操作
			if t0.IfBT {
				// / 因为T+1,所以不用再担心第一天股票交易时间段下单问题。
				if !t0.initBTstate[instID] {
					if datetime[11:] == t0.InitBuyTime {
						if val, ok := SBDE.IndiDataMap["Close"]; ok {
							// 买入到持仓
							tmpSCP := cp.SimpleNewSCPFromMap(CPMap, instID)
							new_order := order.NewStockOrder(instID, false, false, datetime, val, float64(t0.holdNumMap[instID]), "Buy", &tmpSCP)
							if t0.CheckEligible(&new_order, &vAcct.SAcct) {
								orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
								t0.initBTstate[instID] = true
							}
						}
					}
				}
			}
			if t0.initBTstate[instID] && datetime[11:] > t0.StartTime && datetime[11:] < t0.StopTime {
				// / 常规手动操作
				tmpSCP := cp.SimpleNewSCPFromMap(CPMap, instID)
				TradeN := t0.holdNumMap[instID] / t0.Tlimit //这里选择了第一支股票SHoldNum[0] TradeN为每次交易数量

				if bok && sok {
					if val, ok := SBDE.IndiDataMap["Close"]; ok {
						// ~ 建立tmpOrderSlice
						tmpOrderSlice := make([]order.StockOrder, 0)
						// ~ 策略买入逻辑
						if buyval > 0 && lstbv < 0 && t0.tState[instID] <= 0 && t0.tCounter[instID] <= t0.Tlimit*2 {
							// 买入TradeN
							new_order := order.NewStockOrder(instID, false, false, datetime, val, float64(TradeN), "Buy", &tmpSCP)
							if t0.CheckEligible(&new_order, &vAcct.SAcct) {
								tmpOrderSlice = append(tmpOrderSlice, new_order)
								// tCounter计数器+1
								t0.tCounter[instID] += 1
								// tState 赋值
								t0.tState[instID] = 1
							}
						}
						// ~ 策略卖出逻辑
						if sellval > 0 && lstsv < 0 && t0.tState[instID] >= 0 && t0.tCounter[instID] <= t0.Tlimit*2 {
							// 卖出TradeN
							new_order := order.NewStockOrder(instID, false, false, datetime, val, float64(TradeN), "Sell", &tmpSCP)
							if t0.CheckEligible(&new_order, &vAcct.SAcct) {
								tmpOrderSlice = append(tmpOrderSlice, new_order)

								// tCounter计数器+1
								t0.tCounter[instID] += 1
								// tState 赋值
								t0.tState[instID] = -1
							}
						}
						// ~ net tmpOrderSlice
						if len(tmpOrderSlice) > 0 {
							nettedOrders := NetSOrders(tmpOrderSlice)
							orderRes.StockOrderS = append(orderRes.StockOrderS, nettedOrders...)
						}
					}
				}

			} else if t0.initBTstate[instID] && datetime[11:] > t0.StopTime {
				// / 恢复持仓操作
				tmpSCP := cp.SimpleNewSCPFromMap(CPMap, instID)
				if !t0.ifReCoverMap[instID] {
					if val, ok := SBDE.IndiDataMap["Close"]; ok {
						// 恢复持仓
						netNum := t0.holdNumMap[instID] - int(vAcct.SAcct.PosMap[instID].CalPosPrevNum()+vAcct.SAcct.PosMap[instID].CalPosTdyNum())
						if netNum > 0 {
							new_order := order.NewStockOrder(instID, false, false, datetime, val, float64(netNum), "Buy", &tmpSCP)
							if t0.CheckEligible(&new_order, &vAcct.SAcct) {
								orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
							}
						} else if netNum < 0 {
							new_order := order.NewStockOrder(instID, false, false, datetime, val, -float64(netNum), "Sell", &tmpSCP)
							if t0.CheckEligible(&new_order, &vAcct.SAcct) {
								orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
							}
						}
						t0.ifReCoverMap[instID] = true
					}
				}
			}
			// ~lastBuyValue 赋值 lastSellValue 赋值
			t0.lastBuyValue[instID] = buyval
			t0.lastSellValue[instID] = sellval
		}
	}
	t0.lastDate = datetime[:10]
	return
}
