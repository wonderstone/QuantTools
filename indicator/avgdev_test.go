package indicator

import (
	"fmt"
	"testing"
)

func TestAVGDEVEval(t *testing.T) {
	avgdev := NewAvgDev([]int{5}, []string{"Close"})
	avgdev.LoadData(map[string]float64{"Close": 1.0})
	avgdev.LoadData(map[string]float64{"Close": 2.0})
	avgdev.LoadData(map[string]float64{"Close": 3.0})
	avgdev.LoadData(map[string]float64{"Close": 4.0})
	avgdev.LoadData(map[string]float64{"Close": 5.0})

	if avgdev.Eval() != 2.4 {
		fmt.Println("avgdev.Eval() :  ", avgdev.Eval())
		t.Error("Expected --- , got ", avgdev.Eval())
	}
}
