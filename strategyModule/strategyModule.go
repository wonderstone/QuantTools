package strategyModule

import (
	"math"

	"github.com/wonderstone/QuantTools/account/virtualaccount"
	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/dataprocessor"
	"github.com/wonderstone/QuantTools/order"
)

type OrderResult struct {
	StockOrderS   []order.StockOrder
	FuturesOrderS []order.FuturesOrder
	IsExecuted    bool
}

func NewOrderResult() OrderResult {
	return OrderResult{
		StockOrderS:   make([]order.StockOrder, 0),
		FuturesOrderS: make([]order.FuturesOrder, 0),
		IsExecuted:    false,
	}
}

// 此处是为了停盘数据处理设定的规则相检查用的
func ContainNaN(m map[string]float64) bool {
	for _, x := range m {
		if math.IsNaN(x) {
			return true
		}
	}
	return false
}

// 此处设计强制要求形式上有GEP和Manual两类，如果真不想写，对应的地方留空即可
type IStrategy interface {
	ActOnData(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap, Eval func([]float64) []float64) (orderRes OrderResult)
	ActOnDataMAN(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap) (orderRes OrderResult)
}
