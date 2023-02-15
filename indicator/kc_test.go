package indicator

import (
	"fmt"
	"testing"
)

func TestEvalKC(t *testing.T) {
	k := NewKC("KC3", []int{3}, []string{"high","low","closing"})
	k.LoadData(map[string]float64{"high": 40, "low": 10, "closing": 20})
	fmt.Println(k.Eval())
	k.LoadData(map[string]float64{"high": 30, "low": 18, "closing": 20})
	fmt.Println(k.Eval())
	k.LoadData(map[string]float64{"high": 20, "low": 10, "closing": 12})
	fmt.Println(k.Eval())
	k.LoadData(map[string]float64{"high": 30, "low": 8, "closing": 10})
	fmt.Println(k.Eval())
	
	if roundDigits(k.upperBand, 2) != 42.33 {
		fmt.Println("k.upperBand :  ", k.upperBand)
		t.Error("Expected 42.33, got ",k.upperBand)
	}
	
}