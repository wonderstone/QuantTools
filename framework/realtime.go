package framework

import (
	"fmt"
	"time"

	"github.com/wonderstone/QuantTools/order"
	"github.com/wonderstone/QuantTools/strategyModule"
)

func SendOrders(info map[string]interface{}, orders strategyModule.OrderResult) {
	if orders.StockOrderS != nil {
		SendSO(orders.StockOrderS, info)
	}
	if orders.FuturesOrderS != nil {
		SendFO(orders.FuturesOrderS, info)
	}
}

func SendSO(sos []order.StockOrder, info map[string]interface{}) {
	// send real stock order to server
	// get real so
	for _, so := range sos {
		// rsos = append(rsos, order.GetStockOrder(so, info))
		tmpRSO := order.GetStockOrder(so, info)
		// *************************
		// send real order to server
		// *************************
		fmt.Println(tmpRSO)
	}
}

func SendFO(fos []order.FuturesOrder, info map[string]interface{}) {
	// send real futures order to server
	// get real fo
	for _, fo := range fos {
		tmpRFO := order.GetFuturesOrder(fo, info)
		// *************************
		// send real order to server
		// *************************
		fmt.Println(tmpRFO)
	}

}

// + realtime job would use this function
func getRealTimeStamp() string {
	// get current time stamp
	currentTime := time.Now()
	return currentTime.Format("2006/1/2 15:04:05")
}
