package indicator

import (
	"fmt"
	"testing"
)

func TestEvalCMO(t *testing.T) {

	cmo := NewCMO("CMO3", []int{3}, []string{"Close"})
	cmo.DQ.Enqueue(map[string]float64{"Open": 3.0, "Close": 4.0, "High": 6.0, "Low": 2.0, "Vol": 2.0})
	cmo.DQ.Enqueue(map[string]float64{"Open": 4.0, "Close": 3.0, "High": 6.0, "Low": 2.0, "Vol": 1.0})
	cmo.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 3.0, "High": 5.0, "Low": 1.0, "Vol": 1.0})

	if cmo.Eval() != -100 {
		fmt.Println("cmo.Eval() :  ", cmo.Eval())
		t.Error("Expected --- , got ", cmo.Eval())
	}
}

func BenchmarkEvalCMO(b *testing.B) {
	cmo := NewCMO("CMO60", []int{60}, []string{"Close"})
	for i := 0; i < 60; i++ {
		cmo.DQ.Enqueue(map[string]float64{"Open": 3.0, "Close": 4.0, "High": 6.0, "Low": 2.0, "Vol": 2.0})
	}
	for i := 0; i < b.N; i++ {
		cmo.Eval() //cost 3384ns
	}
}
