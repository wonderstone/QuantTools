package dataprocessor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wonderstone/QuantTools/configer"
	"github.com/wonderstone/QuantTools/indicator"
)

// test CsvProcess
func TestCsvProcess(t *testing.T) {

	csvpath := "../tmpdata/stockdata/test/sh510050.csv"
	csvtarpath := "../tmpdata/stockdata/test/res/sh510050.csv"
	indis := []indicator.IndiInfo{
		{Name: "MA3", IndiType: "MA", ParSlice: []int{3}, InfoSlice: []string{"close"}},
		{Name: "Var3", IndiType: "Var", ParSlice: []int{3}, InfoSlice: []string{"close"}},
	}
	isok := CsvProcess(csvpath, csvtarpath, indis)
	assert.Equal(t, isok, true)

}

// test CsvProcessDir
// make sure the csv files have the same header or panic
func TestCsvProcessDir(t *testing.T) {

	csvdir := "../tmpdata/stockdata/1min/"
	csvtardir := "../tmpdata/stockdata/1min/res/"
	indis := []indicator.IndiInfo{
		{Name: "MA3", IndiType: "MA", ParSlice: []int{3}, InfoSlice: []string{"close"}},
		{Name: "Var3", IndiType: "Var", ParSlice: []int{3}, InfoSlice: []string{"close"}},
	}
	isok := CsvProcessDir(csvdir, csvtardir, indis)
	assert.Equal(t, isok, true)

}

// func to get info from Backtest.yaml and run CsvProcessDir
func CsvProcessDirfromConfig(btpath string, iidir string) {
	// read the config file backtest.yaml and get the
	c := configer.New(btpath)
	err := c.Load()
	if err != nil {
		panic(err)
	}
	err = c.Unmarshal()
	if err != nil {
		panic(err)
	}

	adf := c.GetStringSlice("default.sadfields")
	csvdir := c.GetString("default.stockdatadir")
	csvdirfinal := c.GetString("default.stockdatadirfinal")
	fmt.Println(adf)

	// read the indicatorinfo.yaml
	iis := indicator.GetIndiInfoSlice(iidir)
	fmt.Println(iis)

	// make an ii slice for CsvProcessDir whose ii.Name is in adf
	var iis2 []indicator.IndiInfo
	for _, v := range adf {
		isok := false
		for _, ii := range iis {
			if ii.Name == v {
				iis2 = append(iis2, ii)
				isok = true
				break
			}
		}
		if !isok {
			panic("indicator name not found! @ " + v)
		}
	}

	// run CsvProcessDir

	isok := CsvProcessDir(csvdir, csvdirfinal, iis2)
	if isok {
		fmt.Println("CsvProcessDirfromConfig done!")
	} else {
		fmt.Println("CsvProcessDirfromConfig failed!")
	}
}

// test CsvProcessDirfromConfig

func TestCsvProcessDirfromConfig(t *testing.T) {
	BTpath := "../config/Manual/BackTest.yaml"
	IndiDir := "../config/Manual/"
	CsvProcessDirfromConfig(BTpath, IndiDir)

}
