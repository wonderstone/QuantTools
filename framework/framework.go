package framework

import (
	"github.com/wonderstone/QuantTools/account"
	"github.com/wonderstone/QuantTools/account/virtualaccount"
	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/exporter"

	"io/ioutil"

	"github.com/wonderstone/QuantTools/dataprocessor"
	"github.com/wonderstone/QuantTools/matcher"

	"github.com/wonderstone/QuantTools/perfeval"
	"github.com/wonderstone/QuantTools/strategyModule"

	"github.com/rs/zerolog/log"

	"sort"
	"strings"
	// "sync"
	"time"

	"github.com/spf13/viper"
)

// 0. steps: init the backtest struct -> PrepareData -> IterData
type BackTest struct {
	// 账户初始化参数
	StockInitValue   float64
	FuturesInitValue float64
	// Section for data range
	BeginDate string
	EndDate   string
	// Section for Strategy Targets and info fields
	SInstrNames    []string
	SIndiNames     []string
	SCsvDatafields []string
	FInstrNames    []string
	FIndiNames     []string
	FCsvDatafields []string
	// Section for CSV data dir
	StockDataDir      string
	FuturesDataDir    string
	FuturesMTMDataDir string
	// Section for ContractProp
	ConfName  string
	CPDataDir string
	// Section for Matcher
	MatcherSlippage4S float64
	MatcherSlippage4F float64
	// Performance Analytics Parameter
	RiskFreeRate float64
	PAType       string
	// Section for Strategy Module Selection
	StrategyMod string
	SMGEPType   string
	SMName      string
	SMDataDir   string
	// market data
	BCM   *dataprocessor.BarCM
	CPMap cp.CPMap

	// this part is only for test with zerolog and structuring the log
	// fileLogger zerolog.Logger

	// add a sync.RWMutex to make sure BackTest
	// sync.RWMutex
}

type RealTime struct {
	// 实盘任务所需真实信息 IP Port user password等
	Info map[string]interface{}
	// Virtual Account
	VA *virtualaccount.VAcct
	// Section for Strategy Targets and info fields
	SInstrNames []string
	SIndiNames  []string
	SRegisterDF []string
	FInstrNames []string
	FIndiNames  []string
	FRegisterDF []string
	// Section for ContractProp
	ConfName  string
	CPDataDir string
	CPMap     cp.CPMap
	// Section for Matcher
	MatcherSlippage4S float64
	MatcherSlippage4F float64
	// Section for Strategy Module Selection
	StrategyMod string
	SMGEPType   string
	SMName      string
	SMDataDir   string
}

func NewBackTest(SInitVal float64, FInitVale float64, BDt string, EDt string,
	SInstrNs []string, SIndiNs []string, SCDtfields []string, FInstrNs []string, FIndiNs []string, FCDtfields []string,
	SDtDir string, FDtDir string, FMTMDtDir string, ConfName string, CPDataDir string, MatcherSlpg4S float64, MatcherSlpg4F float64,
	StrategyMod string, SMGEPType string, SMName string, SMDataDir string, RiskFR float64, PAType string) BackTest {
	return BackTest{
		// 所有项目均为用户设置
		// 账户初始化参数，用户资金
		StockInitValue:   SInitVal,
		FuturesInitValue: FInitVale,
		// Section for data range
		BeginDate: BDt,
		EndDate:   EDt,
		// Section for Strategy Targets and info fields
		SInstrNames:    SInstrNs,
		SIndiNames:     SIndiNs,
		SCsvDatafields: SCDtfields,
		FInstrNames:    FInstrNs,
		FIndiNames:     FIndiNs,
		FCsvDatafields: FCDtfields,
		// Section for CSV data dir
		StockDataDir:      SDtDir,
		FuturesDataDir:    FDtDir,
		FuturesMTMDataDir: FMTMDtDir,
		// Section for ContractProp
		ConfName:  ConfName,
		CPDataDir: CPDataDir,
		// Section for Matcher
		MatcherSlippage4S: MatcherSlpg4S,
		MatcherSlippage4F: MatcherSlpg4F,
		// Section for Strategy Module Selection
		StrategyMod: StrategyMod,
		SMGEPType:   SMGEPType,
		SMName:      SMName,
		SMDataDir:   SMDataDir,
		// Performance Analytics Parameter
		RiskFreeRate: RiskFR,
		PAType:       PAType,
	}
}

