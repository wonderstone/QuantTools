package indicator

import (
	"fmt"
	"testing"
)

func TestEval(t *testing.T) {
	bolling := NewBolling("Bolling3", []int{3, 1}, []string{"Close"})
	bolling.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 2.0, "High": 4.0, "Low": 1.0, "Vol": 2.0})
	bolling.DQ.Enqueue(map[string]float64{"Open": 4.0, "Close": 3.0, "High": 6.0, "Low": 2.0, "Vol": 1.0})
	bolling.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 4.0, "High": 5.0, "Low": 1.0, "Vol": 1.0})
	mid, upper, lower := bolling.Eval()

	if !(mid == 3 && upper == 4 && lower == 2) {
		fmt.Println("bolling.Eval() :  ", mid, upper, lower)
		t.Error("Expected --- , got ", mid, upper, lower)
	}
}
