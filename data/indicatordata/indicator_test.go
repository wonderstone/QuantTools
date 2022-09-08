// test indicator.go
package indicatordata

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//test NewIndiData
func TestNewIndiData(t *testing.T) {
	expected := IndiData{
		IndiName: "MA",
		InstID:   "cu",
		IndiTime: "2022-05-10 12:12:12 500",
		Value:    3459.3,
	}
	actual := NewIndiData("MA", "cu", "2022-05-10 12:12:12 500", 3459.3)
	assert.Equal(t, expected, actual, fmt.Sprintf("NewIndiData()=%v,expected=%v", actual, expected))
}

//benchmark NewIndiData @ 0.3800 ns/op	       0 B/op	       0 allocs/op
func BenchmarkNewIndiData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewIndiData("MA", "cu", "2022-05-10 12:12:12 500", 3459.3)
	}
}
