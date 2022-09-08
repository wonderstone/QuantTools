package marketdata

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wonderstone/QuantTools/data"
)

//test NewFuturesBar
func TestNewFuturesBar(t *testing.T) {
	expected := Bar{
		InstID:  "cu",
		BarTime: "2022-05-10 12:12:12 500",
		Open:    3459.3,
		Close:   3459.4,
		High:    3459.5,
		Low:     3459.2,
		Vol:     123456789.0,
	}
	actual := NewBar("cu", "2022-05-10 12:12:12 500", 3459.3, 3459.4, 3459.5, 3459.2, 123456789.0)
	assert.Equal(t, expected, actual, fmt.Sprintf("NewFuturesBar()=%v,expected=%v", actual, expected))

}

// test Futures GetUpdateInfo
func TestGetUpdateInfo_F(t *testing.T) {
	fb := Bar{
		InstID:  "cu",
		BarTime: "2022-05-10 12:12:12 500",
		Open:    3459.3,
		Close:   3459.4,
		High:    3459.5,
		Low:     3459.2,
		Vol:     123456789.0,
	}
	expected := data.UpdateMI{
		UpdateTimeStamp: fb.BarTime,
		InstID:          fb.InstID,
		Value:           fb.Close,
	}
	actual := fb.GetUpdateInfo("close")
	assert.Equal(t, expected, actual, fmt.Sprintf("GetUpdateInfo()=%v,expected=%v", actual, expected))
}

//test NewStockBar
func TestNewStockBar(t *testing.T) {
	expected := Bar{
		InstID:  "SZ000058",
		BarTime: "2022-05-10 14:52",
		Open:    8.6,
		Close:   8.5,
		High:    8.7,
		Low:     8.4,
		Vol:     123456789.0,
	}
	actual := NewBar("SZ000058", "2022-05-10 14:52", 8.6, 8.5, 8.7, 8.4, 123456789.0)
	assert.Equal(t, expected, actual, fmt.Sprintf("NewStockBar()=%v,expected=%v", actual, expected))
}

// test stock GetUpdateInfo
func TestGetUpdateInfo_S(t *testing.T) {
	sb := Bar{
		InstID:  "SZ000058",
		BarTime: "2022-05-10 14:52",
		Open:    8.6,
		Close:   8.5,
		High:    8.7,
		Low:     8.4,
		Vol:     123456789.0,
	}
	expected := data.UpdateMI{
		UpdateTimeStamp: sb.BarTime,
		InstID:          sb.InstID,
		Value:           sb.Close,
	}
	actual := sb.GetUpdateInfo("close")
	assert.Equal(t, expected, actual, fmt.Sprintf("GetUpdateInfo()=%v,expected=%v", actual, expected))
}
