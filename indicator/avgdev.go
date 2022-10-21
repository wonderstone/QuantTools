package indicator

// use gods to generate the queue
import (
	"math"

	cb "github.com/emirpasic/gods/queues/circularbuffer"
)

// AvgDev is the AvgDev indicator
type AvgDev struct {
	ParSlice  []int
	InfoSlice []string
	DQ        *cb.Queue
	MA        *MA
}

// NewAvgDev returns a new AvgDev indicator
func NewAvgDev(ParSlice []int, infoslice []string) *AvgDev {
	// * 嵌套指标infoslice不同时，记得单独处理
	tmpma := NewMA(ParSlice, infoslice)
	return &AvgDev{
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        tmpma.DQ,
		MA:        tmpma,
	}
}

// LoadData loads 1 tick info datas into the indicator
func (a *AvgDev) LoadData(data map[string]float64) {
	a.DQ.Enqueue(data)
}

// Eval evaluates the indicator
func (a *AvgDev) Eval() float64 {
	var sum, devSum float64
	sum = a.MA.Eval()
	for _, v := range a.DQ.Values() {
		devSum += math.Abs(v.(map[string]float64)[a.InfoSlice[0]] - sum/float64(a.DQ.Size()))
	}

	return devSum / float64(a.DQ.Size())
}
