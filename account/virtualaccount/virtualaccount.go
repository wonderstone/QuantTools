package virtualaccount

import (
	"github.com/wonderstone/QuantTools/account"
	"github.com/wonderstone/QuantTools/account/futuresaccount"
	"github.com/wonderstone/QuantTools/account/stockaccount"
	cp "github.com/wonderstone/QuantTools/contractproperty"
)

type VAcct struct {
	SAcct stockaccount.StockAccount
	FAcct futuresaccount.FuturesAccount
}

func NewVirtualAccount(BeginDate string, StockInitValue float64, FuturesInitValue float64) VAcct {
	return VAcct{
		SAcct: stockaccount.NewStockAccount(BeginDate, StockInitValue),
		FAcct: futuresaccount.NewFuturesAccount(BeginDate, FuturesInitValue),
	}
}

func NewVirtualAccountFromConfig(configPath string) VAcct {
	// read config to get cpm
	cpm := cp.NewCPMap("ContractProp", configPath)
	return VAcct{
		SAcct: stockaccount.NewSAFromConfig("realtime", configPath, "VA.sacct", cpm),
		// FAcct: futuresaccount.NewFuturesAccountFromConfig(configPath string),
	}

}

// 未完成！ 因期货夜盘因素，这个地方以何种方式合并存在讨论的必要。
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
