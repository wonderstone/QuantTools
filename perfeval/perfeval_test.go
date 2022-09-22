package perfeval

import (
	"fmt"
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

	if PE.CalcPerfEvalResult("TR") != 1.016490765171504 {
		fmt.Println(PE.CalcPerfEvalResult("TR"))
		t.Error("CalcPerfEvalResult failed")
	}
	if PE.CalcPerfEvalResult("AR") != 1.5808703270674531 {
		fmt.Println(PE.CalcPerfEvalResult("AR"))
		t.Error("CalcPerfEvalResult failed")
	}
	if PE.CalcPerfEvalResult("MR") != 32.06005023292793 {
		fmt.Println(PE.CalcPerfEvalResult("MR"))
		t.Error("CalcPerfEvalResult failed")
	}
	if PE.CalcPerfEvalResult("SR") != 3.354371706040631 {
		fmt.Println(PE.CalcPerfEvalResult("SR"))
		t.Error("CalcPerfEvalResult failed")
	}

}
