package indicator

import (
	"fmt"
	"testing"
)

func TestEvalMI(t *testing.T) {
	m := NewMI("MI3", []int{3}, []string{"high", "low"})
	m.LoadData(map[string]float64{"high": 10, "low": 6})
	fmt.Println(m.Eval())
	m.LoadData(map[string]float64{"high": 9, "low": 7})
	fmt.Println(m.Eval())
	m.LoadData(map[string]float64{"high": 12, "low": 9})
	fmt.Println(m.Eval())
	m.LoadData(map[string]float64{"high": 14, "low": 12})
	fmt.Println(m.Eval())

	if roundDigits(m.Eval(), 2) != 4.39 {
		fmt.Println("m.Eval() :  ", m.Eval())
		t.Error("Expected 4.38, got ", m.Eval())
	}

}
