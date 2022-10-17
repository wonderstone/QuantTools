package indicator

// use gods to generate the queue
import (
	cb "github.com/emirpasic/gods/queues/circularbuffer"
)

// MA is the moving average indicator
type MA struct {
	period           int
	DefaultInfoSlice []string
	DQ               *cb.Queue
}

// NewMA returns a new MA indicator
func NewMA(period int) *MA {
	return &MA{
		period:           period,
		DefaultInfoSlice: []string{"Close"},
		DQ:               cb.New(period),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (m *MA) LoadData(data []float64) {
	m.DQ.Enqueue(data[0])
}

// Eval evaluates the indicator
func (m *MA) Eval() float64 {
	var sum float64
	for _, v := range m.DQ.Values() {

		sum += v.(float64)
	}
	return sum / float64(m.DQ.Size())
}
