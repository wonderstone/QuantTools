// // All rights reserved. This is part of West Securities ltd. proprietary source code.
// // No part of this file may be reproduced or transmitted in any form or by any means,
// // electronic or mechanical, including photocopying, recording, or by any information
// // storage and retrieval system, without the prior written permission of West Securities ltd.

// // author:  Zhangweixuan (Digital Office Product Department #2)
package strategyModule

// import (
// 	"golang/account/virtualaccount"
// 	"golang/configer"
// 	cp "golang/contractproperty"
// 	"golang/dataprocessor"
// 	"golang/order"
// )

// // Strategy single object strategy
// type Strategy struct {
// 	// T0 part
// 	InstNames     []string           // 合约的名称
// 	IndiNames     []string           // 合约参与GEP指标名称，注意其数量不大于BarDE内信息数量，且strategy内可见BarDE的数据
// 	Tlimit        int                // 交易次数限制
// 	tCounter      map[string]int     // 交易次数计数器(内部)
// 	tState        int                // 交易状态(内部) 0:未买入 1:已买入 -1:已卖出
// 	InitBuyTime   string             // 初始化买入时间
// 	StartTime     string             // 交易开始时间
// 	StopTime      string             // 交易结束时间
// 	lastDate      string             // 上一交易日日期(内部)
// 	lastBuyValue  map[string]float64 // 公式
// 	lastSellValue map[string]float64 // 公式上一次卖出信号数值
// 	// future part
// 	checkMain bool    // 是否检查过主力合约
// 	Overnight bool    // 是否持仓过夜
// 	checkInit float64 // 防止开仓的钱被平掉
// 	Name      string  // 本策略负责的标的
// 	MainID    string  // 主力合约的名称
// 	OperateID string  // 当前操作的合约
// }

// func NewStrategy(InstNames []string, IndiNames []string, Tlimit int, overnight bool, Name, InitBuyTime, StartTime, StopTime string) Strategy {
// 	tCounter := make(map[string]int)
// 	for _, name := range InstNames {
// 		tCounter[name] = 0
// 	}

// 	return Strategy{
// 		InstNames:     InstNames,
// 		IndiNames:     IndiNames,
// 		Name:          Name,
// 		Tlimit:        Tlimit,
// 		tCounter:      tCounter,
// 		tState:        0,
// 		InitBuyTime:   InitBuyTime,
// 		StartTime:     StartTime,
// 		StopTime:      StopTime,
// 		checkMain:     false,
// 		Overnight:     overnight,
// 		lastBuyValue:  make(map[string]float64),
// 		lastSellValue: make(map[string]float64),

// 		OperateID: "NIL",
// 	}
// }

// func NewStrategyFromConfig(dir string, BTConfile string, sec string, StgConfile string) Strategy {
// 	c := configer.New(dir + BTConfile)
// 	err := c.Load()
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = c.Unmarshal()
// 	if err != nil {
// 		panic(err)
// 	}
// 	tmpMap := c.GetStringMap(sec)
// 	var sinstrnames []string
// 	for _, v := range tmpMap["sinstrnames"].([]interface{}) {
// 		sinstrnames = append(sinstrnames, v.(string))
// 	}
// 	var sindinames []string
// 	for _, v := range tmpMap["sindinames"].([]interface{}) {
// 		sindinames = append(sindinames, v.(string))
// 	}
// 	c = configer.New(dir + StgConfile)
// 	err = c.Load()
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = c.Unmarshal()

// 	return NewStrategy(sinstrnames, sindinames, c.GetInt(sec+".Tlimit"), c.GetBool(sec+".overnight"), c.GetString(sec+".Name"), c.GetString(sec+".InitBuyTime"), c.GetString(sec+".StartTime"), c.GetString(sec+".StopTime"))
// }

