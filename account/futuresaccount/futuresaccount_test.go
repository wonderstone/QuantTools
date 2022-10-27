package futuresaccount

import (
	"testing"

	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/data"
	"github.com/wonderstone/QuantTools/order"

	"github.com/stretchr/testify/assert"
)

func TestNewFuturesAccount(t *testing.T) {
	// 注意，这里使用testify的assert进行对象比对，如果是指针或切片，比对的是具体的对象而非地址
	// 由于添加uuid，这个肯定过不了
	cash := 10000.0
	initTime := "20220515 13:35:27 500"
	res := NewFuturesAccount(initTime, cash)
	expected := FuturesAccount{
		InitTime:  "20220515 13:35:27 500",
		UdTime:    "20220515 13:35:27 500",
		BmkVal:    cash,
		MktVal:    cash,
		Fundavail: cash,
		PosMap:    make(map[string]*PositionSlice),
	}
	assert.Equal(t, &expected, &res, "account创建不符合预期")
}

func TestActOnOrder(t *testing.T) {
	cash := 1000000.0
	initTime := "20220515 13:25:00 500"
	account := NewFuturesAccount(initTime, cash)

	// 一个ContractProp给单子提供计算信息
	var orderfcp cp.FCP = cp.FCP{ContractSize: 1000, TickSize: 0.02, MarginLong: 10, MarginShort: 10, MarginBroker: 0, IsCommRateType: false, CommOpen: 2, CommCloseToday: 0, CommClosePrevious: 2, CommBroker: 0.01}
	// 下一单新 检查
	tmporder0 := order.NewFuturesOrder("au2210", true, true, "20220515 13:35:27 500", 400.00, 2, "Buy", "Open", &orderfcp)
	account.ActOnOrder(&tmporder0)

	assert.Equal(t, "20220515 13:25:00 500", account.InitTime, "账户初始化时间更新不符合预期")
	assert.Equal(t, "20220515 13:35:27 500", account.UdTime, "账户刷新时间更新不符合预期")
	assert.Equal(t, 0.0, account.AllProfit, "账户总利润更新不符合预期")
	assert.Equal(t, 4.02, account.AllCommission, "账户总佣金更新不符合预期")
	assert.Equal(t, 80000.0, account.Margin(), "账户总保证金更新不符合预期")
	assert.Equal(t, 999995.98, account.MktVal, "账户总市值更新不符合预期")
	assert.Equal(t, 999995.98, account.BmkVal, "账户总基准市值更新不符合预期")
	assert.Equal(t, 919995.98, account.Fundavail, "账户可用资金更新不符合预期")

	tmporder1 := order.NewFuturesOrder("au2210", true, true, "20220515 13:35:29 500", 410.00, 2, "Sell", "Open", &orderfcp)
	account.ActOnOrder(&tmporder1)

	assert.Equal(t, "20220515 13:25:00 500", account.InitTime, "账户初始化时间更新不符合预期")
	assert.Equal(t, "20220515 13:35:29 500", account.UdTime, "账户刷新时间更新不符合预期")
	assert.Equal(t, 0.0, account.AllProfit, "账户总利润更新不符合预期")
	assert.Equal(t, 8.04, account.AllCommission, "账户总佣金更新不符合预期")
	assert.Equal(t, 164000.0, account.Margin(), "账户总保证金更新不符合预期")
	assert.Equal(t, 1019991.96, account.MktVal, "账户总市值更新不符合预期")
	assert.Equal(t, 999991.96, account.BmkVal, "账户总基准市值更新不符合预期")
	assert.Equal(t, 835991.96, account.Fundavail, "账户可用资金更新不符合预期")

	tmporder2 := order.NewFuturesOrder("au2210", true, true, "20220515 13:45:00 500", 404.00, 1, "Sell", "CloseToday", &orderfcp)
	account.ActOnOrder(&tmporder2)
	assert.Equal(t, "20220515 13:25:00 500", account.InitTime, "账户初始化时间更新不符合预期")
	assert.Equal(t, "20220515 13:45:00 500", account.UdTime, "账户刷新时间更新不符合预期")
	assert.Equal(t, 4000.0, account.AllProfit, "账户总利润更新不符合预期")
	assert.Equal(t, 8.049999999999999, account.AllCommission, "账户总佣金更新不符合预期")
	assert.Equal(t, 121200.0, account.Margin(), "账户总保证金更新不符合预期")
	assert.Equal(t, 1019991.95, account.MktVal, "账户总市值更新不符合预期")
	assert.Equal(t, 1003991.95, account.BmkVal, "账户总基准市值更新不符合预期")
	assert.Equal(t, 882791.95, account.Fundavail, "账户可用资金更新不符合预期")

	umi := data.NewUpdateMI("20220515 14:35:00 500", "au2210", 415.00)

	account.ActOnUpdateMI(umi.UpdateTimeStamp, umi.InstID, umi.Value)
	assert.Equal(t, "20220515 13:25:00 500", account.InitTime, "账户初始化时间更新不符合预期")
	assert.Equal(t, "20220515 14:35:00 500", account.UdTime, "账户刷新时间更新不符合预期")
	assert.Equal(t, 4000.0, account.AllProfit, "账户总利润更新不符合预期")
	assert.Equal(t, 8.049999999999999, account.AllCommission, "账户总佣金更新不符合预期")
	assert.Equal(t, 124500.0, account.Margin(), "账户总保证金更新不符合预期")
	assert.Equal(t, 1008991.95, account.MktVal, "账户总市值更新不符合预期")
	assert.Equal(t, 1003991.95, account.BmkVal, "账户总基准市值更新不符合预期")
	assert.Equal(t, 879491.95, account.Fundavail, "账户可用资金更新不符合预期")

	account.ActOnMTM(umi.UpdateTimeStamp, umi.InstID, umi.Value)
	assert.Equal(t, "20220515 13:25:00 500", account.InitTime, "账户初始化时间更新不符合预期")
	assert.Equal(t, "20220515 14:35:00 500", account.UdTime, "账户刷新时间更新不符合预期")
	assert.Equal(t, 9000.0, account.AllProfit, "账户总利润更新不符合预期")
	assert.Equal(t, 8.049999999999999, account.AllCommission, "账户总佣金更新不符合预期")
	assert.Equal(t, 124500.0, account.Margin(), "账户总保证金更新不符合预期")
	assert.Equal(t, 1008991.95, account.MktVal, "账户总市值更新不符合预期")
	assert.Equal(t, 1008991.95, account.BmkVal, "账户总基准市值更新不符合预期")
	assert.Equal(t, 884491.95, account.Fundavail, "账户可用资金更新不符合预期")

}
