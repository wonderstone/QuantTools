package perfeval

import (
	"strconv"
	"testing"

	"github.com/wonderstone/QuantTools/account"
)

// test NewPerfEval
func TestNewPerfEval(t *testing.T) {
	PE := NewPerfEval()
	// New a MktValSlice
	PE.MktValSlice = make([]account.MktValDataType, 0)
	// add 10 datas to MktValSlice
	for i := 1; i < 9; i++ {
		PE.MktValSlice = append(PE.MktValSlice, account.MktValDataType{
			Time:   "2018/01/0" + strconv.Itoa(i+1),
			MktVal: float64(i) + 300.2,
		})

	}
	// add the last data
	PE.MktValSlice = append(PE.MktValSlice, account.MktValDataType{
		Time:   "2018/01/01",
		MktVal: float64(3) + 300.2,
	})

	// change one value
	PE.MktValSlice[4].MktVal = float64(1) + 288.2
	// test Sort
	PE.Sort()
	if PE.MktValSlice[0].MktVal != 303.2 {
		t.Error("Sort failed")
	}

	// test CalcPerfEvalResult
	res := PE.CalcPerfEvalResult()
	if res.TotalReturn != 1.016490765171504 {
		t.Error("CalcPerfEvalResult failed")
	}
	if res.AnnualizedReturn != 1.5808703270674531 {
		t.Error("CalcPerfEvalResult failed")
	}
	if res.MaxDrawDown != 0.04930966469428011 {
		t.Error("CalcPerfEvalResult failed")
	}
	if res.SharpeRatio != 3.41925859530848 {
		t.Error("CalcPerfEvalResult failed")
	}

}
