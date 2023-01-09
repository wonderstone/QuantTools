package indicator

import (
	"testing"
	"fmt"
)

func TestEvalSum(t *testing.T) {
	s := NewSum("Sum4",[]int{4}, []string{"Close"}) 
	s.LoadData(map[string]float64{"Close": 1.0})
	fmt.Println(s.Eval())
	s.LoadData(map[string]float64{"Close": 2.0})
	fmt.Println(s.Eval())
	s.LoadData(map[string]float64{"Close": 3.0})
	fmt.Println(s.Eval())
	s.LoadData(map[string]float64{"Close": 4.0})
	fmt.Println(s.Eval())
	s.LoadData(map[string]float64{"Close": 5.0})
	fmt.Println(s.Eval())
	

	if s.Eval() != float64(14) {
		fmt.Println("s.Eval() :  ", s.Eval())
		t.Error("Expected 14, got ", s.Eval()) 
	}

}
