package main

import (
	"github.com/rs/zerolog/log"
	"github.com/wonderstone/QuantTools/framework"
	"github.com/wonderstone/QuantTools/strategyModule"
)

// declare the manager struct used for aggregating the backtest and strategy module
// backtest has the parameters and the market data
// strategy interface relates to the strategy module
type manager struct {
	BT  *framework.BackTest      // BackTest framework component
	STG strategyModule.IStrategy // 在一个BackTest framework下  有多个策略实例，每个策略实例都对应着不一样的GEP表达式
}

// NewManager creates a new manager instance
func NewManagerfromConfig(secBT string, secSTG string, dir string) *manager {
	BT := framework.NewBackTestConfig(secBT, dir)
	STG := BT.GetStrategy(secSTG, dir)
	return &manager{
		BT:  &BT,
		STG: STG,
	}
}

func main() {
	// create a manager instance:
	m := NewManagerfromConfig("Default", "Default", "./config/Manual")
	// manager prepares the market data
	m.BT.PrepareData()
	log.Info().Msg("Data Prepared!")

}
