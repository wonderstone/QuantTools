package indicator

import (
	"fmt"
	"testing"
)

func TestEvalVar(t *testing.T) {
	vari := NewVar([]int{3}, []string{"Close"})

	vari.LoadData(map[string]float64{"Close": 1.0})
	fmt.Println(vari.Eval(), vari.DQ.Full())
	vari.LoadData(map[string]float64{"Close": 2.0})
	fmt.Println(vari.Eval(), vari.DQ.Full())
	vari.LoadData(map[string]float64{"Close": 3.0})
	fmt.Println(vari.Eval(), vari.DQ.Full())
	vari.LoadData(map[string]float64{"Close": 4.0})
	fmt.Println(vari.Eval(), vari.DQ.Full())
	vari.LoadData(map[string]float64{"Close": 5.0})
	fmt.Println(vari.Eval(), vari.DQ.Full())
	vari.LoadData(map[string]float64{"Close": 6.0})
	fmt.Println(vari.Eval(), vari.DQ.Full())
	vari.LoadData(map[string]float64{"Close": 7.0})
	fmt.Println(vari.Eval(), vari.DQ.Full())

	if vari.Eval() != float64(2.0/3.0) {
		fmt.Println("vari.Eval() :  ", vari.Eval())
		t.Error("Expected 3.0, got ", vari.Eval())
	}

}
