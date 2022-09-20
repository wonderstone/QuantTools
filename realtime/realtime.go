package realtime

import (
	"fmt"
	"time"

	"github.com/wonderstone/QuantTools/account/virtualaccount"
	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/dataprocessor"
	"github.com/wonderstone/QuantTools/order"
	"github.com/wonderstone/QuantTools/strategyModule"
)

type QuotePublisher interface {
	Pub(bc chan<- dataprocessor.BarC, mc chan<- map[string]map[string]float64)
}

type QuoteSubscriber interface {
	// receive data from bc and mc and run strategy
	Sub(bc <-chan dataprocessor.BarC, mc <-chan map[string]map[string]float64, strategymodule strategyModule.IStrategy, CPMap cp.CPMap, Eval func([]float64) []float64, mode string)
	// sub should check if sos or fos is nil and SendSO and SendFO
}

func SendOrders(info map[string]string, orders strategyModule.OrderResult) {
	if orders.StockOrderS != nil {
		SendSO(orders.StockOrderS, info)
	}
	if orders.FuturesOrderS != nil {
		SendFO(orders.FuturesOrderS, info)
	}
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

func RealTimeProcess(qp QuotePublisher, qs QuoteSubscriber, bc chan dataprocessor.BarC, mc chan map[string]map[string]float64, strategymodule strategyModule.IStrategy, CPMap cp.CPMap, Eval func([]float64) []float64, mode string) {
	// start a goroutine to get data and publish to bc and mc
	go qp.Pub(bc, mc)
	// run the strategy and give the order
	qs.Sub(bc, mc, strategymodule, CPMap, Eval, mode)
}

type FakeQPublisher struct {
	Info       map[string]string
	InstSIDS   []string
	IndiSNames []string
	InstFIDS   []string
	IndiFNames []string
}

func NewFakeQPublisher(info map[string]string, instSIDS []string, indiSNames []string, instFIDS []string, indiFNames []string) *FakeQPublisher {
	return &FakeQPublisher{
		Info:       info,
		InstSIDS:   instSIDS,
		IndiSNames: indiSNames,
		InstFIDS:   instFIDS,
		IndiFNames: indiFNames,
	}
}

type FakeQSubscriber struct {
	Info map[string]string
}

func NewFakeQSubscriber(info map[string]string) *FakeQSubscriber {
	return &FakeQSubscriber{
		Info: info,
	}
}

func (fqp *FakeQPublisher) Pub(bc chan<- dataprocessor.BarC, mc chan<- map[string]map[string]float64) {
	// 1. get data from server
	fmt.Println("log in and get data from server", fqp.Info)
	// fake some data from csv file,replace when real data is available

	// 2. publish to bc and mc
}

func (fqs *FakeQSubscriber) Sub(bc <-chan dataprocessor.BarC, mc <-chan map[string]map[string]float64, va *virtualaccount.VAcct, strategymodule strategyModule.IStrategy, CPMap cp.CPMap, Eval func([]float64) []float64, mode string) {
	// 1. get data from bc and mc
	go ActOnCM(va, mc)
	ActOnBarC(fqs.Info, va, bc, strategymodule, CPMap, Eval, mode)
	// 2. run strategy
	// 3. get orders
	// 4. send orders to server
}

func getRealTimeStamp() string {
	// get current time stamp
	currentTime := time.Now()
	return currentTime.Format("2006/1/2 15:04:05")
}

func ActOnBarC(info map[string]string, va *virtualaccount.VAcct, bc <-chan dataprocessor.BarC, strategymodule strategyModule.IStrategy, CPMap cp.CPMap, Eval func([]float64) []float64, mode string) {
	// 1. get data from bc
	for data := range bc {
		// 2. run strategy
		// ts:=getRealTimeStamp()
		ts, err := data.GetTimeStamp()
		if err != nil {
			switch mode {
			case "GEP":
				tmpOrderRes := strategymodule.ActOnData(ts, &data, va, CPMap, Eval)
				SendOrders(info, tmpOrderRes)
			case "Manual":
				tmpOrderRes := strategymodule.ActOnDataMAN(ts, &data, va, CPMap)
				SendOrders(info, tmpOrderRes)
			default:
				panic("mode is not defined")
			}
		}
	}
}
func ActOnCM(va *virtualaccount.VAcct, mc <-chan map[string]map[string]float64) {
	// 1. get data from mc and va update with the data
	for data := range mc {
		for timestamp, kv := range data {
			for k, v := range kv {
				va.FAcct.ActOnMTM(timestamp, k, v)
			}
		}
	}
}
