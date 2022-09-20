package realinfo

import "github.com/spf13/viper"

type StockInfo struct {
	SUser    string
	SPass    string
	SMQIPAdr string
	SMQPort  string
}

type FuturesInfo struct {
	FUser    string
	FPass    string
	FMQIPAdr string
	FMQPort  string
}

func NewStockInfo(user string, pass string, ipadr string, port string) StockInfo {
	return StockInfo{
		SUser:    user,
		SPass:    pass,
		SMQIPAdr: ipadr,
		SMQPort:  port,
	}
}

func NewFuturesInfo(user string, pass string, ipadr string, port string) FuturesInfo {
	return FuturesInfo{
		FUser:    user,
		FPass:    pass,
		FMQIPAdr: ipadr,
		FMQPort:  port,
	}
}

func NewStockInfoFromConfig(configpath string, filename string) StockInfo {
	viper.SetConfigName(filename)
	viper.AddConfigPath(configpath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return StockInfo{
		SUser:    viper.GetString("StockAccount.Username"),
		SPass:    viper.GetString("StockAccount.Password"),
		SMQIPAdr: viper.GetString("StockMarketQuotesInfo.IPAddr"),
		SMQPort:  viper.GetString("StockMarketQuotesInfo.Port"),
	}
}

func NewFuturesInfoFromConfig(configpath string, filename string) FuturesInfo {
	viper.SetConfigName(filename)
	viper.AddConfigPath(configpath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return FuturesInfo{
		FUser:    viper.GetString("FuturesAccount.Username"),
		FPass:    viper.GetString("FuturesAccount.Password"),
		FMQIPAdr: viper.GetString("FuturesMarketQuotesInfo.IPAddr"),
		FMQPort:  viper.GetString("FuturesMarketQuotesInfo.Port"),
	}
}
