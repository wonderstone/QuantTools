package indicator

// use gods to generate the queue
import (
	cb "github.com/emirpasic/gods/queues/circularbuffer"
)

// Var 方差指标
type Var struct {
	// Period int
	ParSlice  []int
	InfoSlice []string
	DQ        *cb.Queue
	MA        *MA
}

// NewVar returns a new Variance indicator
func NewVar(ParSlice []int, infoslice []string) *Var {
	tmpma := NewMA(ParSlice, infoslice)
	return &Var{
		// Period: period,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        tmpma.DQ,
		MA:        tmpma,
	}
}

// LoadData loads 1 tick info datas into the indicator
func (v *Var) LoadData(Data []float64) {
	v.DQ.Enqueue(Data[0])
}

// Eval evaluates the indicator
func (v *Var) Eval() float64 {
	var sum float64
	avg := v.MA.Eval()
	for _, val := range v.DQ.Values() {
		sum += (val.(float64) - avg) * (val.(float64) - avg)
	}
	return sum / float64(v.DQ.Size())
}
