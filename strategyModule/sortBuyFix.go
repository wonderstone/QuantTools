package strategyModule

// market timing strategy
// one expression for one strategy, therefore one market for one strategy
import (
	//"fmt"

	"math"
	"sort"

	"github.com/elliotchance/orderedmap/v2"
	// "github.com/rs/zerolog"

	"github.com/wonderstone/QuantTools/account/stockaccount"
	"github.com/wonderstone/QuantTools/account/virtualaccount"
	"github.com/wonderstone/QuantTools/configer"
	"github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/dataprocessor"
	"github.com/wonderstone/QuantTools/order"
)

type Pair struct {
	Key   string
	Value float64
	Price float64
}
type PairList []Pair
type SortBuyStrategy struct {
	init_money  float64  //初始资金
	SInstNames  []string // 股票标的名称
	SIndiNames  []string // 股票参与GEP指标名称，注意其数量不大于BarDE内信息数量，且strategy内可见BarDE的数据
	STimeCritic string   // 时间关键字，用于判断是否需要进行交易

	stimeCondition bool    //排序买入的时间条件是否已经满足
	numHolding     int     //持仓标的数量
	max_min_ratio  float64 //最大最小比例
	position_ratio float64 //使用账户总资金比例

	holdingsM     *orderedmap.OrderedMap[string, float64]
	targetsM      *orderedmap.OrderedMap[string, float64]
	holdingRatios []float64
}

func NewSortBuyStrategy(init_money float64, SInstNms []string, SIndiNms []string, STimeCritic string,
	numHolding int, max_min_ratio float64, position_ratio float64) SortBuyStrategy {
	return SortBuyStrategy{
		init_money:  init_money,
		SInstNames:  SInstNms,
		SIndiNames:  SIndiNms,
		STimeCritic: STimeCritic,

		stimeCondition: false,
		numHolding:     numHolding,
		max_min_ratio:  max_min_ratio,
		position_ratio: position_ratio,

		holdingsM:     orderedmap.NewOrderedMap[string, float64](),
		targetsM:      orderedmap.NewOrderedMap[string, float64](),
		holdingRatios: arithmeticSequence(numHolding, max_min_ratio),
	}
}

// this function is nessary for the framework
func NewSortBuyStrategyFromConfig(dir string, BTConfile string, sec string, StgConfile string) SortBuyStrategy {
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
	init_money := c.GetFloat64("default.stockinitvalue")
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
	STimeCritic := c.GetString("default.stimecritic")
	numHolding := c.GetInt("default.numHolding")
	ratio := c.GetFloat64("default.max_min_ratio")
	cash_used_ratio := c.GetFloat64("default.position_ratio")
	return NewSortBuyStrategy(init_money, sinstrnames, sindinames, STimeCritic, numHolding, ratio, cash_used_ratio)
}

