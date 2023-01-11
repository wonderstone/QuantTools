package indicator

import (
	"fmt"
	"testing"
)

func TestEvalEMA(t *testing.T) {
	e := NewEMA("EMA3", []int{3}, []string{"Close"})
	e.LoadData(map[string]float64{"Close": 8996.96})
	fmt.Println(e.Eval())
	e.LoadData(map[string]float64{"Close": 9003.19})
	fmt.Println(e.Eval())
	e.LoadData(map[string]float64{"Close": 9010.41})
	fmt.Println(e.Eval())
	e.LoadData(map[string]float64{"Close": 9008.07})
	fmt.Println(e.Eval())
	e.LoadData(map[string]float64{"Close": 9018.03})
	fmt.Println(e.Eval())
	e.LoadData(map[string]float64{"Close": 9009.80})
	fmt.Println(e.Eval())

	if roundDigits(e.Eval(), 2) != 9011.07 {
		fmt.Println("e.Eval() :  ", e.Eval())
		t.Error("Expected 9012.34, got ", e.Eval())
	}

}
