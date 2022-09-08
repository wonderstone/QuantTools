/* 1. 未完结，等待确认极速交易系统数据切片字段
2. 未完结，等待追加方法实现多态更新账户 */
package marketdata

import (
	"strings"

	"github.com/wonderstone/QuantTools/data"
)

// Bar类型不再区分股票和期货，统一为Bar。如果有特异数据，可以在indicator中进行处理

// Declaring futuresBar struct with key fields
type Bar struct {
	InstID  string
	BarTime string
	Open    float64
	Close   float64
	High    float64
	Low     float64
	Vol     float64
}

func NewBar(InstID string, BarTime string, Open float64, Close float64, High float64, Low float64, Vol float64) Bar {
	return Bar{InstID: InstID, BarTime: BarTime, Open: Open, Close: Close, High: High, Low: Low, Vol: Vol}
}

// 常规更新account的信息
func (b *Bar) GetUpdateInfo(tag string) (ui data.UpdateMI) {
	ui.UpdateTimeStamp = b.BarTime
	ui.InstID = b.InstID
	switch strings.ToLower(tag) {
	case "close", "c":
		ui.Value = b.Close
	case "high", "h":
		ui.Value = b.High
	case "low", "l":
		ui.Value = b.Low
	case "open", "o":
		ui.Value = b.Open
	case "vol", "v", "volume":
		ui.Value = b.Vol
	default:
		ui.Value = b.Close
	}

	return ui
}