func (sb *SortBuyStrategy) CheckEligible(o *order.StockOrder, SA *stockaccount.StockAccount) bool {
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

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (sb *SortBuyStrategy) ActOnData(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap contractproperty.CPMap, Eval func([]float64) []float64) (orderRes OrderResult) {
	UpdateAcct(bc, vAcct)
	// 2. iter vAcct.SAcct.PosMap and maintain holdingsM
	for k, v := range vAcct.SAcct.PosMap {
		// 日线级策略 问题不大
		sb.holdingsM.Set(k, v.CalPosTdyNum()+v.CalPosPrevNum())
		if val, _ := sb.holdingsM.Get(k); val == 0 {
			sb.holdingsM.Delete(k)
		}
	}
	// 3. check the datetime is the executable time
	if GetTimeValue(datetime) == sb.STimeCritic {
		sb.stimeCondition = true
	} else {
		sb.stimeCondition = false
	}
	// 4. 判断股票标的切片SInstrNames是否为空 并且 时间准则为真，如果为空，则不操作股票数据循环
	if len(sb.SInstNames) != 0 && sb.stimeCondition {
		// check if the length of the stockdata is bigger than 1
		if len(bc.Stockdata) < 1 {
			// 4.1 如果为空，则panic
			panic("data error! the length of the stockdata is smaller than 1")
		} else if len(bc.Stockdata) == 1 {
			// 4.2 如果长度为1，大概率为回测初始数据量原因，亦可能为实盘停牌或单纯交易量低，
			for instID, SBDE := range bc.Stockdata {
				tmpSCP := contractproperty.SimpleNewSCPFromMap(CPMap, instID)
				// 期望买入的总市值
				targetValue := vAcct.SAcct.MktVal * sb.position_ratio
				// 期望买入的股票数量，如果符合大于开单条件(此处不为0)同步到targetsM
				tmpTM := math.Floor(((targetValue / SBDE.IndiDataMap["Close"]) / contractproperty.SimpleNewSCPFromMap(CPMap, instID).ContractSize))
				if tmpTM > 0 {
					sb.targetsM.Set(instID, tmpTM)
				}
				// 查看是否账户持有该标的
				tmpHM, Hok := sb.holdingsM.Get(instID)
				if !Hok {
					// 如果不持有该标的，且sb.holdingsM为空，则直接买入
					// （其他情况大概率为实盘数据异常，不操作）
					if sb.holdingsM.Len() == 0 {
						new_order := order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["Close"], tmpTM, "Buy", &tmpSCP)
						orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
					}
				} else {
					// 如果持有该标的，且仅持有该标的，则调整到期望数量
					// （其他情况大概率为实盘数据异常，不操作）
					if sb.holdingsM.Len() == 1 {
						if tmpHM < tmpTM {
							new_order := order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["Close"], math.Abs(tmpTM-tmpHM), "Buy", &tmpSCP)
							orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
						} else if tmpHM > tmpTM {
							new_order := order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["Close"], math.Abs(tmpTM-tmpHM), "Sell", &tmpSCP)
							orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
						}
					}
				}
			}
		} else {
			// 4.3 如果长度大于1，正常循环
			// 得到数据个数
			n := len(bc.Stockdata)
			// 修改holdingRatios
			if n < sb.numHolding {
				sb.holdingRatios = arithmeticSequence(n, sb.max_min_ratio)
			} else {
				sb.holdingRatios = arithmeticSequence(sb.numHolding, sb.max_min_ratio)
			}
			// 循环标的产生PL 标的排序并得到对应标的的持仓比例
			PL := PairList{}

			// * GEP 引入

			for instID, SBDE := range bc.Stockdata {
				// % GEP 引入
				var GEPSlice = make([]float64, len(sb.SIndiNames))
				for i := 0; i < len(sb.SIndiNames); i++ {
					GEPSlice[i] = SBDE.IndiDataMap[sb.SIndiNames[i]]
				}
				tradeval := Eval(GEPSlice)
				// % GEP 引入完毕
				if !ContainNaN(SBDE.IndiDataMap) {
					PL = append(PL, Pair{Key: instID, Value: tradeval[0], Price: SBDE.IndiDataMap["Close"]})
				}
				// % GEP 应用完毕
			}
			sort.Sort(PL)
			// if debug {
			// 	log.Info().Str("PL details: ", fmt.Sprint(PL)).Str("TimeStamp", datetime).
			// 		Msg("PL details")
			// }
			// 循环PL，得到对应标的的持仓比例
			tmptargetM := orderedmap.NewOrderedMap[string, float64]()
			totalValue := vAcct.SAcct.MktVal * sb.position_ratio
			for i := 0; i < len(sb.holdingRatios); i++ {
				tmptargetM.Set(PL[len(PL)-i-1].Key, math.Floor(((totalValue * sb.holdingRatios[len(sb.holdingRatios)-i-1] / PL[len(PL)-i-1].Price) / contractproperty.SimpleNewSCPFromMap(CPMap, PL[len(PL)-i-1].Key).ContractSize)))
			}
			sb.targetsM = tmptargetM

			// 循环持仓holdingsM，对比targetsM得到卖出方向的订单
			for _, instID := range sb.holdingsM.Keys() {
				tmpHM, _ := sb.holdingsM.Get(instID)
				tmpTM, Tok := sb.targetsM.Get(instID)
				if !Tok {
					// 如果不在目标持仓中，则卖出
					if tmpHM != 0 {
						tmpSCP := contractproperty.SimpleNewSCPFromMap(CPMap, instID)
						if v, ok := bc.Stockdata[instID]; ok {
							new_order := order.NewStockOrder(instID, false, false, datetime, v.IndiDataMap["Close"], tmpHM, "Sell", &tmpSCP)
							orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
						}

					}
				} else if tmpHM > tmpTM {
					// 如果目标持仓小于当前持仓，则卖出对应数量
					tmpSCP := contractproperty.SimpleNewSCPFromMap(CPMap, instID)
					if v, ok := bc.Stockdata[instID]; ok {
						new_order := order.NewStockOrder(instID, false, false, datetime, v.IndiDataMap["Close"], (tmpHM - tmpTM), "Sell", &tmpSCP)
						orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
					}
				}
			}

			// 循环targetsM，对比holdingsM得到买入方向的订单
			for _, instID := range sb.targetsM.Keys() {
				tmpHM, Hok := sb.holdingsM.Get(instID)
				tmpTM, _ := sb.targetsM.Get(instID)
				if !Hok {
					// 如果不在持仓中，则买入
					if tmpTM != 0 {
						tmpSCP := contractproperty.SimpleNewSCPFromMap(CPMap, instID)
						if v, ok := bc.Stockdata[instID]; ok {
							new_order := order.NewStockOrder(instID, false, false, datetime, v.IndiDataMap["Close"], tmpTM, "Buy", &tmpSCP)
							orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
						}
					}
				} else if tmpHM < tmpTM {
					// 如果目标持仓大于当前持仓，则买入对应数量
					tmpSCP := contractproperty.SimpleNewSCPFromMap(CPMap, instID)
					if v, ok := bc.Stockdata[instID]; ok {
						new_order := order.NewStockOrder(instID, false, false, datetime, v.IndiDataMap["Close"], (tmpTM - tmpHM), "Buy", &tmpSCP)
						orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
					}
				}
			}
		}
	}
	// 5. 返回结果
	// if debug {
	// 	log.Info().Str("Order details: ", outputOrderRes(orderRes)).Str("TimeStamp", datetime).
	// 		Msg("Order details")
	// }

	return orderRes
}

