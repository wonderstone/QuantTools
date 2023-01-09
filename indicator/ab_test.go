package indicator

import (
	"fmt"
	"testing"
)

func TestEvalAB(t *testing.T) {
	ab := NewAB("AB3",[]int{3}, []string{"high", "low", "closing"})
	ab.LoadData(map[string]float64{"high": 40, "low": 10, "closing": 20})
	fmt.Println(ab.Eval())
	ab.LoadData(map[string]float64{"high": 30, "low": 10, "closing": 10})
	fmt.Println(ab.Eval())

	if ab.Eval() != 90 {
		fmt.Println("ab.Eval() :  ", ab.Eval())
		t.Error("Expected 90, got ", ab.Eval())
	}

}
