package futuresaccount

import (
	"testing"

	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/order"

	"github.com/stretchr/testify/assert"
)

// test NewPositionSlice
// add closenum > pos num scenario
func TestPositionSlice(t *testing.T) {
	fcp := cp.NewFCP(1000, 0.02, 10, 10, 0.0, false, 2.0, 0.0, 2.0, 0.01)
	fo := order.NewFuturesOrder("au2210", true, "20220515 13:35:27 500", 400.00, 2, order.Buy, order.Open, &fcp)
	pd := NewPositionDetail(&fo)
	ps := NewPosSlice()

	ps.UpdateWithOrder(&fo)
	expected := pd
	assert.Equal(t, expected, ps.PosTdys[0], "NewPositionSlice不符合预期")

	// add one more buy position
	fo2 := order.NewFuturesOrder("au2210", true, "20220515 14:35:27 500", 410.00, 2, order.Buy, order.Open, &fcp)
	pd2 := NewPositionDetail(&fo2)
	ps.UpdateWithOrder(&fo2)
	expected2 := pd2
	assert.Equal(t, expected2, ps.PosTdys[1], "NewPositionSlice不符合预期")

	// add one more sell position
	fo3 := order.NewFuturesOrder("au2210", true, "20220515 15:35:27 500", 420.00, 2, order.Sell, order.Open, &fcp)
	pd3 := NewPositionDetail(&fo3)
	ps.UpdateWithOrder(&fo3)
	expected3 := pd3
	assert.Equal(t, expected3, ps.PosTdys[2], "NewPositionSlice不符合预期")

	// close one buy position
	fo4 := order.NewFuturesOrder("au2210", true, "20220515 16:35:27 500", 430.00, 4, order.Sell, order.CloseToday, &fcp)

	RealizedProfit, Comm, UnRealizedProfit := ps.UpdateWithOrder(&fo4)
	PosTdyNumL, PosTdyNumS := ps.CalPosTdyNum()
	assert.Equal(t, 0.0, PosTdyNumL, "NewPositionSlice不符合预期")
	assert.Equal(t, 2.0, PosTdyNumS, "NewPositionSlice不符合预期")
	assert.Equal(t, 100000.0, RealizedProfit, "NewPositionSlice不符合预期")
	assert.Equal(t, 0.04, Comm, "NewPositionSlice不符合预期")
	assert.Equal(t, -20000.0, UnRealizedProfit, "NewPositionSlice不符合预期")

}

// test UpdateWithUMI
func TestUpdateWithUMI(t *testing.T) {
	fcp := cp.NewFCP(1000, 0.02, 10, 10, 0.0, false, 2.0, 0.0, 2.0, 0.01)
	fo := order.NewFuturesOrder("au2210", true, "20220515 13:35:27 500", 400.00, 2, order.Buy, order.Open, &fcp)
	pd := NewPositionDetail(&fo)
	ps := NewPosSlice()

	ps.UpdateWithOrder(&fo)
	expected := pd
	assert.Equal(t, expected, ps.PosTdys[0], "UpdateWithOrder不符合预期")

	ps.UpdateWithUMI("20220515 14:35:00 500", 415.00)
	expected.UpdateLastPrice("20220515 14:35:00 500", 415.00)
	assert.Equal(t, expected, ps.PosTdys[0], "UpdateWithUMI不符合预期")
}

// test UpdateWithMTM
func TestUpdateWithMTM(t *testing.T) {
	fcp := cp.NewFCP(1000, 0.02, 10, 10, 0.0, false, 2.0, 0.0, 2.0, 0.01)
	fo := order.NewFuturesOrder("au2210", true, "20220515 13:35:27 500", 400.00, 2, order.Buy, order.Open, &fcp)
	pd := NewPositionDetail(&fo)
	ps := NewPosSlice()

	ps.UpdateWithOrder(&fo)
	expected := pd
	assert.Equal(t, expected, ps.PosTdys[0], "UpdateWithOrder不符合预期")

	ps.UpdateWithUMI("20220515 14:35:00 500", 415.00)
	expected.UpdateLastPrice("20220515 14:35:00 500", 415.00)
	assert.Equal(t, expected, ps.PosTdys[0], "UpdateWithUMI不符合预期")

	ps.UpdateWithMTM("20220515 15:35:00 500", 420.00)
	expected.UpdateBasePrice("20220515 15:35:00 500", 420.00)
	expected.UpdateLastPrice("20220515 15:35:00 500", 420.00)
	assert.Equal(t, expected, ps.PosPrevs[0], "UpdateWithMTM不符合预期")
}
