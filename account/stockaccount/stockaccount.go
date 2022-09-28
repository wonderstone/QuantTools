/* 股票账户模拟，仅保留最基本核心字段和方法 */
package stockaccount

import (
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/wonderstone/QuantTools/account"
	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/order"
)

type StockAccount struct {
	//共8个字段  其中  数值类型6个  真实账户核心字段仅保留MktVal、FundAvail
	InitTime         string
	UdTime           string  // update everytime
	MktVal           float64 // 权益 最终采用MktVal
	Fundavail        float64
	AllProfit        float64                   // 根据计算更新， 保留可以脱离数据进行查询上一状态的Profit
	AllCommission    float64                   // 单独标记
	PosMap           map[string]*PositionSlice //用指针版本
	MarketValueSlice []account.MktValDataType
	// add one tmp UUID field for log testing
	UUID string
}

// 生成一个自定义的stockAccount
func NewStockAccount(initTime string, cash float64) StockAccount {
	return StockAccount{
		// 此处初始化字段4个，创建时间、更新时间、总市值和可用资金
		InitTime:         initTime,
		UdTime:           initTime,
		MktVal:           cash,
		Fundavail:        cash,
		UUID:             uuid.New().String(),
		PosMap:           make(map[string]*PositionSlice),
		MarketValueSlice: []account.MktValDataType{{Time: initTime, MktVal: cash}},
	}
}

