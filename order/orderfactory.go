package order

// **********************************************
// 这个部分只是功能设定实现，在具体交投场景下需要做出调整
// 因Quanttools核心只完成订单核心部分，对于真实交投的用户名、密码、IP和Port等信息实际上都是此处添加
// **********************************************
type RealStockOrder struct {
	InstID         string
	OrderTime      string
	OrderPrice     float64
	OrderNum       float64
	OrderDirection string
	// OrderDirection OrderDir
	User     string
	Password string
	// IP             string
	// Port           string
}

type RealFuturesOrder struct {
	InstID         string
	IsExecuted     bool
	OrderTime      string
	OrderPrice     float64
	OrderNum       float64
	OrderDirection string
	OrderType      string
	// OrderDirection OrderDir
	// OrderType      FuturesOrderTYP
	User     string
	Password string
	// IP             string
	// Port           string
}

func GetStockOrder(so StockOrder, info map[string]string) RealStockOrder {
	return RealStockOrder{
		InstID:         so.InstID,
		OrderTime:      so.OrderTime,
		OrderPrice:     so.OrderPrice,
		OrderNum:       so.OrderNum,
		OrderDirection: so.OrderDirection,
		User:           info["user"],
		Password:       info["password"],
	}
}

func GetFuturesOrder(fo FuturesOrder, info map[string]string) RealFuturesOrder {
	return RealFuturesOrder{
		InstID:         fo.InstID,
		IsExecuted:     fo.IsExecuted,
		OrderTime:      fo.OrderTime,
		OrderPrice:     fo.OrderPrice,
		OrderNum:       fo.OrderNum,
		OrderDirection: fo.OrderDirection,
		OrderType:      fo.OrderType,
		User:           info["user"],
		Password:       info["password"],
	}
}