func NewBackTestConfig(sec string, dir string) BackTest {
	viper.SetConfigName("BackTest")
	viper.AddConfigPath(dir)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	tmpMap := viper.GetStringMap(sec)
	var sinstrnames []string
	for _, v := range tmpMap["sinstrnames"].([]interface{}) {
		sinstrnames = append(sinstrnames, v.(string))
	}
	var sindinames []string
	for _, v := range tmpMap["sindinames"].([]interface{}) {
		sindinames = append(sindinames, v.(string))
	}
	var scsvdatafields []string
	for _, v := range tmpMap["scsvdatafields"].([]interface{}) {
		scsvdatafields = append(scsvdatafields, v.(string))
	}
	var finstrnames []string
	for _, v := range tmpMap["finstrnames"].([]interface{}) {
		finstrnames = append(finstrnames, v.(string))
	}
	var findinames []string
	for _, v := range tmpMap["findinames"].([]interface{}) {
		findinames = append(findinames, v.(string))
	}
	var fcsvdatafields []string
	for _, v := range tmpMap["fcsvdatafields"].([]interface{}) {
		fcsvdatafields = append(fcsvdatafields, v.(string))
	}
	return NewBackTest(tmpMap["stockinitvalue"].(float64), tmpMap["futuresinitvalue"].(float64),
		tmpMap["begindate"].(string), tmpMap["enddate"].(string),
		sinstrnames, sindinames, scsvdatafields, finstrnames, findinames, fcsvdatafields,
		tmpMap["stockdatadir"].(string), tmpMap["futuresdatadir"].(string), tmpMap["futuresmtmdatadir"].(string),
		tmpMap["confname"].(string), tmpMap["cpdatadir"].(string), tmpMap["matcherslippage4s"].(float64), tmpMap["matcherslippage4f"].(float64),
		tmpMap["strategymodule"].(string), tmpMap["smgeptype"].(string), tmpMap["smname"].(string), tmpMap["smdatadir"].(string),
		tmpMap["riskfreerate"].(float64), tmpMap["patype"].(string))
}

func NewRealTime(info map[string]interface{}, va *virtualaccount.VAcct, SInstrNs []string, SIndiNs []string, SRDtfields []string, FInstrNs []string, FIndiNs []string, FRDtfields []string,
	ConfName string, CPDataDir string, cpm cp.CPMap, MatcherSlpg4S float64, MatcherSlpg4F float64,
	StrategyMod string, SMGEPType string, SMName string, SMDataDir string) RealTime {
	return RealTime{
		// 所有项目均为用户设置
		Info: info,
		// 账户初始化参数，用户资金
		VA: va,
		// Section for Strategy Targets and info fields
		SInstrNames: SInstrNs,
		SIndiNames:  SIndiNs,
		SRegisterDF: SRDtfields,
		FInstrNames: FInstrNs,
		FIndiNames:  FIndiNs,
		FRegisterDF: FRDtfields,
		// Section for ContractProp
		ConfName:  ConfName,
		CPDataDir: CPDataDir,
		CPMap:     cpm,
		// Section for Matcher
		MatcherSlippage4S: MatcherSlpg4S,
		MatcherSlippage4F: MatcherSlpg4F,
		// Section for Strategy Module Selection
		StrategyMod: StrategyMod,
		SMGEPType:   SMGEPType,
		SMName:      SMName,
		SMDataDir:   SMDataDir,
	}
}

