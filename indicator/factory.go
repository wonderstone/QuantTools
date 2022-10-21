package indicator

type IIndicator interface {
	LoadData(data map[string]float64)
	Eval() float64
}

// factory pattern
func GetIndicator()
