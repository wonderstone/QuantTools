package indicator

import (
	"fmt"
	"testing"
)

func TestEvalPPO(t *testing.T) {
	p := NewPPO("PPO234", []int{2, 3, 4}, []string{"price"})
	p.LoadData(map[string]float64{"price": 10})
	fmt.Println(p.Eval())
	p.LoadData(map[string]float64{"price": 12})
	fmt.Println(p.Eval())
	p.LoadData(map[string]float64{"price": 8})
	fmt.Println(p.Eval())
	p.LoadData(map[string]float64{"price": 10})
	fmt.Println(p.Eval())

	if roundDigits(p.Eval(), 2) != (-0.47) {
		fmt.Println("p.Eval() :  ", p.Eval())
		t.Error("Expected -0.47, got ", p.Eval())
	}

}