// NewRealTimeConfig 从配置文件中读取配置信息 filename could be realtime
func NewRealTimeConfig(dir string, filename string, info map[string]interface{}, va *virtualaccount.VAcct) RealTime {
	viper.SetConfigName(filename)
	viper.AddConfigPath(dir)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	SInstrNs := viper.GetStringSlice("DataFields.sinstrnames")
	SIndiNs := viper.GetStringSlice("AFields.SIndiNmsAfter")
	SRDtfields := viper.GetStringSlice("DataFields.scsvdatafields")
	FInstrNs := viper.GetStringSlice("DataFields.finstrnames")
	FIndiNs := viper.GetStringSlice("DataFields.findinames")
	FRDtfields := viper.GetStringSlice("DataFields.fcsvdatafields")
	ConfName := viper.GetString("ContractProp.ConfName")
	CPDataDir := viper.GetString("ContractProp.CPDataDir")

	MatcherSlpg4S := viper.GetFloat64("MatcherParam.MatcherSlippage4S")
	MatcherSlpg4F := viper.GetFloat64("MatcherParam.MatcherSlippage4F")
	StrategyMod := viper.GetString("StgModel.StrategyModule")
	SMGEPType := viper.GetString("StgModel.SMGEPType")
	SMName := viper.GetString("StgModel.SMName")
	SMDataDir := viper.GetString("StgModel.SMDataDir")
	cpm := cp.NewCPMap("ContractProp", dir)
	return NewRealTime(info, va, SInstrNs, SIndiNs, SRDtfields, FInstrNs, FIndiNs, FRDtfields, ConfName, CPDataDir, cpm,
		MatcherSlpg4S, MatcherSlpg4F, StrategyMod, SMGEPType, SMName, SMDataDir)
}

type void struct{}

// 读取文件夹内文件 存在判断是否包含元素的功能屡次调用 不建议使用slice 这里采用map模拟set
func getFileMap(path string) map[string]void {
	res := make(map[string]void)
	var member void
	fs, _ := ioutil.ReadDir(path)
	for _, file := range fs {
		if !file.IsDir() {
			res[strings.TrimSuffix(file.Name(), ".csv")] = member
		}
	}
	return res
}

// 0. 输出strategy
func (BT *BackTest) GetStrategy(sec string, dir string, tag string) strategyModule.IStrategy {
	// switch BT.StrategyMod {
	// case "Simple":
	// 	return strategyModule.NewStrategyFromConfig(sec, dir)
	// default:
	// 	return strategyModule.NewStrategyFromConfig(sec, dir)
	// }

	return strategyModule.GetStrategy(sec, dir, tag)
}

// 1. 准备数据
func (BT *BackTest) PrepareData() {
	sfilemap := getFileMap(BT.StockDataDir)
	if len(sfilemap) != len(BT.SInstrNames) && len(BT.SInstrNames) != 0 {
		panic("股票操作标的数与数据文件个数不匹配")
	}
	ffilemap := getFileMap(BT.FuturesDataDir)
	if len(ffilemap) != len(BT.FInstrNames) && len(BT.FInstrNames) != 0 {
		panic("期货操作标的数与数据文件个数不匹配")
	}
	// 读取文件 准备数据
	BT.BCM = dataprocessor.NewBarCM(BT.SInstrNames, BT.SIndiNames, BT.FInstrNames, BT.FIndiNames, BT.BeginDate, BT.EndDate)
	// 读取股票数据
	if len(BT.SInstrNames) != 0 {
		Sfiles, err := dataprocessor.ListDir(BT.StockDataDir, "csv")
		if err != nil {
			panic(err)
		}
		for _, Sfile := range Sfiles {
			BT.BCM.CsvSBarReader(Sfile)
		}
	}

	// 读取期货数据
	if len(BT.FInstrNames) != 0 {
		Ffiles, err := dataprocessor.ListDir(BT.FuturesDataDir, "csv")
		if err != nil {
			panic(err)
		}
		for _, FPfile := range Ffiles {
			BT.BCM.CsvFBarReader(FPfile)
		}
		// 读取期货结算数据
		MTMfiles, err := dataprocessor.ListDir(BT.FuturesMTMDataDir, "csv")
		if err != nil {
			panic(err)
		}
		for _, MTMfile := range MTMfiles {
			BT.BCM.CsvFMTMReader(MTMfile)
		}
	}

	// 产生合约属性 Map
	BT.CPMap = cp.NewCPMap(BT.ConfName, BT.CPDataDir)
	// 生成升序时间index
	for mapkeydt := range BT.BCM.BarCMap {
		BT.BCM.BarCMapkeydts = append(BT.BCM.BarCMapkeydts, mapkeydt)
	}
	sort.Slice(BT.BCM.BarCMapkeydts, func(i, j int) bool {
		dti, _ := time.Parse("2006/1/2 15:04", BT.BCM.BarCMapkeydts[i])
		dtj, _ := time.Parse("2006/1/2 15:04", BT.BCM.BarCMapkeydts[j])
		return dti.Before(dtj)
	})

	// this part is for test only with zerolog
	// tmpFile, err := ioutil.TempFile(os.TempDir(), "zerolog_framework")
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Fail to create tmp file")
	// }
	// BT.fileLogger = zerolog.New(tmpFile).With().Timestamp().Logger()
	// fmt.Printf("The log file is allocated at %s\n", tmpFile.Name())

}