// func (s *Strategy) ActOnData(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap, Eval func([]float64) []float64) (orderRes OrderResult) {
// 	//获取当前时间并判断是否介于StartTime与StopTime,是则进行操作
// 	if datetime[11:] >= s.StartTime && datetime[11:] <= s.StopTime {
// 		// 检查主力合约，并且添加/删除合约,每天只算一次
// 		if !s.checkMain {
// 			tmp := 0.0
// 			var tmpID []string
// 			for instID, SBDE := range bc.Futuresdata {
// 				if !ContainNaN(SBDE.IndiDataMap) {
// 					tmpID = append(tmpID, instID)
// 					if SBDE.IndiDataMap["Vol"] > tmp {
// 						s.MainID = instID
// 						tmp = SBDE.IndiDataMap["Vol"]
// 					}
// 				}
// 			}
// 			s.IndiNames = tmpID
// 			s.checkMain = true
// 		}

// 		// 如果当前操作的不是主力合约，则平仓后换仓
// 		if s.OperateID != s.MainID {
// 			// 初始化不平仓
// 			if s.OperateID != "NIL" {
// 				// 平仓
// 				tmpFCP := cp.SimpleNewFCPFromMap(CPMap, s.OperateID)
// 				SBDE := bc.Futuresdata[s.OperateID]
// 				postL, postS := vAcct.FAcct.PosMap[s.OperateID].CalPosPrevNum()
// 				netNum := postL - postS
// 				if netNum > 0 {
// 					orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(s.OperateID, false, false, datetime, SBDE.IndiDataMap["Close"], float64(netNum), "Buy", "ClosePrevious", &tmpFCP))
// 				} else if netNum < 0 {
// 					orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(s.OperateID, false, false, datetime, SBDE.IndiDataMap["Close"], float64(-1*netNum), "Sell", "ClosePrevious", &tmpFCP))
// 				}
// 				}
// 			}
// 			// 80% 开仓
// 			s.OperateID = s.MainID
// 			tmpFCP := cp.SimpleNewFCPFromMap(CPMap, s.OperateID)
// 			SBDE := bc.Futuresdata[s.OperateID]
// 			s.checkInit = vAcct.FAcct.Fundavail * 0.8
// 			orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(s.OperateID, false, false, datetime, SBDE.IndiDataMap["Close"], vAcct.FAcct.Fundavail*0.8, "Buy", "Open", &tmpFCP))
// 			s.tState = 0
// 			s.lastBuyValue[s.OperateID] = 0
// 			s.lastSellValue[s.OperateID] = 0
// 		} else { // 否则进行T0常规操作（GEP引入）
// 			SBDE := bc.Futuresdata[s.OperateID]
// 			tmpFCP := cp.SimpleNewSCPFromMap(CPMap, s.OperateID)
// 			var GEPSlice = make([]float64, len(s.IndiNames))
// 			for i := 0; i < len(s.IndiNames); i++ {
// 				GEPSlice[i] = SBDE.IndiDataMap[s.IndiNames[i]]
// 			}
// 			tmps := Eval(GEPSlice)
// 			buyval := tmps[0]
// 			lstbv, bok := s.lastBuyValue[s.OperateID]
// 			sellval := tmps[1]
// 			lstsv, sok := s.lastSellValue[s.OperateID]
// 			TradeN := int(vAcct.FAcct.Fundavail) / s.Tlimit
// 			if bok && sok {
// 				if buyval > 0 && lstbv <= 0 && s.tState <= 0 && s.tCounter[s.OperateID] <= s.Tlimit*2 {
// 					// 买入TradeN
// 					newOrder := order.NewFuturesOrder(s.OperateID, false, false, datetime, SBDE.IndiDataMap["Close"], float64(TradeN), "Buy", "Open", &tmpFCP)
// 					if newOrder.IsEligible {
// 						orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, newOrder)
// 						// tCounter计数器+1
// 						s.tCounter[s.OperateID] += 1

// 						// lastBuyValue 赋值 lastSellValue 赋值
// 						s.lastBuyValue[s.OperateID] = buyval
// 						s.lastSellValue[s.OperateID] = sellval

