package indicator

import (
	"fmt"
	"testing"
)

func TestBollingUpper_Eval(t *testing.T) {
	bollingUpper := NewBollingUpper("BollingUpper3", []int{3, 1}, []string{"Close"})
	bollingUpper.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 2.0, "High": 4.0, "Low": 1.0, "Vol": 2.0})
	bollingUpper.DQ.Enqueue(map[string]float64{"Open": 4.0, "Close": 3.0, "High": 6.0, "Low": 2.0, "Vol": 1.0})
	bollingUpper.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 4.0, "High": 5.0, "Low": 1.0, "Vol": 1.0})
	upper := bollingUpper.Eval()

	if !(upper == 4.0) {
		fmt.Println("bollingUpper.Eval() :  ", upper)
		t.Error("Expected --- , got ", upper)
	}
}