func (sb *SortBuyStrategy) ActOnDataMAN(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap contractproperty.CPMap) (orderRes OrderResult) {
	// zerolog.TimeFieldFormat = "2006-01-02"
	// 1. update the virtual account
	UpdateAcct(bc, vAcct)
	// 2. iter vAcct.SAcct.PosMap and maintain holdingsM
	for k, v := range vAcct.SAcct.PosMap {
		// 日线级策略 问题不大
		sb.holdingsM.Set(k, v.CalPosTdyNum()+v.CalPosPrevNum())
		if val, _ := sb.holdingsM.Get(k); val == 0 {
			sb.holdingsM.Delete(k)
		}
	}
	// 3. check the datetime is the executable time
	if GetTimeValue(datetime) == sb.STimeCritic {
		sb.stimeCondition = true
	} else {
		sb.stimeCondition = false
	}
	// 4. 判断股票标的切片SInstrNames是否为空 并且 时间准则为真，如果为空，则不操作股票数据循环
	if len(sb.SInstNames) != 0 && sb.stimeCondition {
		// check if the length of the stockdata is bigger than 1
		if len(bc.Stockdata) < 1 {
			// 4.1 如果为空，则panic
			panic("data error! the length of the stockdata is smaller than 1")
		} else if len(bc.Stockdata) == 1 {
			// 4.2 如果长度为1，大概率为回测初始数据量原因，亦可能为实盘停牌或单纯交易量低，
			for instID, SBDE := range bc.Stockdata {
				tmpSCP := contractproperty.SimpleNewSCPFromMap(CPMap, instID)
				// 期望买入的总市值
				targetValue := vAcct.SAcct.MktVal * sb.position_ratio
				// 期望买入的股票数量，如果符合大于开单条件(此处不为0)同步到targetsM
				tmpTM := math.Floor(((targetValue / SBDE.IndiDataMap["Close"]) / contractproperty.SimpleNewSCPFromMap(CPMap, instID).ContractSize))
				if tmpTM > 0 {
					sb.targetsM.Set(instID, tmpTM)
				}
				// 查看是否账户持有该标的
				tmpHM, Hok := sb.holdingsM.Get(instID)
				if !Hok {
					// 如果不持有该标的，且sb.holdingsM为空，则直接买入
					// （其他情况大概率为实盘数据异常，不操作）
					if sb.holdingsM.Len() == 0 {
						new_order := order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["Close"], tmpTM, "Buy", &tmpSCP)
						if sb.CheckEligible(&new_order, &vAcct.SAcct) {
							orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
						}
					}
				} else {
					// 如果持有该标的，且仅持有该标的，则调整到期望数量
					// （其他情况大概率为实盘数据异常，不操作）
					if sb.holdingsM.Len() == 1 {
						if tmpHM < tmpTM {
							new_order := order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["Close"], math.Abs(tmpTM-tmpHM), "Buy", &tmpSCP)
							if sb.CheckEligible(&new_order, &vAcct.SAcct) {
								orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
							}
						} else if tmpHM > tmpTM {
							new_order := order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["Close"], math.Abs(tmpTM-tmpHM), "Sell", &tmpSCP)
							if sb.CheckEligible(&new_order, &vAcct.SAcct) {
								orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
							}
						}
					}
				}
			}
		} else {
			// 4.3 如果长度大于1，正常循环
			// 得到数据个数
			n := len(bc.Stockdata)
			// 修改holdingRatios
			if n < sb.numHolding {
				sb.holdingRatios = arithmeticSequence(n, sb.max_min_ratio)
			} else {
				sb.holdingRatios = arithmeticSequence(sb.numHolding, sb.max_min_ratio)
			}
			// 循环标的产生PL 标的排序并得到对应标的的持仓比例
			PL := PairList{}
			for instID, SBDE := range bc.Stockdata {
				if !ContainNaN(SBDE.IndiDataMap) {
					PL = append(PL, Pair{Key: instID, Value: SBDE.IndiDataMap["MA5"] / SBDE.IndiDataMap["Close"], Price: SBDE.IndiDataMap["Close"]})
				}
			}
			sort.Sort(PL)
			// if debug {
			// 	log.Info().Str("PL details: ", fmt.Sprint(PL)).Str("TimeStamp", datetime).
			// 		Msg("PL details")
			// }
			// 循环PL，得到对应标的的持仓比例
			tmptargetM := orderedmap.NewOrderedMap[string, float64]()
			totalValue := vAcct.SAcct.MktVal * sb.position_ratio
			for i := 0; i < len(sb.holdingRatios); i++ {
				tmptargetM.Set(PL[len(PL)-i-1].Key, math.Floor(((totalValue * sb.holdingRatios[len(sb.holdingRatios)-i-1] / PL[len(PL)-i-1].Price) / contractproperty.SimpleNewSCPFromMap(CPMap, PL[len(PL)-i-1].Key).ContractSize)))
			}
			sb.targetsM = tmptargetM

			// 循环持仓holdingsM，对比targetsM得到卖出方向的订单
			for _, instID := range sb.holdingsM.Keys() {
				tmpHM, _ := sb.holdingsM.Get(instID)
				tmpTM, Tok := sb.targetsM.Get(instID)
				if !Tok {
					// 如果不在目标持仓中，则卖出
					if tmpHM != 0 {
						tmpSCP := contractproperty.SimpleNewSCPFromMap(CPMap, instID)
						if v, ok := bc.Stockdata[instID]; ok {
							new_order := order.NewStockOrder(instID, false, false, datetime, v.IndiDataMap["Close"], tmpHM, "Sell", &tmpSCP)
							if sb.CheckEligible(&new_order, &vAcct.SAcct) {
								orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
							}
						}

					}
				} else if tmpHM > tmpTM {
					// 如果目标持仓小于当前持仓，则卖出对应数量
					tmpSCP := contractproperty.SimpleNewSCPFromMap(CPMap, instID)
					if v, ok := bc.Stockdata[instID]; ok {
						new_order := order.NewStockOrder(instID, false, false, datetime, v.IndiDataMap["Close"], (tmpHM - tmpTM), "Sell", &tmpSCP)
						if sb.CheckEligible(&new_order, &vAcct.SAcct) {
							orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
						}
					}
				}
			}

			// 循环targetsM，对比holdingsM得到买入方向的订单
			for _, instID := range sb.targetsM.Keys() {
				tmpHM, Hok := sb.holdingsM.Get(instID)
				tmpTM, _ := sb.targetsM.Get(instID)
				if !Hok {
					// 如果不在持仓中，则买入
					if tmpTM != 0 {
						tmpSCP := contractproperty.SimpleNewSCPFromMap(CPMap, instID)
						if v, ok := bc.Stockdata[instID]; ok {
							new_order := order.NewStockOrder(instID, false, false, datetime, v.IndiDataMap["Close"], tmpTM, "Buy", &tmpSCP)
							if sb.CheckEligible(&new_order, &vAcct.SAcct) {
								orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
							}
						}
					}
				} else if tmpHM < tmpTM {
					// 如果目标持仓大于当前持仓，则买入对应数量
					tmpSCP := contractproperty.SimpleNewSCPFromMap(CPMap, instID)
					if v, ok := bc.Stockdata[instID]; ok {
						new_order := order.NewStockOrder(instID, false, false, datetime, v.IndiDataMap["Close"], (tmpTM - tmpHM), "Buy", &tmpSCP)
						if sb.CheckEligible(&new_order, &vAcct.SAcct) {
							orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
						}
					}
				}
			}
		}
	}
	// 5. 返回结果
	// if debug {
	// 	log.Info().Str("Order details: ", outputOrderRes(orderRes)).Str("TimeStamp", datetime).
	// 		Msg("Order details")
	// }

	return orderRes
}

