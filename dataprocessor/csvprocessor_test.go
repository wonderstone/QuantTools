package dataprocessor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// test CsvProcess
func TestCsvProcess(t *testing.T) {

	csvpath := "../tmpdata/stockdata/1min/SZ000058.csv"

	isok, _ := CsvProcess(csvpath)
	assert.Equal(t, isok, true)

}
