package dataprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wonderstone/QuantTools/indicator"
)

// test CsvProcess
func TestCsvProcess(t *testing.T) {

	csvpath := "../tmpdata/stockdata/1min/sz000058.csv"
	indis := []indicator.IndiInfo{
		{Name: "MA3", IndiType: "MA", ParSlice: []int{3}, InfoSlice: []string{"close"}},
		{Name: "Var3", IndiType: "Var", ParSlice: []int{3}, InfoSlice: []string{"close"}},
	}
	isok, _ := CsvProcess(csvpath, indis)
	assert.Equal(t, isok, true)

}
