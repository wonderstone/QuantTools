package indicator

import (
	"testing"
	"fmt"
)
func TestEvalQstick(t *testing.T) {
	qs := NewQstick("Qstick5",[]int{5}, []string{"closing","opening"}) 
	qs.LoadData(map[string]float64{"closing": 20,"opening":10})
	fmt.Println(qs.Eval())
	qs.LoadData(map[string]float64{"closing": 15,"opening":20})
	fmt.Println(qs.Eval())
	qs.LoadData(map[string]float64{"closing": 50,"opening":15})
	fmt.Println(qs.Eval())
	qs.LoadData(map[string]float64{"closing": 55,"opening":50})
	fmt.Println(qs.Eval())
	qs.LoadData(map[string]float64{"closing": 42,"opening":40})
	fmt.Println(qs.Eval())
	qs.LoadData(map[string]float64{"closing": 30,"opening":35})
	actual:=qs.Eval()
	fmt.Println(actual)
	actual_2 := roundDigits(actual,2)
	

	if actual_2 != float64(6.4) {
		fmt.Println("qs.Eval() :  ", actual_2)
		t.Error("Expected 6.4, got ", actual_2) 
	}

}
