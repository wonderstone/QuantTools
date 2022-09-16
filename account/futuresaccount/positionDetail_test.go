package futuresaccount

import (
	"testing"

	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/order"

	"github.com/stretchr/testify/assert"
)

// test NewPositionDetail
func TestNewPositionDetail(t *testing.T) {
	fcp := cp.NewFCP(1000, 0.02, 10, 10, 0.0, false, 2.0, 0.0, 2.0, 0.01)
	fo := order.NewFuturesOrder("au2210", true, "20220515 13:35:27 500", 400.00, 2, "Buy", "Open", &fcp)

	pd := NewPositionDetail(&fo)

	expected := PositionDetail{
		UdTime: "20220515 13:35:27 500",
		InstID: "au2210",
		// Dir:       DirType(order.Buy),
		Dir:       "Buy",
		BasePrice: 400.00,
		LastPrice: 400.00,
		Margin:    80000.0,
		Num:       2,
		FCP:       &fcp,
	}
	assert.Equal(t, expected, pd, "NewPositionDetail不符合预期")

	// test UpdateLastPrice
	pd.UpdateLastPrice("20220515 13:35:28 500", 401.00)
	expected = PositionDetail{
		UdTime: "20220515 13:35:28 500",
		InstID: "au2210",
		// Dir:       DirType(order.Buy),
		Dir:       "Buy",
		BasePrice: 400.00,
		LastPrice: 401.00,
		Margin:    80200.0,
		Num:       2,
		FCP:       &fcp,
	}
	assert.Equal(t, expected, pd, "UpdateLastPrice不符合预期")

	// test CalComm
	expectComm := 4.02
	actual := CalComm(&fo)
	assert.Equal(t, expectComm, actual, "CalComm不符合预期")

	// test CalUnRealizedProfit
	expectUnRealizedProfit := 2000.0
	actual = pd.CalUnRealizedProfit(401.00)
	assert.Equal(t, expectUnRealizedProfit, actual, "CalUnRealizedProfit不符合预期")

	// test CalRealizedProfit
	expectRealizedProfit := 1000.0
	actual = pd.CalRealizedProfit(1, 401.00)
	assert.Equal(t, expectRealizedProfit, actual, "CalRealizedProfit不符合预期")

}
