package strategyModule

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test ContainNaN
func TestContainNaN(t *testing.T) {
	m := map[string]float64{"a": 1.0, "b": math.NaN(), "c": 3}
	assert.Equal(t, ContainNaN(m), true)
}
