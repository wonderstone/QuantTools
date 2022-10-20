package indicator

// use gods to generate the queue
import (
	"math"

	cb "github.com/emirpasic/gods/queues/circularbuffer"
)

// AvgDev is the AvgDev indicator
type AvgDev struct {
	period    int
	InfoSlice []string
	DQ        *cb.Queue
	ma        *MA
}

// NewAvgDev returns a new AvgDev indicator
func NewAvgDev(period int, infoslice []string) *AvgDev {
	// * 嵌套指标infoslice不同时，记得单独处理
	tmpma := NewMA(period, infoslice)
	return &AvgDev{
		period:    period,
		InfoSlice: infoslice,
		DQ:        tmpma.DQ,
		ma:        tmpma,
	}
}

// LoadData loads 1 tick info datas into the indicator
func (a *AvgDev) LoadData(data map[string]float64) {
	a.DQ.Enqueue(data)
}

// Eval evaluates the indicator
func (m *AvgDev) Eval() float64 {
	var sum, devSum float64
	sum = m.ma.Eval()
	for _, v := range m.DQ.Values() {
		devSum += math.Abs(v.(map[string]float64)[m.InfoSlice[0]] - sum/float64(m.DQ.Size()))
	}

	return devSum / float64(m.DQ.Size())
}
