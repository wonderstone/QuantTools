package order

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	cp "github.com/wonderstone/QuantTools/contractproperty"
)

// test NewFuturesOrder
func TestNewFuturesOrder(t *testing.T) {
	// new a fcp from instID
	fmt.Println("test NewFuturesOrder")
	confName := "ContractProp"
	dir := "../config/Manual"
	cpm := cp.NewCPMap(confName, dir)
	instID := "au2210"
	fcp := cp.SimpleNewFCPFromMap(cpm, instID)

	// fo := NewFuturesOrder(instID, false, "2022-05-10 14:52", 8.5, 2.0, Buy, Open, &fcp)
	fo := NewFuturesOrder(instID, false, "2022-05-10 14:52", 8.5, 2.0, "Buy", "Open", &fcp)

	assert.Equal(t, instID, fo.InstID, "instID should be equal")
	assert.Equal(t, false, fo.IsExecuted, "isexecuted should be equal")
	assert.Equal(t, "2022-05-10 14:52", fo.OrderTime, "ordertime should be equal")
	assert.Equal(t, 8.5, fo.OrderPrice, "orderprice should be equal")
	assert.Equal(t, 2.0, fo.OrderNum, "ordernum should be equal")
	assert.Equal(t, "Buy", fo.OrderDirection, "orderdir should be equal")
	assert.Equal(t, &fcp, fo.FCP, "fcp should be equal")

}

//test CalMargin
func TestCalMargin(t *testing.T) {
	// new a fcp from instID
	fmt.Println("test CalMargin")
	confName := "ContractProp"
	dir := "../config/Manual"
	cpm := cp.NewCPMap(confName, dir)
	instID := "au2210"
	fcp := cp.SimpleNewFCPFromMap(cpm, instID)

	// fo := NewFuturesOrder(instID, true, "20220515 13:35:27 500", 400.00, 2, Buy, Open, &fcp)
	fo := NewFuturesOrder(instID, true, "20220515 13:35:27 500", 400.00, 2, "Buy", "Open", &fcp)
	assert.Equal(t, 80000.0, fo.CalMargin(), fmt.Sprintf("fo.CalMargin() should be 80000.0, but %v", fo.CalMargin()))
}
