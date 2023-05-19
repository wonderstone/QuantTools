package strategyModule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test NewsortBuyFromConfig function
func TestNewsortBuyFromConfig(t *testing.T) {
	stg := NewSortBuyStrategyFromConfig("../config/Manual/", "BackTest01.yaml", "default", "StrategySB.yaml")
	assert.Equal(t, 5, stg.numHolding)

	// test getTime method
	// make a time type variable at 10:07:00.000 on 2018-01-01

}
