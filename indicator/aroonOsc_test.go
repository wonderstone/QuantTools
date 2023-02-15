package indicator

import (
	"fmt"

	"github.com/wonderstone/QuantTools/indicator/tools"

	"testing"
)

func TestAroonOsc_Eval(t *testing.T) {
	aroonOsc := NewAroonOsc("AroonOsc3", []int{3}, []string{"High", "Low"})
	aroonOsc.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 3.0, "High": 4.0, "Low": 1.0, "Vol": 2.0})
	aroonOsc.DQ.Enqueue(map[string]float64{"Open": 3.0, "Close": 3.0, "High": 6.0, "Low": 2.0, "Vol": 1.0})
	aroonOsc.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 6.0, "High": 5.0, "Low": 1.0, "Vol": 1.0})
	AroonOsc := aroonOsc.Eval()
	if !tools.CompareFloat(AroonOsc, -100/3.0) {
		fmt.Println("aroonOsc.Eval() :  ", AroonOsc)
		t.Error("Expected --- , got ", AroonOsc)
	}
}
