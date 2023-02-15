package tools

import (
	"math"
	// "math/rand"
)

// CompareFloat 比较浮点数的大小
func CompareFloat(a, b float64) bool {
	return math.Abs(a-b) < 0.0000001
}

// VecAdd 向量加法
func VecAdd(a, b []float64) []float64 {
	if len(a) != len(b) {
		return nil
	}
	var c []float64
	c = append(c, a...)

	for i := range a {
		c[i] = a[i] + b[i]
	}
	return c
}

// VecSub   向量减法
func VecSub(a, b []float64) []float64 {
	if len(a) != len(b) {
		return nil
	}
	var c []float64
	c = append(c, a...)
	for i := range a {
		c[i] = a[i] - b[i]
	}
	return c
}

// VecMult 向量乘法
func VecMult(a, b []float64) []float64 {
	if len(a) != len(b) {
		return nil
	}
	var c []float64
	c = append(c, a...)
	for i := range a {
		c[i] = a[i] * b[i]
	}
	return c
}

// VecDiv 向量除法
func VecDiv(a, b []float64) []float64 {
	if len(a) != len(b) {
		return nil
	}

	var c []float64
	c = append(c, a...)
	for i := range a {
		c[i] = a[i] / b[i]
	}
	return c
}

// Repeat1D 生成由一个元素重复形成的向量
func Repeat1D(a float64, n int) []float64 {
	var c []float64
	for i := 0; i < n; i++ {
		c = append(c, a)
	}
	return c
}

func VecAbs(a []float64) []float64 {
	for i := range a {
		if a[i] < 0 {
			a[i] = -a[i]
		}
	}
	return a
}

// VecMax 向量逐元素取较大者
func VecMax(a, b []float64) []float64 {
	if len(a) != len(b) {
		return nil
	}
	var c []float64
	for i := 0; i < len(a); i++ {
		tmp := 0.0
		if a[i] > b[i] {
			tmp = a[i]
		} else {
			tmp = b[i]
		}
		c = append(c, tmp)
	}
	return c
}

// VecMin 向量逐元素取较小者
func VecMin(a, b []float64) []float64 {
	if len(a) != len(b) {
		return nil
	}
	var c []float64
	for i := 0; i < len(a); i++ {
		tmp := 0.0
		if a[i] < b[i] {
			tmp = a[i]
		} else {
			tmp = b[i]
		}
		c = append(c, tmp)
	}
	return c
}
func VecRef1(a []float64) []float64 {
	var b []float64
	b = append(b, a[0])
	for i := 0; i < len(a)-1; i++ {
		b = append(b, a[i])
	}
	return b
}

// Sum1 常规的对数组求和
func Sum1(array []float64) float64 {
	var sum float64
	for _, v := range array {
		sum += v
	}
	return sum
}

// SumN  对该数组能求sum(n数之和)的数都求一遍
func SumN(array []float64, n int) []float64 {
	var sumN []float64
	for i := n - 1; i < cap(array); i++ {
		sum := Sum1(array[i-(n-1) : i+1])
		sumN = append(sumN, sum)
	}
	return sumN
}

// Avg 对该数组能求ma(n数之均值)的数都求一遍
func Avg(array []float64) float64 {

	return AvgN(array, len(array))[0]
}

// AvgN 对该数组能求ma(n数之均值)的数都求一遍
func AvgN(array []float64, n int) []float64 {
	sumN := SumN(array, n)
	for i := range sumN {
		sumN[i] = sumN[i] / float64(n)
	}
	return sumN
}

// Std 标准差
func Std(array []float64) float64 {
	mean := Sum1(array) / float64(len(array))
	var n int = len(array)
	res := 0.0
	for i := 0; i < n; i++ {
		res += (array[i] - mean) * (array[i] - mean)
	}
	return math.Sqrt(res / float64(n-1))
}

// AvgDev 返回平均绝对偏差
func AvgDev(array []float64) float64 {
	mean := Sum1(array) / float64(len(array))
	devSum := 0.0
	for _, v := range array {
		devSum += math.Abs(v - mean)
	}
	return devSum / float64(len(array))
}

// func MakeBarList(n int) []bars.Bar {
// 	//todo
// 	//该函数生成K线模拟数据,要求形似真实K线
// 	//思路为三时间级别,时间级别走势由涨跌幅,时长,波动幅度组成
// 	var barList []bars.Bar
// 	for i := 0; i < n; i++ {
// 		open := 20 + float64(i) + rand.Float64()
// 		close := open + 0.4*rand.Float64()
// 		high := math.Max(open, close) + 0.4*rand.Float64()
// 		low := math.Min(open, close) - 0.4*rand.Float64()
// 		vol := 1000000 * (1 + rand.Float64())
// 		barList = append(barList, bars.Bar{Open: open, Close: close, High: high, Low: low, Vol: vol})
// 	}
// 	return barList
// }

//	func GetArrayFromBarList(queue *Queue) ([]float64, []float64, []float64, []float64, []float64) {
//		var O, C, H, L, V []float64
//		for _, v := range queue.Values() {
//			bar := v.(bars.Bar)
//			C = append(C, bar.Close)
//			O = append(O, bar.Open)
//			H = append(H, bar.High)
//			L = append(L, bar.Low)
//			V = append(V, bar.Vol)
//		}
//		return O, C, H, L, V
//	}
func GetArrayFromMapList(queue *Queue) ([]float64, []float64, []float64, []float64, []float64) {
	var O, C, H, L, V []float64
	for _, v := range queue.Values() {
		bar := v.(map[string]float64)
		C = append(C, bar["Close"])
		O = append(O, bar["Open"])
		H = append(H, bar["High"])
		L = append(L, bar["Low"])
		V = append(V, bar["Vol"])
	}
	return O, C, H, L, V
}

// PushBarList 往Queue里面塞Bar列表
// func PushBarList(queue *Queue, barList []bars.Bar) {
// 	for _, bar := range barList {
// 		queue.Enqueue(bar)
// 	}
// }

// Ema 计算指数移动平均值,n为求EMA的参数
func Ema(array []float64, n int) float64 {
	k := 2.0 / (float64(n) + 1)
	ema := array[0]
	for i := 1; i < len(array); i++ {
		ema = array[i]*k + ema*(1-k)
	}
	return ema
}

// Sma 计算简单移动平均值,n为求SMA的参数
// 按tablib.EMA的处理方式，前N个EMA值皆为NAN,第N个EMA值为sum(c[:N]/N),其后为EMA(X,N)=[M*X+(N-M)Y’]/(N+1),N需大于M
func Sma(array []float64, n int, m float64) []float64 {
	if len(array) <= n {
		return nil
	}
	var sma []float64
	for i := 0; i < n-1; i++ {
		sma = append(sma, math.NaN())
	}
	sma = append(sma, Sum1(array[:n])/float64(n))
	for i := n; i < len(array); i++ {
		sma = append(sma, (m*array[i]+(float64(n)-m)*sma[len(sma)-1])/float64(n+1))
	}
	return sma
}

// 意图: fun(*class,)
//给定queue计算指标Eval两种方式,1.计算一个值 2.所有能算的值
