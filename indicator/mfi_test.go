package indicator

import (
	"fmt"
	"testing"
)

func TestEvalMFI(t *testing.T) {
	mfi := NewMFI("MFI3", []int{3}, []string{"Close", "High", "Low", "Vol"})
	mfi.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 2.0, "High": 4.0, "Low": 1.0, "Vol": 2.0})
	mfi.DQ.Enqueue(map[string]float64{"Open": 4.0, "Close": 3.0, "High": 6.0, "Low": 2.0, "Vol": 1.0})
	mfi.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 4.0, "High": 5.0, "Low": 1.0, "Vol": 1.0})

	if mfi.Eval() != 40 {
		fmt.Println("mfi.Eval() :  ", mfi.Eval())
		t.Error("Expected 40 , got ", mfi.Eval())
	}
}
