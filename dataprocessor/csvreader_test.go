package dataprocessor

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCsvReader(t *testing.T) {
	instSIDS := []string{"sz000058", "sh600409"}
	indiSNames := []string{"open", "close"}

	instFIDS := []string{"a1409"}
	indiFNames := []string{"open", "close"}

	beginDate := "2017/10/9 9:39"
	endDate := "2017/10/20 15:00"
	dirSPth := "../tmpdata/stockdata/1min"
	dirFPth := "../tmpdata/futuresdata/1min"
	dirFMTMPth := "../tmpdata/futuresdata/1day"

	files, err := ListDir(dirSPth, "csv")
	if err != nil {
		panic(err)
	}

	fmt.Println(files)

	// files1, err1 := ListDirNew(dirSPth, "csv")
	// if err1 != nil {
	// 	panic(err1)
	// }

	// fmt.Println(files1)

	bcm := NewBarCM(instSIDS, indiSNames, instFIDS, indiFNames, beginDate, endDate)
	for _, file := range files {
		instSID := strings.TrimSuffix(filepath.Base(file), filepath.Ext(filepath.Base(file)))
		fmt.Println(instSID)
		bcm.CsvSBarReader(file)
	}

	fmt.Println("Stock data finishe!")
	fmt.Println(bcm.BarCMap["2017/10/9 9:55"].Stockdata["sz000058"].IndiDataMap["ma10"])

	assert.Equal(t, math.IsNaN(bcm.BarCMap["2017/10/9 9:55"].Stockdata["sz000058"].IndiDataMap["ma10"]), true)

	FPfiles, err := ListDir(dirFPth, "csv")
	if err != nil {
		panic(err)
	}
	fmt.Println(FPfiles)
	for _, FPfile := range FPfiles {
		instSID := strings.TrimSuffix(filepath.Base(FPfile), filepath.Ext(filepath.Base(FPfile)))
		fmt.Println(instSID)
		bcm.CsvFBarReader(FPfile)
	}

	assert.Equal(t, bcm.BarCMap["2015/1/5 9:16"].Futuresdata["a1409"].IndiDataMap["open"], 3619.0)

	MTMfiles, err := ListDir(dirFMTMPth, "csv")
	if err != nil {
		panic(err)
	}

	fmt.Println(MTMfiles)
	for _, MTMfile := range MTMfiles {
		instSID := strings.TrimSuffix(filepath.Base(MTMfile), filepath.Ext(filepath.Base(MTMfile)))
		fmt.Println(instSID)
		bcm.CsvFMTMReader(MTMfile)
	}
	assert.Equal(t, bcm.FMTMDataMap["2013/3/15"]["a1409"], 4824.0)

}
