package indicator

import (
	"fmt"
	"testing"
)

func TestEvalUI(t *testing.T) {
	ui := NewUI("UI", []int{14}, []string{"closing"})
	ui.LoadData(map[string]float64{"closing": 9})
	fmt.Println(ui.Eval())
	ui.LoadData(map[string]float64{"closing": 11})
	fmt.Println(ui.Eval())
	ui.LoadData(map[string]float64{"closing": 7})
	fmt.Println(ui.Eval())
	ui.LoadData(map[string]float64{"closing": 10})
	fmt.Println(ui.Eval())
	if roundDigits(ui.Eval(),2) != float64(18.74) {
		fmt.Println("ui.Eval() :  ", ui.Eval())
		t.Error("Expected 18.74, got ", ui.Eval()) 
	}

}