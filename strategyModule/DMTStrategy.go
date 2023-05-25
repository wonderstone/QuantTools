// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  CheYang (Digital Office Product Department #2)
package strategyModule

import (
	"fmt"
	"math"

	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/wonderstone/QuantTools/account/stockaccount"
	"github.com/wonderstone/QuantTools/account/virtualaccount"
	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/dataprocessor"

	"github.com/wonderstone/QuantTools/order"

	"github.com/wonderstone/QuantTools/configer"
)

type DMTStrategy struct {
	init_money     float64            //初始资金
	SInstNames     []string           // 股票标的名称
	SIndiNames     []string           // 股票参与GEP指标名称，注意其数量不大于BarDE内信息数量，且strategy内可见BarDE的数据
	STimeCritic    string             // 时间关键字，用于判断是否需要进行交易
	lastTradeValue map[string]float64 // 公式上一次买入信号数值
	stimeCondition bool               // 时间条件是否满足
	faMap          map[string]float64 //各支股票对应可用资金  fund available
	min_cash_ratio float64
}

func NewDMTStrategy(init_money float64, SInstNms []string, SIndiNms []string, STimeCritic string, faMap map[string]float64, min_cash_ratio float64) DMTStrategy {
	return DMTStrategy{
		init_money:     init_money,
		SInstNames:     SInstNms,
		SIndiNames:     SIndiNms,
		STimeCritic:    STimeCritic,
		lastTradeValue: make(map[string]float64),
		stimeCondition: false,
		faMap:          faMap,
		min_cash_ratio: min_cash_ratio,
	}
}

// = this function is designed for the framework
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
	// init the faMap
	faMap := make(map[string]float64, len(sinstrnames))
	average, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(init_money)/float64(len(sinstrnames))), 64)
	for _, s := range sinstrnames {
		// 账户初始资金充当市值情况
		faMap[s] = average
	}
	min_cash_ratio := c.GetFloat64("default.min_cash_ratio")
	return NewDMTStrategy(init_money, sinstrnames, sindinames, STimeCritic, faMap, min_cash_ratio)
}

// = 字符串简易处理时间判断
func GetTimeValue(timeString string) string {
	// get the time value
	timeValue := timeString[11:16]
	return timeValue
}

// ~ 利用fundavailable概念，卖出时增加该数值，买入时其充当管理市值的剩余现金概念
func update_value(dmt *DMTStrategy, vAcct *virtualaccount.VAcct) {
	for InstID, fa := range dmt.faMap {
		order_info, ok := vAcct.SAcct.RecordOrderMapS[InstID]
		if ok {
			for _, info := range order_info {
				if info.SO.OrderDirection == "Buy" {
					fa = fa - info.SO.CalEquity()
				} else if info.SO.OrderDirection == "Sell" {
					fa = fa + info.SO.CalEquity()
				}
			}
			vAcct.SAcct.RecordOrderMapS[InstID] = nil
		}
		dmt.faMap[InstID] = fa
	}
}

// check the order if eligible
// / CheckEligible机制为何不能删除？
// * 原因之一：虚拟账户与真实账户下单同步机制。策略下单实际上基于虚拟账户，下单的一刻即可以同步然后模拟撮合更新虚拟账户。
// *         真实账户情况与虚拟账户可能完全不同。虚拟撮合可能存在虚拟下单失败，但是真实账户下单成功的情况。
// / 为何不使用虚拟撮合后同步下单操作？
// * 原因之一：撮合需要下一笔数据进入进行判断。这会导致下单时间延迟，且逻辑上需要特判该笔数据。
func (dmt *DMTStrategy) CheckEligible(s_name string, o *order.StockOrder, SA *stockaccount.StockAccount) bool {
	switch o.OrderDirection {
	case "Buy":
		if o.CalEquity() <= dmt.faMap[s_name] {
			return true
		}
	case "Sell":
		if _, ok := SA.PosMap[o.InstID]; ok {
			// check the previous position is enough
			if o.OrderNum <= SA.PosMap[o.InstID].CalPosPrevNum() {
				return true
			}
		}
	}
	return false
}

