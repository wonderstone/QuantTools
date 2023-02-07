package strategyModule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test NewT0StrategyFromConfig function
func TestNewT0StrategyFromConfig(t *testing.T) {
	stg := NewST0StrategyFromConfig("../config/Manual/", "BackTest01.yaml", "Default", "StrategyT0.yaml")
	assert.Equal(t, 4, stg.Tlimit)

	// test getTime method
	// make a time type variable at 10:07:00.000 on 2018-01-01

}

// time string check
func TestTimeCheck(t *testing.T) {
	assert.Greater(t, "10:07:00.000", "09:06:50.999")
}
