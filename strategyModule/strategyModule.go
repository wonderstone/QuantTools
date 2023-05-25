package strategyModule

import (
	"math"

	"github.com/wonderstone/QuantTools/account/virtualaccount"
	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/dataprocessor"
	"github.com/wonderstone/QuantTools/order"
)

const debug = true

type OrderResult struct {
	StockOrderS   []order.StockOrder
	FuturesOrderS []order.FuturesOrder
	// IsExecuted    bool
}

func NewOrderResult() OrderResult {
	return OrderResult{
		StockOrderS:   make([]order.StockOrder, 0),
		FuturesOrderS: make([]order.FuturesOrder, 0),
		// IsExecuted:    false,
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

// update the stock account
func UpdateAcct(bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct) {
	for _, v := range bc.Stockdata {
		vAcct.SAcct.ActOnUpdateMI(v.BarTime, v.InstID, v.IndiDataMap["Close"])
	}
	for _, v := range bc.Futuresdata {
		vAcct.FAcct.ActOnUpdateMI(v.BarTime, v.InstID, v.IndiDataMap["Close"])
	}
}

// net the orders
func NetSOrders(os []order.StockOrder) []order.StockOrder {
	// 注意同一时刻是不需要check ordertime的
	var res []order.StockOrder
	// orders with the same id will be netted
	// * this is a map for the orders
	// * key: orderID
	// * value: order
	tmpMap := make(map[string]order.StockOrder)
	for _, v := range os {
		if _, ok := tmpMap[v.InstID]; ok {
			// temp the order in map
			tmpOrder := tmpMap[v.InstID]
			// chect the direction
			if v.OrderDirection == tmpMap[v.InstID].OrderDirection {
				// add the order num
				tmpOrder.OrderNum += v.OrderNum
				// update the map
				tmpMap[v.InstID] = tmpOrder
			} else {
				if v.OrderNum > tmpOrder.OrderNum {
					// update the order num
					tmpOrder.OrderNum = v.OrderNum - tmpOrder.OrderNum
					// update the order direction
					tmpOrder.OrderDirection = v.OrderDirection
					// update the map
					tmpMap[v.InstID] = tmpOrder
				} else if v.OrderNum == tmpMap[v.InstID].OrderNum {
					// delete the order in map
					delete(tmpMap, v.InstID)
				} else {
					// update the order num
					tmpOrder.OrderNum = tmpOrder.OrderNum - v.OrderNum
					// update the map
					tmpMap[v.InstID] = tmpOrder
				}
			}
		} else {
			tmpMap[v.InstID] = v
		}
	}
	for _, v := range tmpMap {
		res = append(res, v)
	}
	return res
}

// 此处设计强制要求形式上有GEP和Manual两类，如果真不想写，对应的地方留空即可
type IStrategy interface {
	ActOnData(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap, Eval func([]float64) []float64) (orderRes OrderResult)
	ActOnDataMAN(datetime string, bc *dataprocessor.BarC, vAcct *virtualaccount.VAcct, CPMap cp.CPMap) (orderRes OrderResult)
}

// something like simple factory pattern
func GetStrategy(dir string, BTConfile string, sec string, StgConfile string, tag string) IStrategy {
	switch tag {
	case "simple":
		return NewSimpleStrategyFromConfig(dir, BTConfile, sec, StgConfile)
	case "DMT":
		s := NewDMTStrategyFromConfig(dir, BTConfile, sec, StgConfile)
		return &s
	case "T0":
		t := NewST0StrategyFromConfig(dir, BTConfile, sec, StgConfile)
		return &t
	case "SortBuy":
		sb := NewSortBuyStrategyFromConfig(dir, BTConfile, sec, StgConfile)
		return &sb
	default:
		return NewSimpleStrategyFromConfig(dir, BTConfile, sec, StgConfile)
	}
}