// ~ GEP Style with only 5 lines different from the Manually Style
func (dmt *DMTStrategy) ActOnData(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap, Eval func([]float64) []float64) (orderRes OrderResult) {
	update_value(dmt, vAcct)
	// 1. check the datetime is the exact time to trade
	if GetTimeValue(datetime) == dmt.STimeCritic {
		dmt.stimeCondition = true
	} else {
		dmt.stimeCondition = false
	}
	// 判断股票标的切片SInstrNames是否为空 并且 时间准则为真，如果为空，则不操作股票数据循环
	if len(dmt.SInstNames) != 0 && dmt.stimeCondition {
		// 依据标的循环Data得到数据
		for instID, SBDE := range bc.Stockdata {
			// 判断是否数据为NaN，如果为NaN，则跳过
			if !ContainNaN(SBDE.IndiDataMap) {
				// % GEP 引入
				var GEPSlice = make([]float64, len(dmt.SIndiNames))
				for i := 0; i < len(dmt.SIndiNames); i++ {
					GEPSlice[i] = SBDE.IndiDataMap[dmt.SIndiNames[i]]
				}
				tradeval := Eval(GEPSlice)
				// % GEP 引入完毕
				tmpSCP := cp.SimpleNewSCPFromMap(CPMap, instID)
				// + 查看上一次记录的交易依托数值
				lsttv, tok := dmt.lastTradeValue[instID]
				//buy condition check，如果上一次交易依托数值为负数，且当前交易依托数值为正数，则进行买入操作
				if tok && lsttv < 0 && tradeval[0] > 0 {
					// 判断数据SBDE.IndiDataMap是否包含Close，如不包含，则跳过
					if val, ok := SBDE.IndiDataMap["Close"]; ok {
						//使用该股票可支配的(1 - dmt.min_cash_ratio)资金计算buy_num
						buy_num := math.Floor(((dmt.faMap[instID] / val) / cp.SimpleNewSCPFromMap(CPMap, instID).ContractSize) * (1 - dmt.min_cash_ratio))
						if buy_num != 0 && val != 0 {
							// buy_num != 0是可以交易的前提，SBDE.IndiDataMap["Close"] != 0是避免数据不包含的情况。
							new_order := order.NewStockOrder(instID, false, false, datetime, val, buy_num, "Buy", &tmpSCP)
							if dmt.CheckEligible(instID, &new_order, &vAcct.SAcct) {
								orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
							}
							// DCE: debug info
							if debug {
								// this part is for test only
								log.Info().Str("Account UUID", vAcct.SAcct.UUID).Str("TimeStamp", datetime).Float64("Close", SBDE.IndiDataMap["Close"]).
									Float64("Open", SBDE.IndiDataMap["Open"]).Str("InstID", instID).
									Msg("Strategy buy")
							}
						}
					}
				}
				//sell condition check
				if tok && tradeval[0] < 0 && lsttv > 0 {
					if _, ok := vAcct.SAcct.PosMap[instID]; ok {
						// daily strategy so no position today. sell them all
						if vAcct.SAcct.PosMap[instID].CalPosPrevNum() > 0 {
							// 判断数据SBDE.IndiDataMap是否包含Close，如不包含，则跳过
							if val, ok := SBDE.IndiDataMap["Close"]; ok {
								// 以防极端情况 一般val不会为0
								if val != 0 {
									new_order := order.NewStockOrder(instID, false, false, datetime, val, vAcct.SAcct.PosMap[instID].CalPosPrevNum(), "Sell", &tmpSCP)
									if dmt.CheckEligible(instID, &new_order, &vAcct.SAcct) {
										orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
									}
								}
							}
						}
					}
					if debug {
						log.Info().Str("Account UUID", vAcct.SAcct.UUID).Str("TimeStamp", datetime).
							Float64("Close", SBDE.IndiDataMap["Close"]).Float64("Open", SBDE.IndiDataMap["Open"]).Str("InstID", instID).
							Msg("Strategy sell")
					}
				}
				dmt.lastTradeValue[instID] = tradeval[0]
			}
		}
	}
	return orderRes
}

