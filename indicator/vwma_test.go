package indicator

import (
	"fmt"
	"testing"
)

func TestEvalVWMA(t *testing.T) {
	v := NewVWMA("VWMA3",[]int{3}, []string{"closing", "volume"})
	v.LoadData(map[string]float64{"closing": 20, "volume": 10})
	fmt.Println(v.Eval())
	v.LoadData(map[string]float64{"closing": 10, "volume": 20})
	fmt.Println(v.Eval())

	if v.Eval() != float64(203/23) {
		fmt.Println("v.Eval() :  ", v.Eval())
		t.Error("Expected 8.826, got ", v.Eval())
	}

}
