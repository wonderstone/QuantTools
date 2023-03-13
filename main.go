package main

import (
	"flag"

	"github.com/rs/zerolog/log"
	"github.com/wonderstone/QuantTools/account/virtualaccount"
	"github.com/wonderstone/QuantTools/framework"
	"github.com/wonderstone/QuantTools/strategyModule"
)

type manager struct {
	BT  *framework.BackTest
	STG strategyModule.IStrategy
}

const debug = true

// * Normally NewManager from Config file
func NewManagerfromConfig(secBT string, secSTG string, dir string) *manager {
	BT := framework.NewBackTestConfig(dir, "BackTest.yaml", secBT)
	STG := BT.GetStrategy(dir, "BackTest.yaml", secSTG, "Strategy.yaml", "DMT")
	return &manager{
		BT:  &BT,
		STG: STG,
	}
}
func main() {
	var configdirPtr = flag.String("configdir", "./config/Manual/", "a string")
	// * **********************This part is for the Backtesting!**********************
	// * New a manager instance:
	m := NewManagerfromConfig("default", "default", *configdirPtr)
	m.BT.PrepareData("VDS")
	if debug {
		log.Info().Msg("Data Prepared!")
	}
	// * new a strategy from backtest
	pstg := m.STG
	// * new a virtual account
	// ! be careful about the futures part
	va := virtualaccount.NewVirtualAccount(m.BT.BeginDate, m.BT.StockInitValue, m.BT.FuturesInitValue)
	// DCE: debug info
	// ? this is log part! Virtual Account Created!
	if debug {
		log.Info().Str("Account UUID", va.SAcct.UUID).Float64("AccountVal", va.SAcct.MktVal).Msg("Virtual Account Created!")
	}
	// * Iterate the Market data for backtest!
	// for _, instID := range m.BT.SInstrNames {
	// 	fmt.Println(instID)
	//      var new_posmap map[string]*stockaccount.PositionDetail
	// 	 new_posmap
	// 	&va.SAcct.PosMap.append()
	// }
	m.BT.IterData(&va, m.BT.BCM, pstg, m.BT.CPMap, func(in []float64) []float64 { return nil }, "Manual")
}
