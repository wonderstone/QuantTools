package strategyModule

import (
	// "math"
	"strings"

	"github.com/wonderstone/QuantTools/account/virtualaccount"
	"github.com/wonderstone/QuantTools/configer"
	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/dataprocessor"
	"github.com/wonderstone/QuantTools/order"
)

// Stock T0Strategy

type ST0Strategy struct {
	IfBT          bool               // 是否为回测/实盘状态
	initBTstate   map[string]bool    // 是否已经初始化了回测买入状态(内部)
	InstNames     []string           // 股票标的名称
	IndiNames     []string           // 股票参与GEP指标名称，注意其数量不大于BarDE内信息数量，且strategy内可见BarDE的数据
	HoldNums      []int              // 股票标的持仓数量
	holdNumMap    map[string]int     // 股票标的持仓数量, key is the stock name,get from SInstNames and SHoldNums(内部)
	Tlimit        int                // 交易次数限制
	tCounter      map[string]int     // 交易次数计数器(内部)
	tState        map[string]int     // 交易状态(内部) 0:未买入 1:已买入 -1:已卖出
	InitBuyTime   string             // 初始化买入时间
	StartTime     string             // 交易开始时间
	StopTime      string             // 交易结束时间
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

	for i, name := range InstNames {
		initBTstate[name] = false
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
	tmp := c.GetString(("default.IfBT"))
	// declear a bool variable
	var IfBT bool
	// 如果lower后的值为true或t，则为true

	if strings.ToLower(tmp) == "true" || strings.ToLower(tmp) == "t" {
		IfBT = true
	} else {
		IfBT = false
	}
	return NewST0Strategy(IfBT, sinstrnames, sindinames, c.GetIntSlice("default.HoldNums"), c.GetInt("default.Tlimit"), c.GetString("default.InitBuyTime"), c.GetString("default.StartTime"), c.GetString("default.StopTime"))

	// return NewST0Strategy(c.GetBool(sec, "IfBT"), c.GetStringSlice(sec, "InstNames"), c.GetStringSlice(sec, "IndiNames"), c.GetIntSlice(sec, "HoldNums"), c.GetInt(sec, "Tlimit"), c.GetString(sec, "StartTime"), c.GetString(sec, "StopTime"))

}

func (t0 *ST0Strategy) ActOnData(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap, Eval func([]float64) []float64) (orderRes OrderResult) {

	return
}

func (t0 *ST0Strategy) ActOnDataMAN(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap) (orderRes OrderResult) {

	// a new day, change all ifReCoverMap to false
	if t0.lastDate != datetime[:10] {
		for k := range t0.ifReCoverMap {
			t0.ifReCoverMap[k] = false
		}
	}

	// 获取当前时间并判定是否介于StartTime与StopTime，是则进行常规操作
	// 否 且大于StopTime,则查看恢复操作状态并进行恢复持仓操作。
	for instID, SBDE := range bc.Stockdata {
		if !ContainNaN(SBDE.IndiDataMap) {
			if t0.IfBT {
				// 因为T+1,所以不用再担心第一天股票交易时间段下单问题。
				if !t0.initBTstate[instID] {
					if datetime[11:] == t0.InitBuyTime {
						// 买入到持仓
						tmpSCP := cp.SimpleNewSCPFromMap(CPMap, instID)
						orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["Close"], float64(t0.holdNumMap[instID]), "Buy", &tmpSCP))
						t0.initBTstate[instID] = true
					}
				}
			}

			if datetime[11:] >= t0.StartTime && datetime[11:] <= t0.StopTime {
				// // 常规手动操作
				tmpSCP := cp.SimpleNewSCPFromMap(CPMap, instID)
				buyval := SBDE.IndiDataMap["Close"] - SBDE.IndiDataMap["Amount"]/SBDE.IndiDataMap["Volume"]
				lstbv, bok := t0.lastBuyValue[instID]
				sellval := SBDE.IndiDataMap["Close"] - SBDE.IndiDataMap["Open"]
				lstsv, sok := t0.lastSellValue[instID]
				TradeN := t0.holdNumMap[instID] / t0.Tlimit //这里选择了第一支股票SHoldNum[0] TradeN为每次交易数量

				if bok && sok {
					if buyval > 0 && lstbv < 0 && t0.tState[instID] <= 0 && t0.tCounter[instID] <= t0.Tlimit*2 {
						// 买入TradeN
						orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["Close"], float64(TradeN), "Buy", &tmpSCP))

						// tCounter计数器+1
						t0.tCounter[instID] += 1

						// lastBuyValue 赋值 lastSellValue 赋值
						t0.lastBuyValue[instID] = buyval
						t0.lastSellValue[instID] = sellval

						// tState 赋值
						t0.tState[instID] = 1
					}

					if sellval > 0 && lstsv < 0 && t0.tState[instID] >= 0 && t0.tCounter[instID] <= t0.Tlimit*2 {
						// 卖出TradeN
						orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["Close"], float64(TradeN), "Sell", &tmpSCP))

						// tCounter计数器+1
						t0.tCounter[instID] += 1

						// lastBuyValue 赋值 lastSellValue 赋值
						t0.lastBuyValue[instID] = buyval
						t0.lastSellValue[instID] = sellval

						// tState 赋值
						t0.tState[instID] = -1
					}
				}

			} else if datetime[11:] > t0.StopTime {

				tmpSCP := cp.SimpleNewSCPFromMap(CPMap, instID)
				if !t0.ifReCoverMap[instID] {
					// 恢复持仓
					netNum := t0.holdNumMap[instID] - int(vAcct.SAcct.PosMap[instID].CalPosPrevNum()+vAcct.SAcct.PosMap[instID].CalPosTdyNum())
					if netNum > 0 {
						orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["Close"], vAcct.SAcct.PosMap[instID].CalPosPrevNum(), "Buy", &tmpSCP))
					} else if netNum < 0 {
						orderRes.StockOrderS = append(orderRes.StockOrderS, order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["Close"], vAcct.SAcct.PosMap[instID].CalPosPrevNum(), "Sell", &tmpSCP))
					}
					t0.ifReCoverMap[instID] = true
				}

			}

		}
	}

	t0.lastDate = datetime[:10]
	return
}
