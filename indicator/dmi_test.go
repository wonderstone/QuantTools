package indicator

import (
	"fmt"
	"github.com/wonderstone/QuantTools/indicator/tools"
	"math"
	"testing"
)

func TestDMIEval(t *testing.T) {
	dmi := NewDMI("DMI422", []int{4, 2, 2}, []string{"Open", "Close"})
	dmi.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 2.0, "High": 4.0, "Low": 1.0, "Vol": 2.0})
	dmi.DQ.Enqueue(map[string]float64{"Open": 4.0, "Close": 3.0, "High": 6.0, "Low": 2.0, "Vol": 1.0})
	dmi.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 4.0, "High": 5.0, "Low": 1.0, "Vol": 1.0})
	dmi.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 4.0, "High": 5.0, "Low": 1.0, "Vol": 1.0})

	if !tools.CompareFloat(dmi.Eval(), 200/3.0) {
		fmt.Println("dmi.Eval() :  ", dmi.Eval())
		t.Error("Expected --- , got ", dmi.Eval())
	}
}

func TestDMIEval2(t *testing.T) {
	dmi := NewDMI("DMI422", []int{10, 2, 2}, []string{"Open", "Close"})
	dmi.DQ.Enqueue(map[string]float64{"Open": 21.32, "Close": 21.51, "High": 21.81, "Low": 21.32})
	dmi.DQ.Enqueue(map[string]float64{"Open": 21.66, "Close": 22.09, "High": 22.09, "Low": 21.48})
	dmi.DQ.Enqueue(map[string]float64{"Open": 22.01, "Close": 22.16, "High": 22.23, "Low": 21.91})
	dmi.DQ.Enqueue(map[string]float64{"Open": 22.07, "Close": 22.26, "High": 22.6, "Low": 22})
	dmi.DQ.Enqueue(map[string]float64{"Open": 22.2, "Close": 22.65, "High": 22.69, "Low": 22.15})
	dmi.DQ.Enqueue(map[string]float64{"Open": 22.65, "Close": 22.28, "High": 22.75, "Low": 22.23})
	dmi.DQ.Enqueue(map[string]float64{"Open": 22.3, "Close": 22.9, "High": 23, "Low": 22.2})
	dmi.DQ.Enqueue(map[string]float64{"Open": 22.99, "Close": 23, "High": 23.06, "Low": 22.73})

	dmi.DQ.Enqueue(map[string]float64{"Open": 22.93, "Close": 23.08, "High": 23.14, "Low": 22.87})
	dmi.DQ.Enqueue(map[string]float64{"Open": 22.98, "Close": 22.64, "High": 23.2, "Low": 22.58})
	//dmi.DQ.Enqueue(map[string]float64{"Open": 600.3, "Close": 603.18, "High": 604.01, "Low": 599.26})
	//dmi.DQ.Enqueue(map[string]float64{"Open": 602.62, "Close": 598.33, "High": 606.62, "Low": 597.43})

	//dmi.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 1.0, "High": 3.0, "Low": 1.0, "Vol": 1.0})
	//dmi.DQ.Enqueue(map[string]float64{"Open": 1.0, "Close": 2.0, "High": 2.0, "Low": 1.0, "Vol": 2.0})
	//dmi.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 1.0, "High": 3.0, "Low": 1.0, "Vol": 1.0})

	if math.Abs(dmi.Eval()-78.38) > 0.01 {
		fmt.Println("dmi.Eval() :  ", dmi.Eval())
		t.Error("Expected --- , got ", dmi.Eval())
	}
}