// 函数：give an arithmetic sequence with n elements, whose all elements add up equals to 100, and the first element is 2 times bigger than the last element
func arithmeticSequence(n int, ratio float64) []float64 {
	// 1. create a slice with n elements
	slice := make([]float64, n)
	if n < 1 {
		panic("n must be bigger than 1")
	} else if n == 1 {
		slice[0] = 1
		return slice
	} else {
		// 2. calculate the first element
		first_element := 2.0 / (float64(n) * (1 + ratio))
		// 3. calculate the difference between the first element and the last element
		diff := first_element*ratio - first_element
		// 4. calculate the common difference
		d := diff / float64(n-1)
		// 5. calculate the arithmetic sequence
		for i := 0; i < n; i++ {
			slice[i] = first_element + float64(i)*d
		}
		return slice
	}
}

// func output  a map[string]*stockaccount.PositionSlice as string
// for debug use
// func outputMap(m map[string]*stockaccount.PositionSlice) string {
// 	var output string
// 	for k, v := range m {
// 		output += fmt.Sprintf("instID: %s, tdyPosition: %v, prePosition: %v\n", k, v.CalPosTdyNum(), v.CalPosPrevNum())
// 	}
// 	return output
// }

// func outputOrderRes(ors OrderResult) string {
// 	var output string
// 	for _, v := range ors.StockOrderS {
// 		// output IsEligible  IsExecuted OrderTime OrderPrice OrderNum
// 		output += fmt.Sprintf("instID: %s, IsEligible: %s, IsExecuted: %v, OrderTime: %s, OrderPrice: %v, OrderNum: %v, OrderDirection: %s\n", v.InstID, v.IsEligible, v.IsExecuted, v.OrderTime, v.OrderPrice, v.OrderNum, v.OrderDirection)
// 	}
// 	return output
// }
