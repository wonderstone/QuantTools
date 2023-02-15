package indicator

import (
	"fmt"
	"testing"
)

func TestEvalASI(t *testing.T) {
	asi := NewASI("ASI32", []int{3, 2}, []string{"Open", "Close", "High", "Low"})
	asi.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 3.0, "High": 4.0, "Low": 1.0, "Vol": 2.0})
	asi.DQ.Enqueue(map[string]float64{"Open": 4.0, "Close": 3.0, "High": 6.0, "Low": 2.0, "Vol": 1.0})
	asi.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 6.0, "High": 5.0, "Low": 1.0, "Vol": 1.0})
	astt := asi.Eval()

	if astt != 84 {
		fmt.Println("astt.Eval() :  ", astt)
		t.Error("Expected --- , got ", astt)
	}
}
