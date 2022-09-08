package contractproperty

import (
	"regexp"

	"github.com/spf13/viper"
)

type CPMap struct {
	StockPropMap   map[string]SCP
	FuturesPropMap map[string]FCP
}

// stock contract property 包含5个类别 但因为与QMT看齐 涨跌幅问题被忽略了。
// ST 			：主板 涨跌5%
// MainBoard	：涨跌10%
// STAR   		：涨跌20% 上海科创板 Sci-Tech innovation board  STAR Market 688XXX st也是20%
// ChiNext  	：涨跌20% 深圳创业板 ChiNext Board 300XXX，创业板st也为20%
// ETF    		：ETF and LOF 没有印花税 SZ159XXX 和 SH510XXX

type SCP struct {
	ContractSize    float64
	TransferFeeRate float64 // 过户费双向收取 历史上2015年 沪深交易所 万分之0.2 现在2022-04-29 沪深万分之0.1
	TaxRate         float64 // 印花税  千分之一  卖方收取
	CommBrokerRate  float64 // 万分之一~万分之五 最低5元 包含交易所收取的经手费0.00487% 证监会最终收取的证管费0.002%(合计0.00687%)
	// Limit           float64 // 涨跌幅度限制 QMT未添加 注释以提高性能
}

// new StockContractProp 策略中会进行调用
func NewSCP(contractsize float64, transferfeerate float64, taxrate float64, commbrokerrate float64) SCP {
	return SCP{
		ContractSize:    contractsize,
		TransferFeeRate: transferfeerate,
		TaxRate:         taxrate,
		CommBrokerRate:  commbrokerrate,
		// Limit:           limit,
	}
}

// viper read config file by input，sample yaml file in config/Manual dir
func NewSCPFromConfig(confName string, sec string, dir string) SCP {
	viper.SetConfigName(confName)
	viper.AddConfigPath(dir)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	tmpMap := viper.GetStringMap("TARGETPROP.STOCK." + sec)
	return NewSCP(tmpMap["contractsize"].(float64), tmpMap["transferfeerate"].(float64), tmpMap["taxrate"].(float64), tmpMap["commbrokerrate"].(float64))
}

// new stockcontractprop from interface{}
func NewSCPFromI(I interface{}) SCP {
	tmpMap := I.(map[string]interface{})
	return SCP{
		ContractSize:    tmpMap["contractsize"].(float64),
		TransferFeeRate: tmpMap["transferfeerate"].(float64),
		TaxRate:         tmpMap["taxrate"].(float64),
		CommBrokerRate:  tmpMap["commbrokerrate"].(float64),
	}
}

// NewSCPFromMap
func SimpleNewSCPFromMap(cpm CPMap, code string) SCP {
	re4char := regexp.MustCompile("[a-zA-Z]*")
	re4num := regexp.MustCompile("[0-9]*$")
	chars := re4char.FindString(code)
	nums := re4num.FindString(code)
	if chars == "sh" && nums[0:3] == "688" {
		return cpm.StockPropMap["star"]
	} else if chars == "sz" && nums[0:3] == "300" {
		return cpm.StockPropMap["chinext"]
	} else if chars == "st" {
		return cpm.StockPropMap["st"]
	} else if chars == "sz" && nums[0:3] == "159" {
		return cpm.StockPropMap["etf"]
	} else if chars == "sh" && nums[0:3] == "510" {
		return cpm.StockPropMap["etf"]
	} else {
		return cpm.StockPropMap["mainboard"]
	}
}

// func ComplexNewSCPFromMap(cpm ContractPropMap, dt string, code string) StockContractProp {
// 	// check dt in some timeseries info map

