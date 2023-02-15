package indicator

import (
	"fmt"
	"testing"

	"github.com/wonderstone/QuantTools/indicator/tools"
)

func TestEvalIMI(t *testing.T) {

	imi := NewIMI("IMI3", []int{3}, []string{"Open", "Close"})
	imi.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 2.0, "High": 4.0, "Low": 1.0, "Vol": 2.0})
	imi.DQ.Enqueue(map[string]float64{"Open": 4.0, "Close": 3.0, "High": 6.0, "Low": 2.0, "Vol": 1.0})
	imi.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 4.0, "High": 5.0, "Low": 1.0, "Vol": 1.0})

	if !tools.CompareFloat(imi.Eval(), 200/3.0) {
		fmt.Println("imi.Eval() :  ", imi.Eval())
		t.Error("Expected --- , got ", imi.Eval())
	}
}
