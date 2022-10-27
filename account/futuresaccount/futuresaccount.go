/* 期货账户模拟，仅保留最基本核心字段和方法 */
package futuresaccount

import (
	"github.com/google/uuid"
	"github.com/wonderstone/QuantTools/account"
	"github.com/wonderstone/QuantTools/configer"
	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/order"
)

type FuturesAccount struct {
	// * 共9个字段  其中  数值类型7个  真实账户核心字段仅保留MktVal、FundAvail
	InitTime         string
	UdTime           string // update everytime
	BmkVal           float64
	MktVal           float64 // 权益 最终采用MktVal
	Fundavail        float64
	AllProfit        float64                   // 根据计算更新， 保留可以脱离数据进行查询上一状态的Profit
	AllCommission    float64                   // 单独标记
	PosMap           map[string]*PositionSlice //用指针版本
	MarketValueSlice []account.MktValDataType
	// todo: add one tmp UUID field for log testing
	UUID string
}

// 生成一个自定义的FuturesAccount
func NewFuturesAccount(initTime string, cash float64) FuturesAccount {
	return FuturesAccount{
		InitTime:         initTime,
		UdTime:           initTime,
		BmkVal:           cash,
		MktVal:           cash,
		Fundavail:        cash,
		UUID:             uuid.New().String(),
		PosMap:           make(map[string]*PositionSlice),
		MarketValueSlice: []account.MktValDataType{{Time: initTime, MktVal: cash}},
	}
}

// load from a yaml file to generate a FuturesAccount with specific fields
func NewFAFromConfig(filename string, configpath string, sec string, cpm cp.CPMap) FuturesAccount {
	// 获取配置文件记录的 Va.FACCT map
	c := configer.New(configpath + filename)
	err := c.Load()
	if err != nil {
		panic(err)
	}
	err = c.Unmarshal()
	if err != nil {
		panic(err)
	}
	FAcctMap := c.GetStringMap(sec)
	// get the posmap, key is code, value is the position slice:
	posmap := make(map[string]*PositionSlice)
	// process interface{} of SAcctMap["posmap"] to yield a map with instid as key and position slice as value
	tmpPosMap := FAcctMap["posmap"].(map[string]interface{})

	for key, value := range tmpPosMap {
		// key is the instid, value should be the position slice
		fcp := cp.SimpleNewFCPFromMap(cpm, key)
		tmpPSlice := NewPosSlice()
		// position slice is another map with 3 elements as its keys and interface{} as its values to copy with different types
		tmpPSliceMap := value.(map[string]interface{})

		// one key has one scp
		tmpPSlice.UdTime = tmpPSliceMap["udtime"].(string)
		// copy the value from postdys slice and iter the tmpPSliceMap["postdys"]
		tmpPosTdys := make([]PositionDetail, len(tmpPSliceMap["postdys"].([]interface{})))
		for _, PDI := range tmpPSliceMap["postdys"].([]interface{}) {
			// PDIMap with pd element names as keys and interface{} as values
			PDIMap := PDI.(map[string]interface{})
			// make pd and append to tmpPosTdys
			// 有一个dir的问题！！！！！！！！

			//
			tmpPD := PositionDetail{
				UdTime:    PDIMap["udtime"].(string),
				InstID:    PDIMap["instid"].(string),
				BasePrice: PDIMap["baseprice"].(float64),
				LastPrice: PDIMap["lastprice"].(float64),
				Num:       float64(PDIMap["num"].(int)),
				Margin:    PDIMap["margin"].(float64),
				FCP:       &fcp,
			}
			tmpPosTdys = append(tmpPosTdys, tmpPD)
		}
		tmpPSlice.PosTdys = tmpPosTdys
		// copy the value from  posprevs slice
		tmpPosPrev := make([]PositionDetail, len(tmpPSliceMap["posprevs"].([]interface{})))
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
				Margin:    PDIMap["margin"].(float64),
				FCP:       &fcp,
			}
			tmpPosPrev = append(tmpPosPrev, tmpPD)
		}
		tmpPSlice.PosPrevs = tmpPosPrev
		posmap[key] = tmpPSlice

	}

	return FuturesAccount{
		// 此处初始化字段4个，创建时间、更新时间、总市值和可用资金
		InitTime:      FAcctMap["inittime"].(string),
		UdTime:        FAcctMap["udtime"].(string),
		BmkVal:        account.GetFloat64(FAcctMap["bmkval"]),
		MktVal:        account.GetFloat64(FAcctMap["mktval"]),
		Fundavail:     account.GetFloat64(FAcctMap["fundavail"]),
		AllProfit:     account.GetFloat64(FAcctMap["allprofit"]),
		AllCommission: account.GetFloat64(FAcctMap["allcommission"]),
		UUID:          FAcctMap["uuid"].(string),

		PosMap: posmap,
		// MarketValueSlice: tmpMap["marketvalueslice"].([]account.MktValDataType),

	}

}

// Eligible check for order
func (FA *FuturesAccount) CheckEligible(o *order.FuturesOrder) {
	switch o.OrderType {
	case "Open":
		if o.CalMargin() < FA.Fundavail {
			o.IsEligible = true
		}

	case "CloseToday":
		// check if o.InstID is in FA.PosMap
		if _, ok := FA.PosMap[o.InstID]; ok {
			numl, nums := FA.PosMap[o.InstID].CalPosTdyNum()
			if o.OrderDirection == "Buy" && numl >= o.OrderNum {
				o.IsEligible = true
			}
			if o.OrderDirection == "Sell" && nums >= o.OrderNum {
				o.IsEligible = true
			}
		}
	case "ClosePrevious":
		// check if o.InstID is in FA.PosMap
		if _, ok := FA.PosMap[o.InstID]; ok {
			numl, nums := FA.PosMap[o.InstID].CalPosPrevNum()
			if o.OrderDirection == "Buy" && numl >= o.OrderNum {
				o.IsEligible = true
			}
			if o.OrderDirection == "Sell" && nums >= o.OrderNum {
				o.IsEligible = true
			}
		}
	}
}

