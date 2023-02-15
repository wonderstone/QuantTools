package indicator

import (
	"fmt"
	"testing"
)

func TestEvalCMF(t *testing.T) {

	cmf := NewCMF("CMF3", []int{3}, []string{"Close", "High", "Low", "Vol"})
	cmf.DQ.Enqueue(map[string]float64{"Open": 3.0, "Close": 4.0, "High": 6.0, "Low": 2.0, "Vol": 2.0})
	cmf.DQ.Enqueue(map[string]float64{"Open": 4.0, "Close": 3.0, "High": 6.0, "Low": 2.0, "Vol": 1.0})
	cmf.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 4.0, "High": 5.0, "Low": 1.0, "Vol": 1.0})

	if cmf.Eval() != -0.125 { //
		fmt.Println("cmf.Eval() :  ", cmf.Eval())
		t.Error("Expected -0.125 , got ", cmf.Eval())
	}
}
