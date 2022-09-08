package indicatordata

// Declaring indiData struct with key fields
type IndiData struct {
	IndiName string
	InstID   string
	IndiTime string
	Value    float64
}

func NewIndiData(IndiName string, InstID string, IndiTime string, Value float64) IndiData {
	return IndiData{IndiName: IndiName, InstID: InstID, IndiTime: IndiTime, Value: Value}
}
