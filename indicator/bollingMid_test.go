package indicator

import (
	"fmt"
	"testing"
)

func TestBollingMid_Eval(t *testing.T) {
	bollingMid := NewBollingMid("BollingMid3", []int{3, 1}, []string{"Close"})
	bollingMid.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 2.0, "High": 4.0, "Low": 1.0, "Vol": 2.0})
	bollingMid.DQ.Enqueue(map[string]float64{"Open": 4.0, "Close": 3.0, "High": 6.0, "Low": 2.0, "Vol": 1.0})
	bollingMid.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 4.0, "High": 5.0, "Low": 1.0, "Vol": 1.0})
	lower := bollingMid.Eval()

	if !(lower == 3.0) {
		fmt.Println("bollingMid.Eval() :  ", lower)
		t.Error("Expected --- , got ", lower)
	}
}
