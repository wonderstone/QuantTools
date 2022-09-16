package dataprocessor

import (
	"fmt"

	"github.com/wonderstone/QuantTools/order"
)

type QuotePublisher interface {
	Pub(info map[string]string, bc chan<- BarC, mc chan<- map[string]float64)
}

type QuoteSubscriber interface {
	Sub(info map[string]string, bc <-chan BarC, mc <-chan map[string]float64)
	// sub should check if sos or fos is nil and SendSO and SendFO
}

func SendSO(sos []order.StockOrder, info map[string]string) {
	// send real stock order to server
	// get real so
	for _, so := range sos {
		// rsos = append(rsos, order.GetStockOrder(so, info))
		tmpRSO := order.GetStockOrder(so, info)
		// send real order to server
		fmt.Println(tmpRSO)
	}
}

func SendFO(fos []order.FuturesOrder, info map[string]string) {
	// send real futures order to server
	// get real fo
	for _, fo := range fos {
		tmpRFO := order.GetFuturesOrder(fo, info)
		// send real order to server
		fmt.Println(tmpRFO)
	}

}

func RealTimeProcess(qp QuotePublisher, qs QuoteSubscriber, info map[string]string, bc chan BarC, mc chan map[string]float64) {
	// start a goroutine to get data and publish to bc and mc
	go qp.Pub(info, bc, mc)
	// run the strategy and give the order
	qs.Sub(info, bc, mc)
}
