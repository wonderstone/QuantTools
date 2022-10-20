package indicator

import (
	"fmt"
	"testing"
)

func TestEval(t *testing.T) {
	ma := NewMA(3, []string{"Close"})
	ma.LoadData(map[string]float64{"Close": 1.0})
	ma.LoadData(map[string]float64{"Close": 2.0})
	ma.LoadData(map[string]float64{"Close": 3.0})

	if ma.Eval() != 2.0 {
		fmt.Println("ma.Eval() :  ", ma.Eval())
		t.Error("Expected 2.0, got ", ma.Eval())
	}
}
