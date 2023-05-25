package strategyModule

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wonderstone/QuantTools/order"
)

// test ContainNaN
func TestContainNaN(t *testing.T) {
	m := map[string]float64{"a": 1.0, "b": math.NaN(), "c": 3}
	assert.Equal(t, ContainNaN(m), true)
}

// test NetSOrders
func TestNetSOrders(t *testing.T) {
	// new 3 stock orders
	so1 := order.NewStockOrder("000001", true, false, "2016-01-01", 10.0, 110.0, "buy", nil)
	so2 := order.NewStockOrder("000001", true, false, "2016-01-01", 10.0, 100.0, "sell", nil)
	so3 := order.NewStockOrder("000002", true, false, "2016-01-01", 10.0, 100.0, "buy", nil)
	tmporslice := []order.StockOrder{so1, so2, so3}
	// net the orders
	res := NetSOrders(tmporslice)
	// check the result
	assert.Equal(t, len(res), 2)
	assert.Equal(t, res[0].InstID, "000001")
	assert.Equal(t, res[0].OrderNum, 10.0)
	assert.Equal(t, res[0].OrderDirection, "buy")

}