// 2. 遍历数据
// VAcct 引用(指针)传递 ，BCM 引用(指针)传递，strategymodule值传递，CPMap值传递，Eval函数传递(回调)
func (BT *BackTest) IterData(VAcct *virtualaccount.VAcct, BCM *dataprocessor.BarCM, strategymodule strategyModule.IStrategy, CPMap cp.CPMap, Eval func([]float64) []float64, mode string) {
	// when SInstrNames and FInstrNames are empty, panic
	if len(BT.FInstrNames) == 0 && len(BT.SInstrNames) == 0 {
		panic("没有操作标的")
	}

	var lastdatetime string
	// get a matcher and a temp orderResult
	simplematcher := matcher.NewSimpleMatcher(BT.MatcherSlippage4S, BT.MatcherSlippage4F)
	tmpOrderRes := strategyModule.NewOrderResult()

	//循环更新数据
	for _, mapkeydt := range BCM.BarCMapkeydts {
		// fmt.Println("MktVal:", VAcct.SAcct.MktVal)

		//  2.3 循环股票和期货的orderslice 基于bar的open撮合
		for i := range tmpOrderRes.StockOrderS {
			// 验证数据是否存在,存在时才撮合
			if matchinfo, isOk := BCM.BarCMap[mapkeydt].Stockdata[tmpOrderRes.StockOrderS[i].InstID]; isOk {
				if !strategyModule.ContainNaN(matchinfo.IndiDataMap) {
					// this part is for test only
					log.Info().Str("Account UUID", VAcct.SAcct.UUID).Str("TimeStamp", mapkeydt).
						Str("Target", tmpOrderRes.StockOrderS[i].InstID).Float64("MatchPrice", matchinfo.IndiDataMap["open"]).
						Msg("Match details")
					VAcct.SAcct.CheckEligible(&tmpOrderRes.StockOrderS[i])
					simplematcher.MatchStockOrder(&tmpOrderRes.StockOrderS[i], matchinfo.IndiDataMap["open"], mapkeydt)
					//
					// tmpOrderRes.IsExecuted = true
					VAcct.SAcct.ActOnOrder(&tmpOrderRes.StockOrderS[i])
					// this part is for test only
					log.Info().Str("Account UUID", VAcct.SAcct.UUID).Str("TimeStamp", mapkeydt).Msg("Stock Order Executed")

				}

			}
		}
		for i := range tmpOrderRes.FuturesOrderS {
			// 验证数据是否存在,存在时才撮合
			if matchinfo, isOk := BCM.BarCMap[mapkeydt].Futuresdata[tmpOrderRes.FuturesOrderS[i].InstID]; isOk {
				if !strategyModule.ContainNaN(matchinfo.IndiDataMap) {
					VAcct.FAcct.CheckEligible(&tmpOrderRes.FuturesOrderS[i])
					simplematcher.MatchFuturesOrder(&tmpOrderRes.FuturesOrderS[i], matchinfo.IndiDataMap["open"], mapkeydt)
					// tmpOrderRes.IsExecuted = true
					VAcct.FAcct.ActOnOrder(&tmpOrderRes.FuturesOrderS[i])

				}
			}
		}

		//  循环股票和期货的orderslice 账户对order进行更新,原则上可以与上面合并 已经合并  测试一下
		// for i := range tmpOrderRes.StockOrderS {
		// 	if tmpOrderRes.StockOrderS[i].IsExecuted {
		// 		VAcct.SAcct.ActOnOrder(&tmpOrderRes.StockOrderS[i])
		// 		// this part is for test only
		// 		log.Info().Str("Account UUID", VAcct.SAcct.UUID).Str("TimeStamp", mapkeydt).Msg("Stock Order Executed")
		// 	}
		// }
		// for i := range tmpOrderRes.FuturesOrderS {
		// 	if tmpOrderRes.FuturesOrderS[i].IsExecuted {
		// 		VAcct.FAcct.ActOnOrder(&tmpOrderRes.FuturesOrderS[i])
		// 	}
		// }

		//2.0 判断是否符合close或MTM条件 确认是否需收盘
		if lastdatetime != "" {
			if len(BCM.BarCMap[mapkeydt].Stockdata) != 0 && strings.Fields(lastdatetime)[0] != strings.Fields(mapkeydt)[0] {
				//2.0.1 如果符合 账户进行对应操作
				VAcct.SAcct.ActOnCM()
				// this part is for test only
				log.Info().Str("Account UUID", VAcct.SAcct.UUID).Str("TimeStamp", mapkeydt).Msg("Market Close")
			}
			// 期货这个需要留意一下具体情况
			if len(BCM.BarCMap[mapkeydt].Futuresdata) != 0 && (strings.Fields(mapkeydt)[1] > "15:15" && strings.Fields(lastdatetime)[1] <= "15:15") {
				for _, instrname := range BT.FInstrNames {
					if v, ok := BCM.FMTMDataMap[mapkeydt][instrname]; ok {

						VAcct.FAcct.ActOnMTM(mapkeydt, instrname, v)
					} else {
						panic("MTM数据与Bar数据不匹配")
					}
				}
			}
		}

		//  2.1 账户接收数据刷新
		if len(BCM.BarCMap[mapkeydt].Stockdata) != 0 {
			for instID, barC := range BCM.BarCMap[mapkeydt].Stockdata {
				if !strategyModule.ContainNaN(barC.IndiDataMap) {
					VAcct.SAcct.ActOnUpdateMI(mapkeydt, instID, barC.IndiDataMap["close"])
					// this part is for test only
					log.Info().Str("Account UUID", VAcct.SAcct.UUID).Str("TimeStamp", mapkeydt).
						Float64("AccountVal", VAcct.SAcct.MktVal).Float64("close", barC.IndiDataMap["close"]).
						Float64("open", barC.IndiDataMap["open"]).Float64("high", barC.IndiDataMap["high"]).
						Float64("vol", barC.IndiDataMap["vol"]).Float64("ma1", barC.IndiDataMap["ma1"]).
						Str("Target", instID).
						Msg("Data")
					// if instID is in PosMap then log
					if _, ok := VAcct.SAcct.PosMap[instID]; ok {
						log.Info().Str("Account UUID", VAcct.SAcct.UUID).Str("TimeStamp", mapkeydt).Str("target", instID).
							Float64("positdy", VAcct.SAcct.PosMap[instID].CalPosTdyNum()).
							Float64("posipre", VAcct.SAcct.PosMap[instID].CalPosPrevNum()).
							Float64("Equity", VAcct.SAcct.PosMap[instID].CalEquity()).
							Float64("UnRealProfit", VAcct.SAcct.PosMap[instID].CalUnRealizedProfit()).
							Float64("AllCommission", VAcct.SAcct.AllCommission).Float64("AllProfit", VAcct.SAcct.AllProfit).
							Float64("Fundavail", VAcct.SAcct.Fundavail).Float64("Equity4ALL", VAcct.SAcct.Equity()).
							Msg("Account")
					}
				}
			}
		}
		if len(BCM.BarCMap[mapkeydt].Futuresdata) != 0 {
			for instID, barC := range BCM.BarCMap[mapkeydt].Futuresdata {
				if !strategyModule.ContainNaN(barC.IndiDataMap) {
					VAcct.FAcct.ActOnUpdateMI(mapkeydt, instID, barC.IndiDataMap["close"])
				}
			}
		}
		//  2.2 策略接收数据并经过ActOnData得到对应账户的orderslice
		switch mode {
		case "GEP":
			tmpOrderRes = strategymodule.ActOnData(mapkeydt, BCM.BarCMap[mapkeydt], VAcct, CPMap, Eval)
		case "Manual":
			tmpOrderRes = strategymodule.ActOnDataMAN(mapkeydt, BCM.BarCMap[mapkeydt], VAcct, CPMap)
		default:
			panic("mode is not defined")
		}

		// this part is for test only
		log.Info().Str("Account UUID", VAcct.SAcct.UUID).Str("TimeStamp", mapkeydt).Msg("Strategy ActOnData Finished")
		// 临时看一下，记得删除
		// fmt.Println("mapkeydt:", mapkeydt, "lastdatetime:", lastdatetime)
		lastdatetime = mapkeydt

	}

}

