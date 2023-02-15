package indicator

import (
	"fmt"
	"testing"

	"github.com/wonderstone/QuantTools/indicator/tools"
)

func TestEvalMP(t *testing.T) {
	medianprice := NewMedianPrice("MedianPrice1", []int{1}, []string{"High", "Low"})
	medianprice.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 2.0, "High": 4.0, "Low": 1.0, "Vol": 2.0})
	if !tools.CompareFloat(medianprice.Eval(), 2.5) {
		fmt.Println("medianprice.Eval() :  ", medianprice.Eval())
		t.Error("Expected --- , got ", medianprice.Eval())
	}
}
