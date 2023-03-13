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
	"github.com/wonderstone/QuantTools/contractproperty"
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
	s_cash         map[string]float64 //各支股票对应现金
	min_cash_ratio float64
}

func NewDMTStrategy(init_money float64, SInstNms []string, SIndiNms []string, STimeCritic string, s_cash map[string]float64) DMTStrategy {
	return DMTStrategy{
		init_money:     init_money,
		SInstNames:     SInstNms,
		SIndiNames:     SIndiNms,
		STimeCritic:    STimeCritic,
		lastTradeValue: make(map[string]float64),
		stimeCondition: false,
		s_cash:         s_cash,
		min_cash_ratio: 0.05,
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
	s_cash := make(map[string]float64, len(sinstrnames))
	average, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(init_money)/float64(len(sinstrnames))), 64)
	for _, s := range sinstrnames {

		s_cash[s] = average
	}
	return NewDMTStrategy(init_money, sinstrnames, sindinames, STimeCritic, s_cash)
}
func GetTimeValue(timeString string) string {
	// get the time value
	timeValue := timeString[11:16]
	return timeValue
}
func update_value(dmt *DMTStrategy, vAcct *virtualaccount.VAcct) {
	for InstID, cash := range dmt.s_cash {
		order_info, ok := vAcct.SAcct.RecordOrderMapS[InstID]
		if ok {
			for _, info := range order_info {
				if info.SO.OrderDirection == "Buy" {
					cash = cash - info.SO.CalEquity()
				} else if info.SO.OrderDirection == "Sell" {
					cash = cash + info.SO.CalEquity()
				}
			}
			vAcct.SAcct.RecordOrderMapS[InstID] = nil
		}
		dmt.s_cash[InstID] = cash
	}
}

// check the order  if eligible
func (dmt *DMTStrategy) CheckEligible(s_name string, o *order.StockOrder, SA *stockaccount.StockAccount) {
	switch o.OrderDirection {
	case "Buy":
		if o.CalEquity() <= dmt.s_cash[s_name] {
			o.IsEligible = true
		}
	case "Sell":
		o.IsEligible = true
	}
}
func (dmt *DMTStrategy) ActOnData(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap, Eval func([]float64) []float64) (orderRes OrderResult) {
	update_value(dmt, vAcct)
	//check the datetime is the executable time
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
				// * GEP 引入
				var GEPSlice = make([]float64, len(dmt.SInstNames))
				for i := 0; i < len(dmt.SInstNames); i++ {
					GEPSlice[i] = SBDE.IndiDataMap[dmt.SInstNames[i]]
				}
				tmpSCP := cp.SimpleNewSCPFromMap(CPMap, instID)

				tradeval := Eval(GEPSlice)
				lsttv, tok := dmt.lastTradeValue[instID]
				//buy condition check
				if tok && lsttv < 0 && tradeval[0] > 0 {
					//使用该股票可支配的80%资金计算buy_num
					buy_num := math.Floor(((dmt.s_cash[instID] / SBDE.IndiDataMap["Close"]) / contractproperty.SimpleNewSCPFromMap(CPMap, instID).ContractSize) * (1 - dmt.min_cash_ratio))
					if buy_num != 0 && SBDE.IndiDataMap["Close"] != 0 {
						new_order := order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["Close"], buy_num, "Buy", &tmpSCP)
						dmt.CheckEligible(instID, &new_order, &vAcct.SAcct)
						if new_order.IsEligible {
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
					dmt.lastTradeValue[instID] = tradeval[0]
				}
				//sell condition check
				if tok && tradeval[0] < 0 && lsttv > 0 {
					if _, ok := vAcct.SAcct.PosMap[instID]; ok {
						if vAcct.SAcct.PosMap[instID].CalPosPrevNum() > 0 {
							if SBDE.IndiDataMap["Close"] != 0 {
								new_order := order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["Close"], vAcct.SAcct.PosMap[instID].CalPosPrevNum(), "Sell", &tmpSCP)
								dmt.CheckEligible(instID, &new_order, &vAcct.SAcct)
								if new_order.IsEligible {
									orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
								}
							}
						}
					}
					if debug {
						log.Info().Str("Account UUID", vAcct.SAcct.UUID).Str("TimeStamp", datetime).
							Float64("Close", SBDE.IndiDataMap["Close"]).Float64("Open", SBDE.IndiDataMap["Open"]).Str("InstID", instID).
							Msg("Strategy sell")
					}
					dmt.lastTradeValue[instID] = tradeval[0]
				}

			}
		}
	}
	return orderRes
}
func (dmt *DMTStrategy) ActOnDataMAN(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap) (orderRes OrderResult) {
	update_value(dmt, vAcct)
	//check the datetime is the executable time
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
				tmpSCP := cp.SimpleNewSCPFromMap(CPMap, instID)
				tradeval := SBDE.IndiDataMap["MA3"] - SBDE.IndiDataMap["MA5"]
				fmt.Println(SBDE.IndiDataMap["MA3"], SBDE.IndiDataMap["MA5"])
				if tradeval < 0 {
					fmt.Println("有卖的时机")
				}
				lsttv, tok := dmt.lastTradeValue[instID]

				//buy condition check
				if tok && lsttv < 0 && tradeval > 0 {
					//使用该股票可支配的80%资金计算buy_num
					buy_num := math.Floor(((dmt.s_cash[instID] / SBDE.IndiDataMap["Close"]) / contractproperty.SimpleNewSCPFromMap(CPMap, instID).ContractSize) * 0.8)
					if buy_num != 0 && SBDE.IndiDataMap["Close"] != 0 {
						new_order := order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["Close"], buy_num, "Buy", &tmpSCP)
						dmt.CheckEligible(instID, &new_order, &vAcct.SAcct)
						if new_order.IsEligible {
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
				//sell condition check
				if tok && tradeval < 0 && lsttv > 0 {
					if _, ok := vAcct.SAcct.PosMap[instID]; ok {
						if vAcct.SAcct.PosMap[instID].CalPosPrevNum() > 0 {
							if SBDE.IndiDataMap["Close"] != 0 {
								new_order := order.NewStockOrder(instID, false, false, datetime, SBDE.IndiDataMap["Close"], vAcct.SAcct.PosMap[instID].CalPosPrevNum(), "Sell", &tmpSCP)
								dmt.CheckEligible(instID, &new_order, &vAcct.SAcct)
								if new_order.IsEligible {
									orderRes.StockOrderS = append(orderRes.StockOrderS, new_order)
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