// ~ Manually Style trade logic part
func (dmt *DMTStrategy) ActOnDataMAN(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap) (orderRes OrderResult) {
	update_value(dmt, vAcct)
	// 1. check the datetime is the exact time to trade
	if GetTimeValue(datetime) == dmt.STimeCritic {
		dmt.stimeCondition = true
	} else {
		dmt.stimeCondition = false
	}
	// 判断股票标的切片SInstrNames是否为空 并且 时间准则为真，如果为空，则不操作股票数据循环
	if len(dmt.SInstNames) != 0 && dmt.stimeCondition {
		// 依据标的循环Data得到数据
		for instID, SBDE := range bc.Stockdata {
			// 判断是否数据为NaN，如果为NaN，则跳过
			if !ContainNaN(SBDE.IndiDataMap) {
				// % Manually logic definition
				tradeval := SBDE.IndiDataMap["MA3"] - SBDE.IndiDataMap["MA5"]
				// % Manually logic definition end
				tmpSCP := cp.SimpleNewSCPFromMap(CPMap, instID)
				// + 查看上一次记录的交易依托数值
				lsttv, tok := dmt.lastTradeValue[instID]
				//buy condition check
				if tok && lsttv < 0 && tradeval > 0 {
					// 判断数据SBDE.IndiDataMap是否包含Close，如不包含，则跳过
					if val, ok := SBDE.IndiDataMap["Close"]; ok {
						//使用该股票可支配的80%资金计算buy_num
						buy_num := math.Floor(((dmt.faMap[instID] / val) / cp.SimpleNewSCPFromMap(CPMap, instID).ContractSize) * (1 - dmt.min_cash_ratio))
						if buy_num != 0 && val != 0 {
							new_order := order.NewStockOrder(instID, false, false, datetime, val, buy_num, "Buy", &tmpSCP)
							if dmt.CheckEligible(instID, &new_order, &vAcct.SAcct) {
								orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
							}
							// DCE: debug info
							if debug {
								// this part is for test only
								log.Info().Str("Account UUID", vAcct.SAcct.UUID).Str("TimeStamp", datetime).
									Float64("Close", SBDE.IndiDataMap["Close"]).Float64("Open", SBDE.IndiDataMap["Open"]).Str("InstID", instID).
									Msg("Strategy buy")
							}
						}
					}
				}
				//sell condition check
				if tok && tradeval < 0 && lsttv > 0 {
					if _, ok := vAcct.SAcct.PosMap[instID]; ok {
						// daily strategy so no position today. sell them all
						if vAcct.SAcct.PosMap[instID].CalPosPrevNum() > 0 {
							// 判断数据SBDE.IndiDataMap是否包含Close，如不包含，则跳过
							if val, ok := SBDE.IndiDataMap["Close"]; ok {
								// 以防极端情况 一般val不会为0
								if val != 0 {
									new_order := order.NewStockOrder(instID, false, false, datetime, val, vAcct.SAcct.PosMap[instID].CalPosPrevNum(), "Sell", &tmpSCP)
									if dmt.CheckEligible(instID, &new_order, &vAcct.SAcct) {
										orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
									}
								}
							}
						}
					}
					if debug {
						log.Info().Str("Account UUID", vAcct.SAcct.UUID).Str("TimeStamp", datetime).
							Float64("Close", SBDE.IndiDataMap["Close"]).Float64("Open", SBDE.IndiDataMap["Open"]).Str("InstID", instID).
							Msg("Strategy sell")
					}
				}
				dmt.lastTradeValue[instID] = tradeval
			}
		}
	}
	return orderRes
}
