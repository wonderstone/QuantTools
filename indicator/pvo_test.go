package indicator

import (
	"fmt"
	"testing"
)

func TestEvalPVO(t *testing.T) {
	p := NewPVO("PVO234", []int{2,3,4}, []string{"volume"})
	p.LoadData(map[string]float64{"volume":10})
	fmt.Println(p.Eval())
	p.LoadData(map[string]float64{"volume":12})
	fmt.Println(p.Eval())
	p.LoadData(map[string]float64{"volume":8})
	fmt.Println(p.Eval())
	p.LoadData(map[string]float64{"volume":10})
	fmt.Println(p.Eval())

	if roundDigits(p.Eval(), 2) != (-0.47) {
		fmt.Println("p.Eval() :  ", p.Eval())
		t.Error("Expected -0.47, got ", p.Eval())
	}

}
