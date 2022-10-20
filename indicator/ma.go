package indicator

// use gods to generate the queue
import (
	cb "github.com/emirpasic/gods/queues/circularbuffer"
)

// MA is the moving average indicator
type MA struct {
	period int
	// info fields for indicator calculation
	InfoSlice []string
	DQ        *cb.Queue
}

// NewMA returns a new MA indicator
func NewMA(period int, infoslice []string) *MA {
	return &MA{
		period:    period,
		InfoSlice: infoslice,
		DQ:        cb.New(period),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (m *MA) LoadData(data map[string]float64) {
	m.DQ.Enqueue(data)
}

// Eval evaluates the indicator
func (m *MA) Eval() float64 {
	var sum float64
	for _, v := range m.DQ.Values() {
		sum += v.(map[string]float64)[m.InfoSlice[0]]
	}
	return sum / float64(m.DQ.Size())
}
