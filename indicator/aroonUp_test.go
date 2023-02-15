package indicator

import (
	"fmt"
	"github.com/wonderstone/QuantTools/indicator/tools"

	"testing"
)

func TestAroonUp_Eval(t *testing.T) {
	aroonUp := NewAroonUp("AROON3", []int{3}, []string{"High", "Low"})
	aroonUp.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 3.0, "High": 4.0, "Low": 1.0, "Vol": 2.0})
	aroonUp.DQ.Enqueue(map[string]float64{"Open": 3.0, "Close": 3.0, "High": 6.0, "Low": 2.0, "Vol": 1.0})
	aroonUp.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 6.0, "High": 5.0, "Low": 1.0, "Vol": 1.0})
	AroonUp := aroonUp.Eval()
	if !tools.CompareFloat(AroonUp, 200/3.0) {
		fmt.Println("aroon.Eval() :  ", AroonUp)
		t.Error("Expected --- , got ", AroonUp)
	}
}
