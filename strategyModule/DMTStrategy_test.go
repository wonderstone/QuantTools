package strategyModule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test DMTStrategy timecritic function
func TestDMTStrategyTimecritic(t *testing.T) {
	//time string
	timeString := "2020/12/8 11:30"
	// get the time value
	timeValue := getTimeValue(timeString)
	// test the value is correct
	assert.Equal(t, "11:30", timeValue)
}
