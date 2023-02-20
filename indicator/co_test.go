package indicator

import (
	"fmt"
	"testing"
)

func TestEvalCO(t *testing.T) {
	c := NewCO("CO3", []int{2, 3}, []string{"high", "low", "closing", "volume"})
	c.LoadData(map[string]float64{"high": 10, "low": 6, "closing": 9, "volume": 100})
	fmt.Println(c.Ema1.Eval(), c.Ema2.Eval(), c.Eval())
	c.LoadData(map[string]float64{"high": 9, "low": 7, "closing": 11, "volume": 110})
	fmt.Println(c.Ema1.Eval(), c.Ema2.Eval(), c.Eval())
	c.LoadData(map[string]float64{"high": 12, "low": 9, "closing": 7, "volume": 80})
	fmt.Println(c.Ema1.Eval(), c.Ema2.Eval(), c.Eval())
	if roundDigits(c.Eval(), 2) != 14.72 {
		fmt.Println("c.Eval() :  ", c.Eval())
		t.Error("Expected 14.72, got ", c.Eval())
	}

}