// load from a yaml file to generate a stock account with specific fields
func NewSAFromConfig(filename string, configpath string, sec string, cpm cp.CPMap) StockAccount {
	// 获取配置文件记录的 VA.SACCT map
	viper.SetConfigName(filename)
	viper.AddConfigPath(configpath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	SAcctMap := viper.GetStringMap(sec)
	// get the posmap, key is code, value is the position slice:
	posmap := make(map[string]*PositionSlice)
	// process interface{} of SAcctMap["posmap"] to yield a map with instid as key and position slice as value
	tmpPosMap := SAcctMap["posmap"].(map[string]interface{})

	for key, value := range tmpPosMap {
		// key is the instid, value should be the position slice
		scp := cp.SimpleNewSCPFromMap(cpm, key)
		tmpPSlice := NewPosSlice()
		// position slice is another map with 3 elements as its keys and interface{} as its values to copy with different types
		tmpPSliceMap := value.(map[string]interface{})

		// one key has one scp
		tmpPSlice.UdTime = tmpPSliceMap["udtime"].(string)
		// copy the value from postdys slice and iter the tmpPSliceMap["postdys"]
		tmpPosTdys := make([]PositionDetail, 0)
		for _, PDI := range tmpPSliceMap["postdys"].([]interface{}) {
			// PDIMap with pd element names as keys and interface{} as values
			PDIMap := PDI.(map[string]interface{})
			// make pd and append to tmpPosTdys
			tmpPD := PositionDetail{
				UdTime:    PDIMap["udtime"].(string),
				InstID:    PDIMap["instid"].(string),
				BasePrice: PDIMap["baseprice"].(float64),
				LastPrice: PDIMap["lastprice"].(float64),
				Num:       float64(PDIMap["num"].(int)),
				Equity:    PDIMap["equity"].(float64),
				SCP:       &scp,
			}
			tmpPosTdys = append(tmpPosTdys, tmpPD)
		}
		tmpPSlice.PosTdys = tmpPosTdys
		// copy the value from  posprevs slice
		tmpPosPrev := make([]PositionDetail, 0)
		for _, PDI := range tmpPSliceMap["posprevs"].([]interface{}) {
			// PDIMap with pd element names as keys and interface{} as values
			PDIMap := PDI.(map[string]interface{})
			// make pd and append to tmpPosTdys
			tmpPD := PositionDetail{
				UdTime:    PDIMap["udtime"].(string),
				InstID:    PDIMap["instid"].(string),
				BasePrice: PDIMap["baseprice"].(float64),
				LastPrice: PDIMap["lastprice"].(float64),
				Num:       float64(PDIMap["num"].(int)),
				Equity:    PDIMap["equity"].(float64),
				SCP:       &scp,
			}
			tmpPosPrev = append(tmpPosPrev, tmpPD)
		}
		tmpPSlice.PosPrevs = tmpPosPrev
		posmap[key] = tmpPSlice
	}

	return StockAccount{
		// 此处初始化字段4个，创建时间、更新时间、总市值和可用资金
		// InitTime:      SAcctMap["inittime"].(string),
		InitTime:      viper.GetString("VA.SACCT.inittime"),
		UdTime:        viper.GetString("VA.SACCT.udtime"),
		MktVal:        viper.GetFloat64("VA.SACCT.mktval"),
		Fundavail:     viper.GetFloat64("VA.SACCT.fundavail"),
		AllProfit:     viper.GetFloat64("VA.SACCT.allprofit"),
		AllCommission: viper.GetFloat64("VA.SACCT.allcommission"),
		UUID:          viper.GetString("VA.SACCT.uuid"),
		PosMap:        posmap,
		// MarketValueSlice: tmpMap["marketvalueslice"].([]account.MktValDataType),
	}

}

// reset the MarketValueSlice to nil
func (SA *StockAccount) ResetMVSlice() {
	SA.MarketValueSlice = nil
}

// 汇总所有Equity
func (SA *StockAccount) Equity() (Equity float64) {
	for _, ps := range SA.PosMap {
		Equity += ps.CalEquity()
	}
	return
}

// Eligible check for order
func (SA *StockAccount) CheckEligible(o *order.StockOrder) {
	switch o.OrderDirection {
	case "Buy":
		if o.CalEquity() <= SA.Fundavail {
			o.IsEligible = true
		}
	case "Sell":
		// check o.InstID is in SA.PosMap
		if _, ok := SA.PosMap[o.InstID]; ok {
			// check the previous position is enough
			if o.OrderNum <= SA.PosMap[o.InstID].CalPosPrevNum() {
				o.IsEligible = true
			}
		}
	}
}

// 针对order产生反应
func (SA *StockAccount) ActOnOrder(SO *order.StockOrder) {
	if SO.IsExecuted {
		if SA.Fundavail <= SO.CalEquity() {
			// panic("确保账户足够资金")
		} else {
			// in principle, backtest should be done under one mutex lock
			// insurance: add a mutex for stock account write
			// SA.Lock()
			// defer 后进先出
			// defer SA.Unlock()
			// 1. 不修正初始化时间字段
			// 2. 修正刷新时间
			SA.UdTime = SO.OrderTime
			// 7. 调整PosMap内的对应PositionSlice
			RealizedProfit, Comm, Equity := 0.0, 0.0, 0.0
			if _, ok := SA.PosMap[SO.InstID]; ok {
				RealizedProfit, Comm, Equity = SA.PosMap[SO.InstID].UpdateWithOrder(SO)
			} else {
				SA.PosMap[SO.InstID] = NewPosSlice() //&PositionSlice{} //{UdTime: FO.OrderTime}
				RealizedProfit, Comm, Equity = SA.PosMap[SO.InstID].UpdateWithOrder(SO)
			}
			// 3.修正 AllProfit
			SA.AllProfit += RealizedProfit
			// 4.修正 AllCommission
			SA.AllCommission += Comm
			// 由价格变动确认profit确定新MV  由价格变动确认Margin，进而确定FundAvail
			// 5.修正 Fundavail
			SA.Fundavail += RealizedProfit - Comm - Equity
			// 6.修正 MktVal
			SA.MktVal = SA.Fundavail + SA.Equity()
			// 8.不修正字段 MarketValueSlice
		}
	}
}

// 针对数据反应
func (SA *StockAccount) ActOnUpdateMI(UpdateTimeStamp string, InstID string, Value float64) {
	// in principle, backtest should be done under one mutex lock
	// insurance: add a mutex for stock account write
	// SA.Lock()
	// defer 后进先出
	// defer SA.Unlock()
	// 1. 不修正初始化时间字段
	// 2. 修正刷新时间
	SA.UdTime = UpdateTimeStamp
	// 3.不修正 AllProfit
	// 4.不修正 AllCommission

	// 6.不修正 Fundavail
	// 7. 调整PosMap内的对应PositionSlice
	if _, ok := SA.PosMap[InstID]; ok {
		// 更新pd内lastprice数值
		SA.PosMap[InstID].UpdateWithUMI(UpdateTimeStamp, Value)
	}
	// 5.修正 MktVal
	SA.MktVal = SA.Fundavail + SA.Equity()
	// 8.不修正字段 MarketValueSlice
}

// CloseMarket 在账户层面实施
func (SA *StockAccount) ActOnCM() {
	// in principle, backtest should be done under one mutex lock
	// insurance: add a mutex for stock account write
	// SA.Lock()
	// defer 后进先出
	// defer SA.Unlock()
	// 1. 不修正初始化时间字段
	// 2. 不修正刷新时间
	// 3.不修正 AllProfit
	// 4.不修正 AllCommission
	// 5.不修正 MktVal
	// 6.不修正 Fundavail
	// 7. 调整PosMap内的对应PositionSlice
	for key := range SA.PosMap {
		SA.PosMap[key].UpdateWithCM()
	}
	// 8.修正字段 MarketValueSlice
	SA.MarketValueSlice = append(SA.MarketValueSlice, account.MktValDataType{Time: SA.UdTime, MktVal: SA.MktVal})

}
