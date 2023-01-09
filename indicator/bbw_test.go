package indicator
import (
	"testing"
	"fmt"
)
func TestEvalBBW(t *testing.T) {
	bbw := NewBBW("BBW3",[]int{3}, []string{"upperBand", "middleBand", "lowerBand"}) 
	bbw.LoadData(map[string]float64{"upperBand":30, "middleBand":20, "lowerBand":10})
	fmt.Println(bbw.Eval())
	bbw.LoadData(map[string]float64{"upperBand":40, "middleBand":15, "lowerBand":10})
	fmt.Println(bbw.Eval())
	
	
	if bbw.Eval() != float64(2) {
		fmt.Println("bbw.Eval() :  ", bbw.Eval())
		t.Error("Expected 2, got ", bbw.Eval()) 
	}

}