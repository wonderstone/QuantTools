package stockaccount

import (
	"fmt"
	"testing"

	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/order"

	"github.com/stretchr/testify/assert"
)

// test NewPositionDetail
func TestPositionDetail(t *testing.T) {
	// new a scp from code
	fmt.Println("test NewStockOrder")
	confName := "ContractProp"
	dir := "../../config/Manual"
	cpm := cp.NewCPMap(confName, dir)
	instID := "SZ000058"
	scp := cp.SimpleNewSCPFromMap(cpm, instID)

	so := order.NewStockOrder(instID, false, "2022-05-10 14:52", 8.5, 2.0, "Buy", &scp)
	// new 1 positiondetail from stockorder
	pd := NewPositionDetail(&so)
	// test
	expected := PositionDetail{
		UdTime:    "2022-05-10 14:52",
		InstID:    "SZ000058",
		BasePrice: 8.5,
		LastPrice: 8.5,
		Num:       2.0,
		Equity:    1700.0,
		SCP:       &scp,
	}
	assert.Equal(t, expected, pd, "NewPositionDetail不符合预期")

	// test UpdateLastPrice
	pd.UpdateLastPrice("2022-05-10 14:53", 9.5)
	expected = PositionDetail{
		UdTime:    "2022-05-10 14:53",
		InstID:    "SZ000058",
		BasePrice: 8.5,
		LastPrice: 9.5,
		Num:       2.0,
		Equity:    1900.0,
		SCP:       &scp,
	}
	assert.Equal(t, expected, pd, "UpdateLastPrice不符合预期")

	// test CalComm
	expectComm := 5.017
	actual := CalComm(&so)
	assert.Equal(t, expectComm, actual, "CalComm不符合预期")

	// test CalUnRealizedProfit
	expectUnRealizedProfit := 200.0
	actual = pd.CalUnRealizedProfit(9.5)
	assert.Equal(t, expectUnRealizedProfit, actual, "CalUnRealizedProfit不符合预期")

	// test CalRealizedProfit
	expectRealizedProfit := 100.0
	actual = pd.CalRealizedProfit(1, 9.5)
	assert.Equal(t, expectRealizedProfit, actual, "CalRealizedProfit不符合预期")

}