// 						// tState 赋值
// 						s.tState = 1
// 					}
// 				}
// 				if sellval >= 0 && lstsv <= 0 && s.tState >= 0 && s.tCounter[s.OperateID] <= s.Tlimit*2 {
// 					// 卖出TradeN
// 					newOrder := order.NewFuturesOrder(s.OperateID, false, false, datetime, SBDE.IndiDataMap["Close"], float64(TradeN), "Sell", "CloseToday", &tmpFCP)
// 					if newOrder.IsEligible {
// 						orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, newOrder)
// 						// tCounter计数器+1
// 						s.tCounter[s.OperateID] += 1

// 						// lastBuyValue 赋值 lastSellValue 赋值
// 						s.lastBuyValue[s.OperateID] = buyval
// 						s.lastSellValue[s.OperateID] = sellval

// 						// tState 赋值
// 						s.tState = -1
// 					}

// 				}
// 			}
// 		}
// 	} else if datetime[11:] > s.StopTime {
// 		// 非交易时间更新checkMain
// 		s.checkMain = false
// 		// 持仓不过夜的话每天都要平仓
// 		if !s.Overnight {
// 			SBDE := bc.Futuresdata[s.OperateID]
// 			tmpFCP := cp.SimpleNewFCPFromMap(CPMap, s.OperateID)
// 			tdL, tdS := vAcct.FAcct.PosMap[s.OperateID].CalPosTdyNum()
// 			// 第一天需要防止把开仓的钱平掉
// 			if s.checkInit != 0 {
// 				tdL -= s.checkInit
// 				s.checkInit = 0
// 			}
// 			if tdL-tdS > 0 {
// 				orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(s.OperateID, false, false, datetime, SBDE.IndiDataMap["Close"], tdL-tdS, "Sell", "CloseToday", &tmpFCP))
// 			} else if tdL-tdS < 0 {
// 				orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(s.OperateID, false, false, datetime, SBDE.IndiDataMap["Close"], tdS-tdL, "Buy", "CloseToday", &tmpFCP))
// 			}
// 		}
// 	}
// 	s.lastDate = datetime[:10]
// 	return
// }

// // ActOnDataMAN Manual logic to check the correctness of the strategy
// func (s *Strategy) ActOnDataMAN(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap) (orderRes OrderResult) {
// 	//获取当前时间并判断是否介于StartTime与StopTime,是则进行操作
// 	if datetime[11:] >= s.StartTime && datetime[11:] <= s.StopTime {
// 		// 检查主力合约，并且添加/删除合约,每天只算一次
// 		if !s.checkMain {
// 			tmp := 0.0
// 			var tmpID []string
// 			for instID, SBDE := range bc.Futuresdata {
// 				if !ContainNaN(SBDE.IndiDataMap) {
// 					tmpID = append(tmpID, instID)
// 					if SBDE.IndiDataMap["Vol"] > tmp {
// 						s.MainID = instID
// 						tmp = SBDE.IndiDataMap["Vol"]
// 					}
// 				}
// 			}
// 			s.IndiNames = tmpID
// 			s.checkMain = true
// 		}