func (BT *BackTest) EvalPerformance(MarketValueSlice []account.MktValDataType, einfo map[string]interface{}) float64 {
	//  4.0 获得账户的mkvslice 进行评估
	// new a performanceevaluator
	// BT.Lock()
	// defer BT.Unlock()
	PE := perfeval.NewPerfEval()
	PE.MktValSlice = MarketValueSlice
	return PE.CalcPerfEvalResult(einfo)
}

func (RT *RealTime) ActOnRTData(bc <-chan *dataprocessor.BarC, mc <-chan map[string]map[string]float64, strategymodule strategyModule.IStrategy, CPMap cp.CPMap, Eval func([]float64) []float64, mode string) {
	// 期货账户开启goroutine 用于接收mc数据 并更新账户
	go RT.ActOnCM(mc)
	// 2.0 defer 将va数据更新写入到realtime.yaml中
	defer exporter.ReplaceVA("../config/Manual", "realtime", *RT.VA)

	// 3.0 dataprocessor中RealTimeProcess
	var lastdatetime string
	// get a matcher and a temp orderResult
	simplematcher := matcher.NewSimpleMatcher(RT.MatcherSlippage4S, RT.MatcherSlippage4F)
	tmpOrderRes := strategyModule.NewOrderResult()
	// 循环从BarC channel中读取数据 直到channel关闭，获取的data为BarC类型
	for data := range bc {
		// ts:=getRealTimeStamp()
		timestamp, err := data.GetTimeStamp()
		//  2.3 循环股票和期货的orderslice 基于bar的open撮合
		for i := range tmpOrderRes.StockOrderS {
			// 验证数据是否存在,存在时才撮合
			if matchinfo, isOk := data.Stockdata[tmpOrderRes.StockOrderS[i].InstID]; isOk {
				if !strategyModule.ContainNaN(matchinfo.IndiDataMap) {
					// this part is for test only
					log.Info().Str("Account UUID", RT.VA.SAcct.UUID).Str("TimeStamp", lastdatetime).
						Str("Target", tmpOrderRes.StockOrderS[i].InstID).Float64("MatchPrice", matchinfo.IndiDataMap["open"]).
						Msg("Match details")
					// 采用本bar的open价格进行撮合
					RT.VA.SAcct.CheckEligible(&tmpOrderRes.StockOrderS[i])
					simplematcher.MatchStockOrder(&tmpOrderRes.StockOrderS[i], matchinfo.IndiDataMap["open"], lastdatetime)
					RT.VA.SAcct.ActOnOrder(&tmpOrderRes.StockOrderS[i])

					// this part is for test only
					log.Info().Str("Account UUID", RT.VA.SAcct.UUID).Str("TimeStamp", lastdatetime).Msg("Stock Order Executed")

				}

			}
		}

		for i := range tmpOrderRes.FuturesOrderS {
			// 验证数据是否存在,存在时才撮合
			if matchinfo, isOk := data.Futuresdata[tmpOrderRes.FuturesOrderS[i].InstID]; isOk {
				if !strategyModule.ContainNaN(matchinfo.IndiDataMap) {
					// 采用本bar的open价格进行撮合
					RT.VA.FAcct.CheckEligible(&tmpOrderRes.FuturesOrderS[i])
					simplematcher.MatchFuturesOrder(&tmpOrderRes.FuturesOrderS[i], matchinfo.IndiDataMap["open"], lastdatetime)
					RT.VA.FAcct.ActOnOrder(&tmpOrderRes.FuturesOrderS[i])
					// in case you wanna put some log here!
				}
			}
		}
		//2.0 判断是否符合close或MTM条件 确认是否需收盘
		if lastdatetime != "" {
			// 股票收盘
			if len(data.Stockdata) != 0 && strings.Fields(lastdatetime)[0] != strings.Fields(timestamp)[0] {
				//2.0.1 如果符合 账户进行对应操作
				RT.VA.SAcct.ActOnCM()
				// this part is for test only
				log.Info().Str("Account UUID", RT.VA.SAcct.UUID).Str("TimeStamp", timestamp).Msg("Market Close")
			}
			// ! 期货实盘应该以接收到收盘数据为准，这个过程应该由其他goroutine完成
			// 期货收盘

		}
		//  2.1 账户接收数据刷新
		if len(data.Stockdata) != 0 {
			for instID, barC := range data.Stockdata {
				if !strategyModule.ContainNaN(barC.IndiDataMap) {
					RT.VA.SAcct.ActOnUpdateMI(timestamp, instID, barC.IndiDataMap["close"])
					// this part is for test only
					log.Info().Str("Account UUID", RT.VA.SAcct.UUID).Str("TimeStamp", timestamp).
						Float64("AccountVal", RT.VA.SAcct.MktVal).Float64("close", barC.IndiDataMap["close"]).
						Float64("open", barC.IndiDataMap["open"]).Float64("high", barC.IndiDataMap["high"]).
						Float64("vol", barC.IndiDataMap["vol"]).Float64("ma1", barC.IndiDataMap["ma1"]).
						Str("Target", instID).
						Msg("Data")
					// if instID is in PosMap then log
					if _, ok := RT.VA.SAcct.PosMap[instID]; ok {
						log.Info().Str("Account UUID", RT.VA.SAcct.UUID).Str("TimeStamp", timestamp).Str("target", instID).
							Float64("positdy", RT.VA.SAcct.PosMap[instID].CalPosTdyNum()).
							Float64("posipre", RT.VA.SAcct.PosMap[instID].CalPosPrevNum()).
							Float64("Equity", RT.VA.SAcct.PosMap[instID].CalEquity()).
							Float64("UnRealProfit", RT.VA.SAcct.PosMap[instID].CalUnRealizedProfit()).
							Float64("AllCommission", RT.VA.SAcct.AllCommission).Float64("AllProfit", RT.VA.SAcct.AllProfit).
							Float64("Fundavail", RT.VA.SAcct.Fundavail).Float64("Equity4ALL", RT.VA.SAcct.Equity()).
							Msg("Account")
					}
				}
			}
		}
		if len(data.Futuresdata) != 0 {
			for instID, barC := range data.Futuresdata {
				if !strategyModule.ContainNaN(barC.IndiDataMap) {
					RT.VA.FAcct.ActOnUpdateMI(timestamp, instID, barC.IndiDataMap["close"])
				}
			}
		}
		//  2.2 策略接收数据并经过ActOnData得到对应账户的orderslice
		if err != nil {
			switch mode {
			case "GEP":
				tmpOrderRes = strategymodule.ActOnData(timestamp, data, RT.VA, CPMap, Eval)
				SendOrders(RT.Info, tmpOrderRes)
			case "Manual":
				tmpOrderRes = strategymodule.ActOnDataMAN(timestamp, data, RT.VA, CPMap)
				SendOrders(RT.Info, tmpOrderRes)
			default:
				panic("mode is not defined")
			}
		}

		// this part is for test only
		log.Info().Str("Account UUID", RT.VA.SAcct.UUID).Str("TimeStamp", timestamp).Msg("Strategy ActOnData Finished")
		// 临时看一下，记得删除
		// fmt.Println("mapkeydt:", mapkeydt, "lastdatetime:", lastdatetime)
		lastdatetime = timestamp

	}

}

func (RT *RealTime) ActOnCM(mc <-chan map[string]map[string]float64) {
	// 1. get data from mc and va update with the data
	for data := range mc {
		for timestamp, kv := range data {
			for k, v := range kv {
				RT.VA.FAcct.ActOnMTM(timestamp, k, v)
			}
		}
	}
}
