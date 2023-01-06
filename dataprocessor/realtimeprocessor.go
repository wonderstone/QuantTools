package dataprocessor

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/wonderstone/QuantTools/indicator"
)

// realtime processor has 3 parts:
// 1. get the history data from source
func FakeGetHistoryData(dir string) (BarCMap map[string]*BarC, BarCMapkeydts []string) {
	// init the BarCMap
	BarCMap = make(map[string]*BarC)
	// 1. get the history data from source, this fake data is from csvreader.go
	// 1.1 get the data from csv files in dir
	files, err := ListDir(dir, "csv")
	if err != nil {
		panic(err)
	}
	for _, file := range files {

		csvFile, err := os.Open(file)
		if err != nil {
			panic("建立csv文件handler出错")
		}
		defer csvFile.Close()
		// get the instID from the file name
		instSID := strings.TrimSuffix(filepath.Base(file), filepath.Ext(filepath.Base(file)))
		fmt.Println(instSID)
		// 逐行读取csv文件
		csvReader := csv.NewReader(csvFile)
		header, err := csvReader.Read()
		if err != nil {
			panic("第一行读取csv文件头出错@" + instSID)
		}
		fmt.Println(header)
		// store the data
		rows, err := csvReader.ReadAll()
		if err != nil {
			panic("整块读取csv文件出错@" + instSID)
		}
		for _, row := range rows {
			// 需要根据最终csv字段进行调整
			dtstr := row[0]
			// iterate the header backwards and get the data in a temp map
			tmpmap := make(map[string]float64, len(header))
			for i, j := len(header)-1, len(row)-1; i > 0; i, j = i-1, j-1 {
				tmpmap[header[i]], err = strconv.ParseFloat(row[j], 64)
				if err != nil {
					panic("解析csv数据出错@" + instSID)
				}
			}
			//更新map的stockdata
			pSBDE := NewBarDE(dtstr, instSID, tmpmap)
			if BarCMap[dtstr] == nil {
				pBC := NewBarC(len(files))
				pBC.Stockdata[instSID] = pSBDE
				BarCMap[dtstr] = pBC
			} else {
				BarCMap[dtstr].Stockdata[instSID] = pSBDE
			}
		}

	}

	// 1.2 get the BarCMapkeydts
	for k := range BarCMap {
		BarCMapkeydts = append(BarCMapkeydts, k)
	}
	// sort the BarCMapkeydts
	sort.Slice(BarCMapkeydts, func(i, j int) bool {
		dti, _ := time.Parse("2006/1/2 15:04", BarCMapkeydts[i])
		dtj, _ := time.Parse("2006/1/2 15:04", BarCMapkeydts[j])
		return dti.Before(dtj)
	})

	return BarCMap, BarCMapkeydts
}

// 2. add indicators to one data
func AddIndicatorsToSData(pBC *BarC, ID string, pIndicators []indicator.IIndicator) {
	// 2.1 add indicators to one data
	for _, pIndicator := range pIndicators {
		// check if nan, only load data when no nan
		if !ContainNaN(pBC.Stockdata[ID].IndiDataMap) {
			pIndicator.LoadData(pBC.Stockdata[ID].IndiDataMap)
			pBC.Stockdata[ID].IndiDataMap[pIndicator.GetName()] = pIndicator.Eval()
			fmt.Println(pIndicator)
		} //else {
		//fmt.Println("nan")
		//}

	}
}

// 2. get the realtime data from source and process it
