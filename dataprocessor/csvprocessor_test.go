package dataprocessor

import (
	"encoding/csv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// test CsvProcess
func TestCsvProcess(t *testing.T) {

	csvpath := "../tmpdata/stockdata/1min/SZ000058.csv"

	isok, err := CsvProcess(csvpath)
	assert.Equal(t, isok, true)

}
