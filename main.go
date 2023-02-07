package main

import (
	"encoding/csv"
	"time"

	"os"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/wonderstone/QuantTools/account/virtualaccount"
	"github.com/wonderstone/QuantTools/dataprocessor"
	"github.com/wonderstone/QuantTools/exporter"
	"github.com/wonderstone/QuantTools/indicator"
	"github.com/wonderstone/QuantTools/realinfo"

	"github.com/wonderstone/QuantTools/framework"
	"github.com/wonderstone/QuantTools/strategyModule"
)

const debug = true

// * declare the manager struct used for aggregating the backtest and strategy module
// * backtest has the parameters and the market data
// * strategy interface relates to the strategy module
type manager struct {
	BT  *framework.BackTest
	STG strategyModule.IStrategy
}

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
	// * **********************This part is for the Backtesting!**********************
	// * New a manager instance:
	m := NewManagerfromConfig("Default", "Default", "./config/Manual/")
	// todo: download the data to tmpdata dir first no matter what the data source is
	// // pretend the data has been downloaded already
	// * manager prepares the market data
	m.BT.PrepareData()
	// DCE: debug info
	// ? this is log part! Market Data has been prepared!
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
	m.BT.IterData(&va, m.BT.BCM, pstg, m.BT.CPMap, func(in []float64) []float64 { return nil }, "Manual")
	// * Get the result from virtual stock account and write to the records.csv file
	file, err := os.Create("./records.csv")
	if err != nil {
		// ? Fatal level: this is log part! Error when creating the records.csv file
		log.Fatal().Msg(err.Error())
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()
	// ** Using Write
	for _, record := range va.SAcct.MarketValueSlice {
		row := []string{record.Time, strconv.FormatFloat(record.MktVal, 'f', 2, 64)}
		if err := w.Write(row); err != nil {
			// ? Fatal level: this is log part! Error when writing the records.csv file
			log.Fatal().Msg(err.Error())
		}
	}
	// export realtime.yaml
	exporter.ExportRealtimeYaml("./config/Manual/", "BackTest.yaml", "Default", va)

	// * **********************   The end for the Backtesting!   **********************
	// 注意 这是个偷懒的做法  原则上请只包含一个回测或实盘任务
	// ? 分割线

	// * **********************This part is for the Realtime job!**********************
	// 0. read the config file realtime.yaml and get the vitual account and strategy info
	// !0.1 of course, something should happen to realtime.yaml, make sure it fits the needs
	// 1. get the data from VDS for preload using the same way as Backtest csvprocessor
	// 2. get realtime data from VDS and process it
	// !2.1 frequency check: if the frequency is not the same as strategy required, make it the same
	// !2.2 add indicators to the data from 2.1 as *BarC
	// !2.3 pass the data to channel
	// 3. strategy receives the data from channel and process it
	// * **********************       Realtime job starts!      **********************
	// * 0 从realtime.yaml中读取数据信息
	configdir := "./config/Manual/"
	configfile := "realtime.yaml"
	// * 0.1 从realtime.yaml中读取虚拟账户信息 方便用户调整到自己的持仓 当然也可以手动改yaml文件
	vatmp := virtualaccount.NewVirtualAccountFromConfig(configdir, configfile)
	// * 0.2 从accountinfo.yaml中读取真实账户信息，下单用
	info := realinfo.NewInfoFromConfig("./config/Manual/", "accountinfo.yaml")
	// * 0.3 组合为一个realtime结构体
	rt := framework.NewRealTimeConfig(configdir, "realtime.yaml", info.IM, &vatmp)
	// * 0.4 build a barc channel
	bch := make(chan *dataprocessor.BarC)
	// * 0.5 build a cm channel for futures market
	cmch := make(chan map[string]map[string]float64)

	// * 1. get the data from VDS for preload using the same way as Backtest csvprocessor

	// * 2. get realtime data from VDS and process it

	// * for instance: build a ma2 indicatormap and load some data into the ma2 indicator
	ma2map := make(map[string]*indicator.MA)
	// ** iter the target list
	for _, stock := range rt.SInstrNames {
		ma2map[stock] = indicator.NewMA("Ma2", []int{2}, []string{"Close"})
		// ** data preloading for indicators
		ma2map[stock].LoadData(map[string]float64{"Close": 1.0})
		ma2map[stock].LoadData(map[string]float64{"Close": 2.0})
	}
	// //be serious you jackass!!!

	// * pretend that you finished the data subscribe process, and send data to channel
	go func() {
		for _, dts := range m.BT.BCM.BarCMapkeydts {
			// ** add an indicator to the m.BT.BCM.BarCMap[dts]
			for key, value := range m.BT.BCM.BarCMap[dts].Stockdata {
				ma2map[key].LoadData(value.IndiDataMap)
				// ma2map[key].DQ.Enqueue(value.IndiDataMap["close"])
				m.BT.BCM.BarCMap[dts].Stockdata[key].IndiDataMap["ma2_m"] = ma2map[key].Eval()
			}
			// peek the data, do delete when the test is done
			// for k, v := range m.BT.BCM.BarCMap[dts].Stockdata {
			// 	fmt.Println(k, v.IndiDataMap["ma2_m"])
			// }
			bch <- m.BT.BCM.BarCMap[dts]
			// delay for 0.1 second
			time.Sleep(10 * time.Millisecond)

		}
		// * close the channel
		close(bch)
	}()

	// * 3.0 strategy receives the data from channel and do the realtime job!!!
	// ! be sure you add the code to connect the broker transaction server
	rt.ActOnRTData(configdir, configfile, bch, cmch, pstg, rt.CPMap, func(in []float64) []float64 { return nil }, "Manual")

	// * **********************   The end for the Realtime job!   **********************

}
