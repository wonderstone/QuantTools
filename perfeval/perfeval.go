package perfeval

import (
	"math"
	"sort"
	"strings"
	"time"

	"github.com/wonderstone/QuantTools/account"
)

type PerfEval struct {
	MktValSlice []account.MktValDataType //有机会换一下指针试试性能是否提升
	sorted      bool
}

func NewPerfEval() (PE *PerfEval) {
	PE = new(PerfEval)
	return
}

// type PerfEvalRes interface{
// 	Eval() float64
// }

// type PerfEvalResult struct {
// 	TotalReturn      float64
// 	AnnualizedReturn float64
// 	MaxDrawDown      float64
// 	SharpeRatio      float64
// }

func (p *PerfEval) CalcPerfEvalResult(tag string) float64 {
	switch tag {
	case "TR":
		return p.TotalReturn()
	case "AR":
		return p.AnnualizedReturn()
	case "MR":
		return p.AnnualizedReturn() / p.MaxDrawDown()
	case "SR":
		return p.SharpeRatio(0.03)
	default:
		return p.TotalReturn()
	}
}

func (PE *PerfEval) RateofReturns() (RoRs []float64) {
	RoRs = make([]float64, PE.Len()-1)
	for i := 1; i < PE.Len(); i++ {
		RoRs[i-1] = (PE.MktValSlice[i].MktVal / PE.MktValSlice[i-1].MktVal) - 1
	}
	return
}

func (PE *PerfEval) TotalReturn() (TR float64) {
	if !PE.sorted {
		PE.Sort()
	}
	return PE.MktValSlice[PE.Len()-1].MktVal / PE.MktValSlice[0].MktVal
}

func (PE *PerfEval) AnnualizedReturn() (AR float64) {
	//默认了日线级别 偷懒做法  后期有空精细化吧
	return math.Pow(PE.TotalReturn(), float64(252/PE.Len()))
}

func (PE *PerfEval) MaxDrawDown() (maxDrawDown float64) {
	if !PE.sorted {
		PE.Sort()
	}
	maxVal := 0.0
	for i := 0; i < PE.Len(); i++ {
		if PE.MktValSlice[i].MktVal > maxVal {
			maxVal = PE.MktValSlice[i].MktVal
		}
		drawDown := 1.0 - (PE.MktValSlice[i].MktVal / maxVal)
		if drawDown > 0 && drawDown > maxDrawDown {
			maxDrawDown = drawDown
		}

	}
	return
}

func (PE *PerfEval) SharpeRatio(Rf float64) (SR float64) {
	//默认了日线级别 偷懒做法  后期有空精细化吧
	std := Std(PE.RateofReturns(), 1)
	if std == 0 {
		return 0
	}
	return (PE.AnnualizedReturn() - Rf) / (math.Sqrt(252) * std)
}

// sort section
func (PE *PerfEval) Len() int {
	return len(PE.MktValSlice)
}

func (PE *PerfEval) Less(i, j int) bool {
	datei, erri := time.Parse("2006/1/2", strings.Fields(PE.MktValSlice[i].Time)[0]) // golang的时间format是数值固定的 与python不一致
	if erri != nil {
		panic(erri)
	}
	datej, errj := time.Parse("2006/1/2", strings.Fields(PE.MktValSlice[j].Time)[0])
	if errj != nil {
		panic(errj)
	}
	return datei.Before(datej)
}

func (PE *PerfEval) Swap(i, j int) {
	PE.MktValSlice[i], PE.MktValSlice[j] = PE.MktValSlice[j], PE.MktValSlice[i]
}

func (PE *PerfEval) Add(time string, mktval float64) {
	PE.MktValSlice = append(PE.MktValSlice, account.MktValDataType{Time: time, MktVal: mktval})
	PE.sorted = false
}

func (PE *PerfEval) Sort() {
	if !PE.sorted {
		sort.Sort(PE)
		PE.sorted = true
	}
}
