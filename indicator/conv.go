package indicator

// use gods to generate the queue
import (
	cb "github.com/emirpasic/gods/queues/circularbuffer"
)

// Conv is the Conv indicator
type Conv struct {
	ParSlice  []int
	InfoSlice []string
	DQS       *cb.Queue
	DQI       *cb.Queue
	Ma        *MA
}

// NewConv returns a new Conv indicator
func NewConv(ParSlice []int, infoslice []string) *Conv {
	return &Conv{
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQS:       cb.New(ParSlice[0]),
		DQI:       cb.New(ParSlice[0]),
		Ma:        NewMA(ParSlice, infoslice),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (c *Conv) LoadData(Data []float64) {
	c.DQS.Enqueue(Data[0])
	c.DQI.Enqueue(Data[1])
}

// Eval evaluates the indicator
func (v *Conv) Eval() float64 {
	var sum float64
	v.Ma.DQ = v.DQS
	avgStock := v.Ma.Eval()
	v.Ma.DQ = v.DQI
	avgIndex := v.Ma.Eval()
	for i := range v.DQS.Values() {

		sum += (v.DQI.Values()[i].(float64) - avgStock) * (v.DQI.Values()[i].(float64) - avgIndex)
	}
	return sum / float64(v.DQI.Size())
}
