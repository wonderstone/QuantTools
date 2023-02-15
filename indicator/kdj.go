package indicator

// use gods to generate the DQ
import (
	"github.com/wonderstone/QuantTools/indicator/tools"
	"gonum.org/v1/gonum/floats"
)

// KDJ 顺势指标/商品路径指标
// https://max.book118.com/html/2017/1204/142703047.shtm
type KDJ struct {
	Name      string
	ParSlice  []int    //period N M1 M2
	InfoSlice []string // C H L
	DQ        *tools.Queue
}

// NewKDJ  returns a new KDJ indicator
func NewKDJ(Name string, ParSlice []int, infoslice []string) *KDJ {
	return &KDJ{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        tools.New(ParSlice[0]),
	}
}

// GetName returns the name of the indicator
func (k *KDJ) GetName() string {
	return k.Name
}

// LoadData loads 1 tick info datas into the indicator
func (k *KDJ) LoadData(data map[string]float64) {
	k.DQ.Enqueue(data)
}

// Eval evaluates the indicator
func (k *KDJ) Eval() (float64, float64, float64) {
	var rsv []float64
	var low, high []float64
	for _, v := range k.DQ.Values() {
		bar := v.(map[string]float64)
		low = append(low, bar[k.InfoSlice[2]])
		high = append(high, bar[k.InfoSlice[1]])
		begin := len(low) - k.ParSlice[1]
		if begin < 0 {
			begin = 0
		}
		rsv = append(rsv, (bar[k.InfoSlice[0]]-floats.Min(low[begin:]))/(floats.Max(high[begin:])-floats.Min(low[begin:])))
	}
	K := tools.Sma(rsv, k.ParSlice[2], 1)
	D := tools.Sma(K[k.ParSlice[2]-1:], k.ParSlice[3], 1)
	J := 3*K[len(K)-1] - 2*D[len(D)-1]
	return K[len(K)-1], D[len(D)-1], J
}