// 		// 如果当前操作的不是主力合约，则平仓后换仓
// 		if s.OperateID != s.MainID {
// 			// 初始化不平仓
// 			if s.OperateID != "NIL" {
// 				tmpFCP := cp.SimpleNewFCPFromMap(CPMap, s.OperateID)
// 				SBDE := bc.Futuresdata[s.OperateID]
// 				// 平仓
// 				postL, postS := vAcct.FAcct.PosMap[s.OperateID].CalPosPrevNum()
// 				netNum := postL - postS
// 				if netNum > 0 {
// 					orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(s.OperateID, false, false, datetime, SBDE.IndiDataMap["Close"], float64(netNum), "Buy", "ClosePrevious", &tmpFCP))
// 				} else if netNum < 0 {
// 					orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(s.OperateID, false, false, datetime, SBDE.IndiDataMap["Close"], float64(-1*netNum), "Sell", "ClosePrevious", &tmpFCP))
// 				}
// 			}
// 			// 80% 开仓
// 			s.OperateID = s.MainID
// 			tmpFCP := cp.SimpleNewFCPFromMap(CPMap, s.OperateID)
// 			SBDE := bc.Futuresdata[s.OperateID]
// 			s.checkInit = vAcct.FAcct.Fundavail * 0.8
// 			orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(s.OperateID, false, false, datetime, SBDE.IndiDataMap["Close"], vAcct.FAcct.Fundavail*0.8, "Buy", "Open", &tmpFCP))
// 			s.tState = 0
// 			s.lastBuyValue[s.OperateID] = 0
// 			s.lastSellValue[s.OperateID] = 0
// 		} else { // 否则进行T0常规操作（手动）
// 			SBDE := bc.Futuresdata[s.OperateID]
// 			tmpFCP := cp.SimpleNewFCPFromMap(CPMap, s.OperateID)
// 			tradeval := SBDE.IndiDataMap["MA3"] - SBDE.IndiDataMap["MA5"]
// 			lsttv, tok := s.lastBuyValue[s.OperateID]
// 			TradeN := int(vAcct.FAcct.Fundavail) / s.Tlimit
// 			if tok && lsttv <= 0 && tradeval > 0 && s.tCounter[s.OperateID] <= s.Tlimit*2 {
// 				// 买入TradeN
// 				newOrder := order.NewFuturesOrder(s.OperateID, false, false, datetime, SBDE.IndiDataMap["Close"], float64(TradeN), "Buy", "Open", &tmpFCP)
// 				if newOrder.IsEligible {
// 					orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, newOrder)
// 					// tCounter计数器+1
// 					s.tCounter[s.OperateID] += 1

// 					// lastBuyValue 赋值 lastSellValue 赋值
// 					s.lastBuyValue[s.OperateID] = tradeval
// 					s.lastSellValue[s.OperateID] = tradeval

// 					// tState 赋值
// 					s.tState = 1
// 				}
// 			}
// 			if tok && tradeval < 0 && lsttv >= 0 && s.tCounter[s.OperateID] <= s.Tlimit*2 {
// 				// 卖出TradeN
// 				newOrder := order.NewFuturesOrder(s.OperateID, false, false, datetime, SBDE.IndiDataMap["Close"], float64(TradeN), "Sell", "CloseToday", &tmpFCP)
// 				if newOrder.IsEligible {
// 					orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, newOrder)
// 					// tCounter计数器+1
// 					s.tCounter[s.OperateID] += 1

// 					// lastBuyValue 赋值 lastSellValue 赋值
// 					s.lastBuyValue[s.OperateID] = tradeval
// 					s.lastSellValue[s.OperateID] = tradeval

// 					// tState 赋值
// 					s.tState = -1
// 				}
// 			}
// 		}
// 	} else if datetime[11:] > s.StopTime {
// 		// 非交易时间更新checkMain
// 		s.checkMain = false
// 		// 持仓不过夜的话每天都要平仓
// 		if !s.Overnight {
// 			SBDE := bc.Futuresdata[s.OperateID]
// 			tmpFCP := cp.SimpleNewFCPFromMap(CPMap, s.OperateID)
// 			tdL, tdS := vAcct.FAcct.PosMap[s.OperateID].CalPosTdyNum()
// 			// 第一天需要防止把开仓的钱平掉
// 			if s.checkInit != 0 {
// 				tdL -= s.checkInit
// 				s.checkInit = 0
// 			}
// 			if tdL-tdS > 0 {
// 				orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(s.OperateID, false, false, datetime, SBDE.IndiDataMap["Close"], tdL-tdS, "Sell", "CloseToday", &tmpFCP))
// 			} else if tdL-tdS < 0 {
// 				orderRes.FuturesOrderS = append(orderRes.FuturesOrderS, order.NewFuturesOrder(s.OperateID, false, false, datetime, SBDE.IndiDataMap["Close"], tdS-tdL, "Buy", "CloseToday", &tmpFCP))
// 			}
// 		}
// 	}
// 	s.lastDate = datetime[:10]
// 	return
// }
