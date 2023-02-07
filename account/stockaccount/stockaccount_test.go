// 核心测试两类操作：基于order 基于data
// order的情况需要考虑：1. 持仓记录添加是否正常；2. 隔日更新进入昨日持仓是否正常；
//  3. 卖出当日持仓(T+1)是否被拒绝；4. 卖出超过持仓是否报错；
//
// data 的情况需要考虑：1. 提供非持仓数据是否忽略；2. 提供持仓数据是否对应更新；
//  3. 提供新一日数据，当日持仓是否转入前持仓
package stockaccount

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/order"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestNewStockAccount(t *testing.T) {
	// 注意，这里使用testify的assert进行对象比对，如果是指针或切片，比对的是具体的对象而非地址
	// uuid 过不了
	cash := 100000.0
	initTime := "20220515 13:35:27 500"
	res := NewStockAccount(initTime, cash)

	fmt.Printf("Address of res:\t%p\n", &res)
	expected := StockAccount{
		InitTime:  "20220515 13:35:27 500",
		UdTime:    "20220515 13:35:27 500",
		MktVal:    cash,
		Fundavail: cash,
		PosMap:    make(map[string]*PositionSlice),
	}
	fmt.Printf("Address of expected:\t%p\n", &expected)
	// file output for zerolog
	tmpFile, err := ioutil.TempFile(os.TempDir(), "zerolog_test")
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to create tmp file")
	}
	fileLgger := zerolog.New(tmpFile).With().Timestamp().Logger()
	fileLgger.Info().Str("Account UUID", res.UUID).Msg("NewStockAccount")

	fmt.Printf("The log file is allocated at %s\n", tmpFile.Name())
	assert.Equal(t, &expected, &res, "stockaccount创建不符合预期")
}

// 1. 接收Day1一个targetA的order  查看是否增加在map的今日切片里 以及所有字段
func TestActOnOrder(t *testing.T) {
	cash := 1000000.0
	initTime := "20220515 13:25:00 500"
	account := NewStockAccount(initTime, cash)

	// expected = expected.PosMap
	// 一个ContractProp给单子提供计算信息

	var orderscp cp.SCP = cp.SCP{ContractSize: 100, TransferFeeRate: 0.00001, TaxRate: 0.001, CommBrokerRate: 0.0002}
	// 下一单新 检查
	tmporder0 := order.NewStockOrder("600000", true, true, "20220515 13:35:27 500", 8.04, 2, "Buy", &orderscp)
	account.ActOnOrder(&tmporder0)
	// account结构过于复杂，仅核对部分重点字段
	assert.Equal(t, "20220515 13:25:00 500", account.InitTime, "账户初始化时间更新不符合预期")
	assert.Equal(t, "20220515 13:35:27 500", account.UdTime, "账户刷新时间更新不符合预期")
	assert.Equal(t, 998386.98392, account.Fundavail, "账户可用资金更新不符合预期")
	assert.Equal(t, 0.0, account.AllProfit, "账户总利润更新不符合预期")
	assert.Equal(t, 5.01608, account.AllCommission, "账户总佣金更新不符合预期")
	assert.Equal(t, 1607.9999999999998, account.Equity(), "账户总股权更新不符合预期")
	assert.Equal(t, 999994.98392, account.MktVal, "账户总市值更新不符合预期")
	assert.Equal(t, 2.0, account.PosMap["600000"].CalPosTdyNum(), "账户持股数更新不符合预期")

	account.ActOnCM()

	tmporder1 := order.NewStockOrder("600000", true, true, "20220516 13:40:27 500", 8.00, 4, "Buy", &orderscp)
	account.ActOnOrder(&tmporder1)
	assert.Equal(t, "20220515 13:25:00 500", account.InitTime, "账户初始化时间更新不符合预期")
	assert.Equal(t, "20220516 13:40:27 500", account.UdTime, "账户刷新时间更新不符合预期")
	assert.Equal(t, 995181.95192, account.Fundavail, "账户可用资金更新不符合预期")
	assert.Equal(t, 0.0, account.AllProfit, "账户总利润更新不符合预期")
	assert.Equal(t, 10.048079999999999, account.AllCommission, "账户总佣金更新不符合预期")
	assert.Equal(t, 4800.0, account.Equity(), "账户总股权更新不符合预期")
	assert.Equal(t, 4.0, account.PosMap["600000"].CalPosTdyNum(), "账户持股数更新不符合预期")

	tmporder2 := order.NewStockOrder("600000", true, true, "20220516 13:45:27 500", 9.0, 1, "Sell", &orderscp)
	account.ActOnOrder(&tmporder2)
	assert.Equal(t, "20220515 13:25:00 500", account.InitTime, "账户初始化时间更新不符合预期")
	assert.Equal(t, "20220516 13:45:27 500", account.UdTime, "账户刷新时间更新不符合预期")
	assert.Equal(t, 996172.04292, account.Fundavail, "账户可用资金更新不符合预期")
	assert.Equal(t, 96.00000000000009, account.AllProfit, "账户总利润更新不符合预期")
	assert.Equal(t, 15.957079999999998, account.AllCommission, "账户总佣金更新不符合预期")
	assert.Equal(t, 4500.0, account.Equity(), "账户总股权更新不符合预期")
	assert.Equal(t, 4.0, account.PosMap["600000"].CalPosTdyNum(), "账户持股数更新不符合预期")

	tmporder3 := order.NewStockOrder("600002", true, true, "20220516 14:45:00 500", 20.0, 1, "Buy", &orderscp)
	account.ActOnOrder(&tmporder3)
	assert.Equal(t, "20220515 13:25:00 500", account.InitTime, "账户初始化时间更新不符合预期")
	assert.Equal(t, "20220516 14:45:00 500", account.UdTime, "账户刷新时间更新不符合预期")
	assert.Equal(t, 994167.02292, account.Fundavail, "账户可用资金更新不符合预期")
	assert.Equal(t, 96.00000000000009, account.AllProfit, "账户总利润更新不符合预期")
	assert.Equal(t, 20.977079999999997, account.AllCommission, "账户总佣金更新不符合预期")
	assert.Equal(t, 6500.0, account.Equity(), "账户总股权更新不符合预期")
	assert.Equal(t, 4.0, account.PosMap["600000"].CalPosTdyNum(), "账户持股数更新不符合预期")

}

// test the NewStockAccountFromConfig
func TestNewStockAccountFromConfig(t *testing.T) {
	dir := "/Users/wonderstone/go/QuantTools/config/Manual/"
	cpm := cp.NewCPMap(dir, "ContractProp.yaml")

	SA := NewSAFromConfig("/Users/wonderstone/go/QuantTools/", "realtime.yaml", "VA.sacct", cpm)
	fmt.Println(SA)

}
