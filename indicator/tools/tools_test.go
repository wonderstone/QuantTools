package tools

import (
	"fmt"

	// mat "gonum.org/v1/gonum/mat"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestEvalSum1(t *testing.T) {
	array := []float64{1, 2, 3, 4, 5}

	if Sum1(array) != 15.0 {
		fmt.Println("Sum1() :  ", Sum1(array))
		t.Error("Expected 15.0, got ", Sum1(array))
	}
}

func TestEvalSumN(t *testing.T) {
	array := []float64{1, 2, 3, 4, 5}
	sumN := []float64{6, 9, 12}
	b := true
	for i, v := range sumN { //数组比较
		if v != SumN(array, 3)[i] {
			b = false
		}
	}
	if !b {
		fmt.Printf("SumN() : %v \n", SumN(array, 3))
		t.Error("Expected 15.0, got ", SumN(array, 3))
	}
}

func TestEvalAvg(t *testing.T) {
	array := []float64{1, 2, 3, 4, 5}
	sma := Sma(array, 3, 2)
	fmt.Printf("%v", sma) //结果[nan,nan,2,3,4]
	if Avg(array) != 3 {
		fmt.Println("Ma.Eval() :  ", Avg(array))
		t.Error("Expected 3, got ", Avg(array))
	}
}

func TestEvalMaN(t *testing.T) {
	array := []float64{1, 2, 3, 4, 5}
	maN := []float64{2, 3, 4}
	b := true
	for i, v := range maN { //数组比较
		if v != AvgN(array, 3)[i] {
			b = false
		}
	}
	if !b {
		fmt.Printf("AvgN() : %v \n", AvgN(array, 3))
		t.Error("Expected 15.0, got ", AvgN(array, 3))
	}
}

func TestEvalStd(t *testing.T) {
	array := []float64{1, 2, 3, 4, 5}
	if Std(array) != math.Sqrt(2.5) {
		fmt.Println("m.Eval() :  ", Std(array))
		t.Error("Expected 5.0, got ", Std(array))
	}
}

func TestEvalAvgDev(t *testing.T) {
	array := []float64{1, 2, 3, 4, 5}
	if AvgDev(array) == 1.5 {
		fmt.Println("AvgDev() :  ", AvgDev(array))
		t.Error("Expected 5.0, got ", AvgDev(array))
	}
}

// func TestEvalPushBarList(t *testing.T) {
// 	BarList := MakeBarList(100)
// 	queue := cb.New(300)
// 	PushBarList(queue, BarList)
// 	if len(BarList) == 100 {
// 		fmt.Println("len(BarList) :  ", len(BarList))
// 		t.Error("Expected 100, got ", len(BarList))
// 	}
// }

func TestPerformance(t *testing.T) { //用于性能测试  拟废除
	//常规来说界面显示300个K线,最多4000个
	begin := time.Now()

	var rand20, rand100 []float64 //指标参数周期多以20为主,大的100左右
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 20; i++ {
		rand20 = append(rand20, float64(rand.Int31n(200))+10)
	}
	for i := 0; i < 10000; i++ {
		rand100 = append(rand100, float64(rand.Int31n(200))+10)
	}
	for i := 0; i < 100000000; i++ {
		tmp := rand100
		tmp[0] = 0.0
	}
	end := time.Now()

	println(end.Sub(begin).Seconds())
}

// func TestGonumVectorOperate(t *testing.T) {
// 	u := mat.NewVecDense(3, []float64{1, 2, 3})
// 	v := mat.NewVecDense(3, []float64{4, 5, 6})
// 	w := mat.NewVecDense(3, nil)
// 	w.AddVec(u, v)
// 	fmt.Println("u + v: ", w)
// 	// Add u + alpha * v for some scalar alpha
// 	w.AddScaledVec(u, 2, v)
// 	fmt.Println("u + 2 * v: ", w)
// 	// Subtract v from u
// 	w.SubVec(u, v)
// 	fmt.Println("v - u: ", w)
// 	// Scale u by alpha
// 	w.ScaleVec(23, u)
// 	fmt.Println("u * 23: ", w)
// 	// Compute the dot product of u and v
// 	// Since float64’s don’t have a dot method, this is not done
// 	//inplace
// 	d := mat.Dot(u, v)
// 	fmt.Println("u dot v: ", d)
// 	// element-wise product
// 	w.MulElemVec(u, v)
// 	fmt.Println("u element-wise product v: ", w)
// 	// Find length of v
// 	l := v.Len()

// 	fmt.Println("Length of v: ", l)

// }

func TestVec(t *testing.T) {
	a := []float64{1, 2, 3, 4, 5}
	b := []float64{-1, 2, 3, 4, 5}
	fmt.Printf("a+b=%v\n", VecAdd(a, b))
	fmt.Printf("a-b=%v\n", VecSub(a, b))
	fmt.Printf("a*b=%v\n", VecMult(a, b))
	fmt.Printf("a/b=%v\n", VecDiv(a, b))
	fmt.Printf("max(a,b)=%v\n", VecMax(a, b))
	fmt.Printf("min(a,b)=%v\n", VecMin(a, b))
	fmt.Printf("repeat1D(3,2)=%v\n", Repeat1D(3, 2))
	fmt.Printf("VecRef1(a)=%v\n", VecRef1(a))
	fmt.Printf("VecAbs(b)=%v\n", VecAbs(b))

}
