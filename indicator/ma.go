package indicator

import (
	cb "github.com/wonderstone/QuantTools/indicator/tools"
)

// MA is the moving average indicator
type MA struct {
	ParSlice []int
	sum, lv  float64

	// info fields for indicator calculation
	InfoSlice []string
	DQ        *cb.Queue
}

// NewMA returns a new MA indicator
func NewMA(ParSlice []int, infoslice []string) *MA {
	return &MA{
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        cb.New(ParSlice[0]),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (m *MA) LoadData(data map[string]float64) {
	m.sum += data[m.InfoSlice[0]]
	if m.DQ.Full() {
		m.lv = m.DQ.Vals[m.DQ.End].(float64)
	}

	m.DQ.Enqueue(data[m.InfoSlice[0]])
	if m.DQ.Full() {
		m.sum -= m.lv
	}

}

// Eval evaluates the indicator
func (m *MA) EvalOld() float64 {
	var sum float64
	for _, v := range m.DQ.Values() {
		sum += v.(float64)
	}
	return sum / float64(m.DQ.Size())
}

func (m *MA) Eval() float64 {
	return m.sum / float64(m.DQ.Size())
}
