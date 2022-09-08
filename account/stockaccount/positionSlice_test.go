package stockaccount

import (
	"fmt"
	"testing"

	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/data"
	"github.com/wonderstone/QuantTools/order"

	"github.com/stretchr/testify/assert"
)

// test UpdateWithOrder
func TestUpdateWithOrder(t *testing.T) {

	// new a scp from code
	fmt.Println("test NewStockOrder")
	confName := "ContractProp"
	dir := "../../config/Manual"
	cpm := cp.NewCPMap(confName, dir)
	instID := "SZ000058"
	scp := cp.SimpleNewSCPFromMap(cpm, instID)

	so1 := order.NewStockOrder(instID, false, "2022-05-10 14:52", 8.5, 2.0, order.Buy, &scp)

	// new  positiondetails from stockorder
	pd1 := NewPositionDetail(&so1)
	// new 1 positionslice from positiondetail
	ps := NewPosSlice()
	ps.UpdateWithOrder(&so1)
	expected1 := pd1
	assert.Equal(t, expected1, ps.PosTdys[0], "NewPositionSlice不符合预期")

	// test add order to positionslice
	so2 := order.NewStockOrder(instID, false, "2022-05-10 14:53", 9.5, 2.0, order.Buy, &scp)
	ps.UpdateWithOrder(&so2)
	assert.Equal(t, "2022-05-10 14:53", ps.PosTdys[0].UdTime, "UpdateWithOrder不符合预期")
	assert.Equal(t, "2022-05-10 14:53", ps.PosTdys[1].UdTime, "UpdateWithOrder不符合预期")
	assert.Equal(t, 9.5, ps.PosTdys[0].LastPrice, "UpdateWithOrder不符合预期")
	assert.Equal(t, 9.5, ps.PosTdys[1].LastPrice, "UpdateWithOrder不符合预期")

	// test UpdateWithUMI
	umi := data.NewUpdateMI("2022-05-10 14:54", "SZ000058", 9.0)
	ps.UpdateWithUMI(umi.UpdateTimeStamp, umi.Value)
	assert.Equal(t, "2022-05-10 14:54", ps.PosTdys[0].UdTime, "UpdateWithUMI不符合预期")
	assert.Equal(t, 9.0, ps.PosTdys[0].LastPrice, "UpdateWithUMI不符合预期")
	assert.Equal(t, 9.0, ps.PosTdys[1].LastPrice, "UpdateWithUMI不符合预期")

	// test UpdateWithCM
	ps.UpdateWithCM()
	assert.Equal(t, 0, len(ps.PosTdys), "UpdateWithCM不符合预期")

	// test UpdateWithOrder with order.Sell
	fmt.Printf("the num of holdings  before in positionslice is: %d", int(ps.CalPosTdyNum()+ps.CalPosPrevNum()))
	so3 := order.NewStockOrder("SZ000058", false, "2022-05-11 14:54", 8.0, 3.0, order.Sell, &scp)
	ps.UpdateWithOrder(&so3)
	fmt.Printf("the num of holdings  after in positionslice is: %d", int(ps.CalPosTdyNum()+ps.CalPosPrevNum()))
	assert.Equal(t, "2022-05-11 14:54", ps.PosPrevs[0].UdTime, "UpdateWithOrder不符合预期")
	assert.Equal(t, 8.0, ps.PosPrevs[0].LastPrice, "UpdateWithOrder不符合预期")
	assert.Equal(t, 1.0, ps.CalPosTdyNum()+ps.CalPosPrevNum(), "UpdateWithOrder不符合预期")
}
