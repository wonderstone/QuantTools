package indicator

import (
	"fmt"
	"testing"
)

func TestEvalCCI(t *testing.T) {
	cci := NewCCI("CCI3", []int{3}, []string{"Close", "High", "Low"})
	cci.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 2.0, "High": 4.0, "Low": 1.0, "Vol": 2.0})
	cci.DQ.Enqueue(map[string]float64{"Open": 4.0, "Close": 3.0, "High": 6.0, "Low": 2.0, "Vol": 1.0})
	cci.DQ.Enqueue(map[string]float64{"Open": 2.0, "Close": 4.0, "High": 5.0, "Low": 1.0, "Vol": 1.0})

	if cci.Eval() != 28.61 { //通达信显示28.61,注意到显示数据做四舍五入处理,误差可能由此而来
		fmt.Println("cci.Eval() :  ", cci.Eval())
		t.Error("Expected --- , got ", cci.Eval())
	}
}
func BenchmarkEvalCCI(b *testing.B) {
	cci := NewCCI("CCI3", []int{60}, []string{"Close", "High", "Low"})
	for i := 0; i < 60; i++ {
		cci.DQ.Enqueue(map[string]float64{"Open": 4.0, "Close": 3.0, "High": 6.0, "Low": 2.0, "Vol": 1.0})
	}
	for i := 0; i < b.N; i++ {
		cci.Eval() //cost 3295ns
	}
}
