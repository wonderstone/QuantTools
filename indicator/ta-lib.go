package indicator

// use gods to generate the queue
import (
	cb "github.com/emirpasic/gods/queues/circularbuffer"
)

// MA is the moving average indicator
type MA struct {
	period int
	DQ     *cb.Queue
}

// NewMA returns a new MA indicator
func NewMA(period int) *MA {
	return &MA{
		period: period,
		DQ:     cb.New(period),
	}
}

// Eval evaluates the indicator
func (m *MA) Eval() float64 {
	var sum float64
	for _, v := range m.DQ.Values() {

		sum += v.(float64)
	}
	return sum / float64(m.DQ.Size())
}
