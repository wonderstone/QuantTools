package indicator

import (
	"fmt"
	"testing"
)

func TestEvalAtr(t *testing.T) {
	a := NewAtr("Atr3", []int{3}, []string{"high","low","closing"})
	a.LoadData(map[string]float64{"high": 40, "low": 10, "closing": 20})
	fmt.Println(a.Eval())
	a.LoadData(map[string]float64{"high": 30, "low": 18, "closing": 20})
	fmt.Println(a.Eval())
	a.LoadData(map[string]float64{"high": 20, "low": 10, "closing": 12})
	fmt.Println(a.Eval())
	a.LoadData(map[string]float64{"high": 30, "low": 8, "closing": 10})
	fmt.Println(a.Eval())
	
	if roundDigits(a.Eval(), 2) != 14.67 {
		fmt.Println("a.Eval() :  ", a.Eval())
		t.Error("Expected 14.67, got ",a.Eval())
	}
	
}