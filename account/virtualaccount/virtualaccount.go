package virtualaccount

import (
	"github.com/wonderstone/QuantTools/account"
	"github.com/wonderstone/QuantTools/account/futuresaccount"
	"github.com/wonderstone/QuantTools/account/stockaccount"
	cp "github.com/wonderstone/QuantTools/contractproperty"
)

// * VAcct is the composite of stock and futures account
type VAcct struct {
	SAcct stockaccount.StockAccount
	FAcct futuresaccount.FuturesAccount
}

// * Normally new one
func NewVirtualAccount(BeginDate string, StockInitValue float64, FuturesInitValue float64) VAcct {
	return VAcct{
		SAcct: stockaccount.NewStockAccount(BeginDate, StockInitValue),
		FAcct: futuresaccount.NewFuturesAccount(BeginDate, FuturesInitValue),
	}
}

// * Normally new from config file
func NewVirtualAccountFromConfig(configDir string, configFile string) VAcct {
	// read config to get cpm
	cpm := cp.NewCPMap("./config/Manual/", "ContractProp.yaml")
	return VAcct{
		SAcct: stockaccount.NewSAFromConfig(configDir, configFile, "va.sacct", cpm),
		// ! 正式发布前，请务必检查此处！ 还没细看对期货的支持呢！
		// FAcct: futuresaccount.NewFuturesAccountFromConfig(configPath string),
	}

}

// todo 未完成！ 因期货夜盘因素，这个地方以何种方式合并存在讨论的必要。
func (va *VAcct) SumVAcctMktVal() []account.MktValDataType {
	SumVAcctMV := make([]account.MktValDataType, len(va.SAcct.MarketValueSlice))
	for i := range va.SAcct.MarketValueSlice {
		SumVAcctMV[i].Time = va.SAcct.MarketValueSlice[i].Time
		if len(va.SAcct.MarketValueSlice) == 0 && len(va.FAcct.MarketValueSlice) == 0 {
			panic("确保账户有数据")
		} else if len(va.FAcct.MarketValueSlice) == 0 && len(va.FAcct.MarketValueSlice) != 0 {
			SumVAcctMV[i].MktVal = va.FAcct.MarketValueSlice[i].MktVal
		} else if len(va.FAcct.MarketValueSlice) != 0 && len(va.FAcct.MarketValueSlice) == 0 {
			SumVAcctMV[i].MktVal = va.SAcct.MarketValueSlice[i].MktVal
		} else {
			SumVAcctMV[i].MktVal = va.SAcct.MarketValueSlice[i].MktVal + va.FAcct.MarketValueSlice[i].MktVal
		}
	}
	return SumVAcctMV
}
