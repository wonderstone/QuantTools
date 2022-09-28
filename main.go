package main

import (
	"encoding/csv"
	// "fmt"
	// "time"

	// "fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/wonderstone/QuantTools/account/virtualaccount"
	// "github.com/wonderstone/QuantTools/dataprocessor"
	"github.com/wonderstone/QuantTools/framework"
	// "github.com/wonderstone/QuantTools/realinfo"
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
	STG := BT.GetStrategy(secSTG, dir, "simple")
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
	// new a strategy from backtest
	pstg := m.STG
	// new virtual account
	va := virtualaccount.NewVirtualAccount(m.BT.BeginDate, m.BT.StockInitValue, m.BT.FuturesInitValue)
	log.Info().Str("Account UUID", va.SAcct.UUID).Float64("AccountVal", va.SAcct.MktVal).Msg("Virtual Account Created!")
	m.BT.IterData(&va, m.BT.BCM, pstg, m.BT.CPMap, func(in []float64) []float64 { return nil }, "Manual")
	file, err := os.Create("./records.csv")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()
	// Using Write
	for _, record := range va.SAcct.MarketValueSlice {
		row := []string{record.Time, strconv.FormatFloat(record.MktVal, 'f', 2, 64)}
		if err := w.Write(row); err != nil {
			log.Fatal().Msg(err.Error())
		}
	}

	// realtime job part

	// // 1.0 从realtime.yaml中读取数据信息
	// configdir := "./"
	// vatmp := virtualaccount.NewVirtualAccountFromConfig(configdir)
	// info := realinfo.NewInfoFromConfig("./config/Manual", "accountinfo")

	// rt := framework.NewRealTimeConfig(configdir, "realtime", info.IM, &vatmp)
	// fmt.Println(rt)

	// // build a barc channel
	// ch := make(chan *dataprocessor.BarC)
	// // pretend that you finished the data subscribe process, and send data to channel
	// go func() {
	// 	for _, dts := range m.BT.BCM.BarCMapkeydts {
	// 		ch <- m.BT.BCM.BarCMap[dts]
	// 		// delay for 0.1 second
	// 		time.Sleep(100 * time.Millisecond)

	// 	}
	// 	// close the channel
	// 	close(ch)
	// }()
	// rt.ActOnRTData(ch, pstg, rt.CPMap, func(in []float64) []float64 { return nil }, "Manual")
	// //

}