// 	// check dt finished
// 	re4char := regexp.MustCompile("[a-zA-Z]*")
// 	re4num := regexp.MustCompile("[0-9]*$")
// 	chars := re4char.FindString(code)
// 	nums := re4num.FindString(code)
// 	if chars == "sh" && nums[0:3] == "688" {
// 		return cpm.StockPropMap["star"]
// 	} else if chars == "sz" && nums[0:3] == "300" {
// 		return cpm.StockPropMap["chinext"]
// 	} else if chars == "st" {
// 		return cpm.StockPropMap["st"]
// 	} else if chars == "sz" && nums[0:3] == "159" {
// 		return cpm.StockPropMap["etf"]
// 	} else if chars == "sh" && nums[0:3] == "510" {
// 		return cpm.StockPropMap["etf"]
// 	} else {
// 		return cpm.StockPropMap["mainboard"]
// 	}
// }
// futures contract property

type FCP struct {
	ContractSize      float64
	TickSize          float64
	MarginLong        float64 //保证金
	MarginShort       float64
	MarginBroker      float64
	IsCommRateType    bool
	CommOpen          float64 //手续费
	CommCloseToday    float64
	CommClosePrevious float64
	CommBroker        float64
	// Uplimit           float64 //与股票不同，合约涨跌幅限制较为合约个性化 故不再区分sector
	// Downlimit         float64 //
}

func NewFCP(contractsize float64, ticksize float64, marginlong float64, marginshort float64, marginbroker float64, iscomrate bool, commopen float64, commclosetoday float64, commcloseprevious float64, commbroker float64) FCP {
	return FCP{
		ContractSize:      contractsize,
		TickSize:          ticksize,
		MarginLong:        marginlong,
		MarginShort:       marginshort,
		MarginBroker:      marginbroker,
		IsCommRateType:    iscomrate,
		CommOpen:          commopen,
		CommCloseToday:    commclosetoday,
		CommClosePrevious: commcloseprevious,
		CommBroker:        commbroker,
		// Uplimit:           uplimit,
		// Downlimit:         downlimit,
	}
}

// viper read config file by input
func NewFCPFromConfig(confName string, instrID string, dir string) FCP {
	viper.SetConfigName(confName)
	viper.AddConfigPath(dir)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	tmpMap := viper.GetStringMap("TARGETPROP.FUTURES." + instrID)
	if len(tmpMap) == 0 {
		panic("check config file for instrIDs")
	}
	fcp := NewFCP(tmpMap["contractsize"].(float64), tmpMap["ticksize"].(float64), tmpMap["marginlong"].(float64), tmpMap["marginshort"].(float64), tmpMap["marginbroker"].(float64), tmpMap["iscommratetype"].(bool), tmpMap["commopen"].(float64), tmpMap["commclosetoday"].(float64), tmpMap["commcloseprevious"].(float64), tmpMap["commbroker"].(float64))

	return fcp
}

// NewFSCPFromMap
func SimpleNewFCPFromMap(cpm CPMap, code string) FCP {
	return cpm.FuturesPropMap[code]
}

// NewCPMap from config file
func NewCPMap(confName string, dir string) CPMap {
	viper.SetConfigName(confName)
	viper.AddConfigPath(dir)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	tmpMap := viper.GetStringMap("TARGETPROP")
	if len(tmpMap) == 0 {
		panic("check config file for instrIDs")
	}
	scpm := make(map[string]SCP)
	fcpm := make(map[string]FCP)
	for k, m := range tmpMap {
		switch k {
		case "stock":
			for ks, vs := range m.(map[string]interface{}) {
				vsm := vs.(map[string]interface{})
				scpm[ks] = NewSCP(vsm["contractsize"].(float64), vsm["transferfeerate"].(float64), vsm["taxrate"].(float64), vsm["commbrokerrate"].(float64))
			}
		case "futures":
			for kf, vf := range m.(map[string]interface{}) {
				vfm := vf.(map[string]interface{})
				fcpm[kf] = NewFCP(vfm["contractsize"].(float64), vfm["ticksize"].(float64), vfm["marginlong"].(float64), vfm["marginshort"].(float64), vfm["marginbroker"].(float64), vfm["iscommratetype"].(bool), vfm["commopen"].(float64), vfm["commclosetoday"].(float64), vfm["commcloseprevious"].(float64), vfm["commbroker"].(float64))
			}
		}
	}

	return CPMap{
		StockPropMap:   scpm,
		FuturesPropMap: fcpm,
	}
}
