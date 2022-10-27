package matcher

import (
	"fmt"
	"testing"

	cp "github.com/wonderstone/QuantTools/contractproperty"
	"github.com/wonderstone/QuantTools/order"

	"github.com/stretchr/testify/assert"
)

// test MatchFuturesOrder
func TestMatchFuturesOrder(t *testing.T) {
	// new a fcp from instID
	fmt.Println("test MatchFuturesOrder")
	confName := "ContractProp.yaml"
	dir := "../config/Manual/"
	cpm := cp.NewCPMap(confName, dir)
	instID := "au2210"
	fcp := cp.SimpleNewFCPFromMap(cpm, instID)

	// fo := order.NewFuturesOrder("au2210", false, "20220515 13:35:27 500", 400.00, 2, order.Buy, order.Open, &fcp)
	fo := order.NewFuturesOrder("au2210", true, false, "20220515 13:35:27 500", 400.00, 2, "Buy", "Open", &fcp)
	simMatcher := NewSimpleMatcher(0.01, 1.0)

	expected := order.FuturesOrder{
		InstID:         "au2210",
		IsEligible:     true,
		IsExecuted:     true,                    //changed
		OrderTime:      "20220515 13:35:28 500", //changed
		OrderPrice:     401.90,                  //changed
		OrderNum:       2,
		OrderDirection: "Buy",
		OrderType:      "Open",
		FCP:            &fcp,
	}
	simMatcher.MatchFuturesOrder(&fo, 401.88, "20220515 13:35:28 500")
	assert.Equal(t, &expected, &fo, "MatchFuturesOrder不符合预期")

}

// benchmark MatchFuturesOrder @ 1.619 ns/op	       0 B/op	       0 allocs/op
// go test -bench=.
// BenchmarkMatchFuturesOrder @ 18.56 ns/op            0 B/op          0 allocs/op
// fuck!!! a little bit changes in fcp and simple matcher with apple M1
// you give me ten times slower result than the old one?
func BenchmarkMatchFuturesOrder(b *testing.B) {
	fcp := cp.NewFCP(1000, 0.02, 10, 10, 0.0, false, 2.0, 0.0, 2.0, 0.01)
	fo := order.NewFuturesOrder("au2210", true, true, "20220515 13:35:27 500", 400.00, 2, "Buy", "Open", &fcp)
	simMatcher := NewSimpleMatcher(0.02, 1.0)
	for i := 0; i < b.N; i++ {
		simMatcher.MatchFuturesOrder(&fo, 401.88, "20220515 13:35:28 500")
	}
}

// test MatchStockOrder
func TestMatchStockOrder(t *testing.T) {
	// new a fcp from instID
	fmt.Println("test MatchStockOrder")
	confName := "ContractProp"
	dir := "../config/Manual"
	cpm := cp.NewCPMap(confName, dir)
	instID := "SZ000058"
	scp := cp.SimpleNewSCPFromMap(cpm, instID)

	// so := order.NewStockOrder("SZ000058", false, "2022-05-10 14:52", 8.5, 2.0, order.Buy, &scp)
	so := order.NewStockOrder("SZ000058", true, false, "2022-05-10 14:52", 8.5, 2.0, "Buy", &scp)
	simMatcher := NewSimpleMatcher(0.01, 1.0)

	// expected := order.NewStockOrder("SZ000058", true, "2022-05-10 14:53", 8.51, 2.0, order.Buy, &scp)
	expected := order.NewStockOrder("SZ000058", true, true, "2022-05-10 14:53", 8.51, 2.0, "Buy", &scp)

	simMatcher.MatchStockOrder(&so, 8.50, "2022-05-10 14:53")
	assert.Equal(t, &expected, &so, "MatchStockOrder不符合预期")
}

// benchmark MatchStockOrder @ 1.516 ns/op	       0 B/op	       0 allocs/op
// go test -bench=.
// BenchmarkMatchStockOrder-8@ 18.51 ns/op	       0 B/op	       0 allocs/op
// again! ten times slower than the old one
func BenchmarkMatchStockOrder(b *testing.B) {
	scp := cp.NewSCP(100, 0.00001, 0.001, 0.0000687)
	// so := order.NewStockOrder("SZ000058", false, "2022-05-10 14:52", 8.5, 2.0, order.Buy, &scp)
	so := order.NewStockOrder("SZ000058", true, false, "2022-05-10 14:52", 8.5, 2.0, "Buy", &scp)

	simMatcher := NewSimpleMatcher(0.01, 1.0)
	for i := 0; i < b.N; i++ {
		simMatcher.MatchStockOrder(&so, 8.50, "2022-05-10 14:53")
	}
}