// reset the MarketValueSlice to nil
func (FA *FuturesAccount) ResetMVSlice() {
	FA.MarketValueSlice = nil
}

// 汇总所有Margin
func (FA *FuturesAccount) Margin() (Margin float64) {
	for _, ps := range FA.PosMap {
		Margin += ps.CalMargin()
	}
	return
}

// 针对order产生反应
// 注意针对保证金问题 simnow的无限易实时资金页面保证金计算有问题，而真实账户是用真实开仓价处理。
// 账户结构采用  账户的保证金占用(实时更新) + 可用资金(fundavail) = 动态权益市值 MarketValue (实时刷新)  而是通过结算刷新
// 浮动盈亏 = (最新价 - 开仓均价)*手数  相当于基准从开始计算
// 持仓盈亏 = (最新价 - 持仓均价)*手数  相当于基准从昨天计算  持仓均价是用每日MTM的结算价进行替换的
func (FA *FuturesAccount) ActOnOrder(FO *order.FuturesOrder) {
	if FO.IsExecuted {
		// 相比股票，期货只检查了开仓保证金，没有检查平仓合约数
		// 这部分逻辑放在了Account的CheckEligible中，就不再重复了
		if FA.Fundavail <= FO.CalMargin() && FO.OrderType == "Open" {
			// panic("确保账户具有足够资金")
		} else {
			// in principle, backtest should be done under one mutex lock
			// insurance: add a mutex for stock account write
			// FA.Lock()
			// defer 后进先出
			// defer FA.Unlock()
			// 1. 初始化时间字段不修正
			// 2. 修正刷新时间
			FA.UdTime = FO.OrderTime
			// 7. 调整PosMap内的对应PositionSlice
			RealizedProfit, Comm, UnRealizedProfit := 0.0, 0.0, 0.0
			if _, ok := FA.PosMap[FO.InstID]; ok {
				RealizedProfit, Comm, UnRealizedProfit = FA.PosMap[FO.InstID].UpdateWithOrder(FO)
			} else {
				FA.PosMap[FO.InstID] = NewPosSlice() //&PositionSlice{} //{UdTime: FO.OrderTime}
				RealizedProfit, Comm, UnRealizedProfit = FA.PosMap[FO.InstID].UpdateWithOrder(FO)
			}
			// 3.修正 AllProfit
			FA.AllProfit += RealizedProfit
			// 4.修正 AllCommission
			FA.AllCommission += Comm
			// 由价格变动确认profit确定新MV  由价格变动确认Margin，进而确定FundAvail
			// 5.修正 MktVal
			FA.BmkVal += RealizedProfit - Comm
			FA.MktVal = FA.BmkVal + UnRealizedProfit
			// 6.修正 Fundavail
			FA.Fundavail = FA.BmkVal - FA.Margin()
			// 8.不修正字段 MarketValueSlice
		}
	}
}

// 针对数据反应
func (FA *FuturesAccount) ActOnUpdateMI(UpdateTimeStamp string, InstID string, Value float64) {
	// in principle, backtest should be done under one mutex lock
	// insurance: add a mutex for stock account write
	// FA.Lock()
	// defer 后进先出
	// defer FA.Unlock()
	// 1. 初始化时间字段不修正
	// 2. 修正刷新时间
	FA.UdTime = UpdateTimeStamp
	// 7. 调整PosMap内的对应PositionSlice
	UnRealizedProfit := 0.0
	if _, ok := FA.PosMap[InstID]; ok {
		// 更新pd内lastprice数值
		FA.PosMap[InstID].UpdateWithUMI(UpdateTimeStamp, Value)
	}
	for _, ps := range FA.PosMap {
		UnRealizedProfit += ps.CalUnRealizedProfit()
	}
	// 3.不修正 AllProfit
	// 4.不修正 AllCommission
	// 5.修正 MktVal
	FA.MktVal = FA.BmkVal + UnRealizedProfit
	// 6.修正 Fundavail
	FA.Fundavail = FA.BmkVal - FA.Margin()
	// 8.不修正字段 MarketValueSlice
}

// MTM 类似针对数据反应   但修改baseprice和添加MV进入slice
func (FA *FuturesAccount) ActOnMTM(UpdateTimeStamp string, InstID string, Value float64) {
	// in principle, backtest should be done under one mutex lock
	// insurance: add a mutex for stock account write
	// FA.Lock()
	// defer 后进先出
	// defer FA.Unlock()
	// 1. 不修正初始化时间字段
	// 2. 修正刷新时间
	FA.UdTime = UpdateTimeStamp
	// 7. 调整PosMap内的对应PositionSlice
	RealizedProfit := 0.0
	if _, ok := FA.PosMap[InstID]; ok {
		RealizedProfit = FA.PosMap[InstID].UpdateWithMTM(UpdateTimeStamp, Value)
	}
	// 3.修正 AllProfit
	FA.AllProfit += RealizedProfit
	// 4.不修正 AllCommission
	// 由价格变动确认profit确定新MV  由价格变动确认Margin，进而确定FundAvail
	// 5.修正 MktVal
	FA.BmkVal += RealizedProfit
	FA.MktVal = FA.BmkVal
	// 6.修正 Fundavail
	FA.Fundavail = FA.BmkVal - FA.Margin()
	// 8.修正字段 MarketValueSlice
	FA.MarketValueSlice = append(FA.MarketValueSlice, account.MktValDataType{Time: UpdateTimeStamp, MktVal: FA.MktVal})

}
