package contractproperty

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test NewStockContractPropFromConfig
func TestNewStockContractPropFromConfig(t *testing.T) {
	fmt.Println("test NewStockContractPropFromConfig")
	confName := "ContractProp"
	dir := "../config/Manual"
	sec := "st"
	scp := NewSCPFromConfig(confName, sec, dir)
	fmt.Println(scp)
	assert.Equal(t, float64(100), scp.ContractSize)
}

// test NewFuturesContractPropFromConfig
func TestNewFuturesContractPropFromConfig(t *testing.T) {
	fmt.Println("test NewFuturesContractPropFromConfig")
	confName := "ContractProp"
	dir := "../config/Manual"
	instrID := "au2210"
	fcp := NewFCPFromConfig(confName, instrID, dir)
	fmt.Println(fcp)
	assert.Equal(t, float64(1000), fcp.ContractSize)
}

// test NewContractPropMap
func TestNewContractPropMap(t *testing.T) {
	fmt.Println("test NewContractPropMap")
	confName := "ContractProp"
	dir := "../config/Manual"
	cpm := NewCPMap(dir, confName)

	assert.Equal(t, float64(100), cpm.StockPropMap["st"].ContractSize)
	fmt.Println(cpm)
}

// test regexp
func TestRegexp(t *testing.T) {
	str1 := "sh600002"
	str2 := "au2210"
	str3 := "d[12]"
	re4char := regexp.MustCompile("[a-zA-Z]*")
	re4num := regexp.MustCompile("[0-9]*$")
	re4d := regexp.MustCompile("d[[0-9]+]")
	fmt.Println(re4char.FindString(str1))
	assert.Equal(t, "sh", re4char.FindString(str1))
	fmt.Println(re4num.FindString(str2))
	assert.Equal(t, "2210", re4num.FindString(str2))
	fmt.Println(re4d.FindString(str3))
	assert.Equal(t, "d[12]", re4d.FindString(str3))
}

// test SimpleNewSCPFromMap
func TestSimpleNewSCPFromMap(t *testing.T) {
	fmt.Println("test SimpleNewSCPFromMap")
	confName := "ContractProp"
	dir := "../config/Manual"
	cpm := NewCPMap(dir, confName)
	// 原则上需要开头字母标注市场。接受sh600002和sh.600002的格式输入
	// 当输入600002时，只是返回默认的主板类型。
	// 当不核查涨跌幅时，此处区分意义有限，保留只是为了未来扩展。
	scp := SimpleNewSCPFromMap(cpm, "sh.600002")
	fmt.Println(scp)
	assert.Equal(t, float64(100), scp.ContractSize)
}

// test SimpleNewFSCPFromMap
func TestSimpleNewFSCPFromMap(t *testing.T) {
	fmt.Println("test SimpleNewFSCPFromMap")
	confName := "ContractProp"
	dir := "../config/Manual"
	cpm := NewCPMap(dir, confName)
	fcp := SimpleNewFCPFromMap(cpm, "au2210")
	fmt.Println(fcp)
	assert.Equal(t, float64(1000), fcp.ContractSize)

}
