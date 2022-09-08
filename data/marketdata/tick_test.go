/* 未完结  需要等待股票tick数据的具体情况*/
/* 使用基本的testing进行的测试*/
package marketdata

import (
	"testing"

	"github.com/wonderstone/QuantTools/data"
)

type TestCase struct {
	value    FuturesTick
	expected data.UpdateMI
	Actual   data.UpdateMI
}

// test for GetUpdateInfo()
func TestGetUpdateInfo(t *testing.T) {
	testCase := TestCase{
		value: FuturesTick{
			UpdateTimeStamp: "2022-05-10 12:12:12 500",
			InstID:          "cu",
			LastPrice:       3459.2,
		},
		expected: data.UpdateMI{
			UpdateTimeStamp: "2022-05-10 12:12:12 500",
			InstID:          "cu",
			Value:           3459.2,
		},
	}
	testCase.Actual = testCase.value.GetUpdateInfo()
	if testCase.Actual != testCase.expected {
		t.Fatal("Expected Result Not Given")
	}
}
