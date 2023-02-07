package indicator

import (
	"fmt"
	"testing"
)

func TestEvalTema(t *testing.T) {
	tema := NewTema("Tema3", []int{3}, []string{"values"})
	tema.LoadData(map[string]float64{"values": 10})
	fmt.Println(tema.Eval())
	tema.LoadData(map[string]float64{"values": 12})
	fmt.Println(tema.Eval())
	tema.LoadData(map[string]float64{"values": 8})
	fmt.Println(tema.Eval())
	tema.LoadData(map[string]float64{"values": 10})
	fmt.Println(tema.Eval())

	if roundDigits(tema.Eval(), 2) != (10.25) {
		fmt.Println("d.Eval() :  ", tema.Eval())
		t.Error("Expected 10.25, got ", tema.Eval())
	}

}
