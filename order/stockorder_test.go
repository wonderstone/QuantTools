package order

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	cp "github.com/wonderstone/QuantTools/contractproperty"
)

// test NewStockOrder
func TestNewStockOrder(t *testing.T) {
	// new a scp from code
	fmt.Println("test NewStockOrder")
	confName := "ContractProp.yaml"
	dir := "../config/Manual/"
	cpm := cp.NewCPMap(confName, dir)
	instID := "SZ000058"
	scp := cp.SimpleNewSCPFromMap(cpm, instID)

	// so := NewStockOrder(instID, false, "2022-05-10 14:52", 8.5, 2.0, Buy, &scp)
	so := NewStockOrder(instID, true, false, "2022-05-10 14:52", 8.5, 2.0, "Buy", &scp)
	assert.Equal(t, instID, so.InstID, "instID should be equal")
	assert.Equal(t, true, so.IsEligible, "IsEligible should be equal")
	assert.Equal(t, false, so.IsExecuted, "isexecuted should be equal")
	assert.Equal(t, "2022-05-10 14:52", so.OrderTime, "ordertime should be equal")
	assert.Equal(t, 8.5, so.OrderPrice, "orderprice should be equal")
	assert.Equal(t, 2.0, so.OrderNum, "ordernum should be equal")
	assert.Equal(t, "Buy", so.OrderDirection, "orderdir should be equal")
	assert.Equal(t, &scp, so.SCP, "scp should be equal")
}

// test CalEquity
func TestCalEquity(t *testing.T) {
	// new a fcp from instID
	fmt.Println("test NewStockOrder")
	confName := "ContractProp"
	dir := "../config/Manual"
	cpm := cp.NewCPMap(confName, dir)
	instID := "SZ000058"
	scp := cp.SimpleNewSCPFromMap(cpm, instID)

	// so := NewStockOrder(instID, false, "2022-05-10 14:52", 8.5, 2.0, Buy, &scp)
	so := NewStockOrder(instID, true, false, "2022-05-10 14:52", 8.5, 2.0, "Buy", &scp)
	assert.Equal(t, 1717.0, so.CalEquity(), fmt.Sprintf("so.CalEquity() should be 1717.0, but %v", so.CalEquity()))
}
