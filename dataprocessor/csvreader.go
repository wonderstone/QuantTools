package dataprocessor

import (
	"encoding/csv"
	// "io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 假定一个标的的Bar信息和indicator信息都已经整理完存在一个csv文件中
// viper is not thread safe, add mutex to protect the data read process and do it before iteration or goroutine
func (BCM *BarCM) CsvSBarReader(filedir string) {
	if len(BCM.InstSIDS) == 0 {
		panic("股票标的切片为空 检查策略逻辑")
	}
	csvFile, err := os.Open(filedir)
	if err != nil {
		panic("建立csv文件handler出错")
	}
	defer csvFile.Close()
	// get the instID from the file name
	instSID := strings.TrimSuffix(filepath.Base(filedir), filepath.Ext(filepath.Base(filedir)))
	// 逐行读取csv文件
	csvReader := csv.NewReader(csvFile)
	header, err := csvReader.Read()
	if err != nil {
		panic("第一行读取csv文件头出错")
	}
	// if the indiSNames is nil ,then init it
	// else compare indiSNames and header
	if BCM.IndiSNames == nil {
		BCM.IndiSNames = header[1:]
	} else {
		if !CompareStringSlices(BCM.IndiSNames, header[1:]) {
			panic("股票标的的csv文件头与策略配置的不一致@" + instSID)
		}
	}
	// store the data
	rows, err := csvReader.ReadAll()
	if err != nil {
		panic("整块读取csv文件出错")
	}
	for _, row := range rows {
		// 需要根据最终csv字段进行调整
		dtstr := row[0]
		// iterate the header backwards and get the data in a temp map
		tmpmap := make(map[string]float64, len(header))
		for i, j := len(header)-1, len(row)-1; i > 0; i, j = i-1, j-1 {
			tmpmap[header[i]], err = strconv.ParseFloat(row[j], 64)
			if err != nil {
				panic("解析csv数据出错")
			}
		}
		//更新map的stockdata
		pSBDE := NewBarDE(dtstr, instSID, tmpmap)
		if BCM.BarCMap[dtstr] == nil {
			pBC := NewBarC(len(BCM.InstSIDS))
			pBC.Stockdata[instSID] = pSBDE
			BCM.BarCMap[dtstr] = pBC
		} else {
			BCM.BarCMap[dtstr].Stockdata[instSID] = pSBDE
		}
	}

}

func (BCM *BarCM) CsvFBarReader(filedir string) {
	if len(BCM.InstFIDS) == 0 {
		panic("期货标的切片为空 检查策略逻辑")
	}
	csvFile, err := os.Open(filedir)
	if err != nil {
		panic("建立csv文件handler出错")
	}
	defer csvFile.Close()
	// get the instID from the file name
	instFID := strings.TrimSuffix(filepath.Base(filedir), filepath.Ext(filepath.Base(filedir)))
	// 逐行读取csv文件
	csvReader := csv.NewReader(csvFile)
	header, err := csvReader.Read()
	if err != nil {
		panic("第一行读取csv文件头出错")
	}
	// if the indiSNames is nil ,then init it
	// else compare indiSNames and header
	if BCM.IndiFNames == nil {
		BCM.IndiFNames = header[1:]
	} else {
		if !CompareStringSlices(BCM.IndiFNames, header[1:]) {
			panic("股票标的的csv文件头与策略配置的不一致@" + instFID)
		}
	}
	// store the data
	rows, err := csvReader.ReadAll()
	if err != nil {
		panic("整块读取csv文件出错")
	}
	for _, row := range rows {
		// 需要根据最终csv字段进行调整
		dtstr := row[0]
		// iterate the csvdatafields backwards and get the data in a temp map
		tmpmap := make(map[string]float64, len(header))
		for i, j := len(header)-1, len(row)-1; i > 0; i, j = i-1, j-1 {
			tmpmap[header[i]], err = strconv.ParseFloat(row[j], 64)
			if err != nil {
				panic("解析csv数据出错")
			}
		}
		//更新map的futuresdata
		pFBDE := NewBarDE(dtstr, instFID, tmpmap)
		if BCM.BarCMap[dtstr] == nil {
			pBC := NewBarC(len(BCM.InstFIDS))
			pBC.Futuresdata[instFID] = pFBDE
			BCM.BarCMap[dtstr] = pBC
		} else {
			BCM.BarCMap[dtstr].Futuresdata[instFID] = pFBDE
		}
	}

}

func (BCM *BarCM) CsvFMTMReader(filedir string) {
	if len(BCM.InstFIDS) == 0 {
		panic("期货标的切片为空 检查策略逻辑")
	}
	csvFile, err := os.Open(filedir)
	if err != nil {
		panic("建立csv文件handler出错")
	}
	// get the instID from the file name
	instFID := strings.TrimSuffix(filepath.Base(filedir), filepath.Ext(filepath.Base(filedir)))
	defer csvFile.Close()
	// 逐行读取csv文件
	csvReader := csv.NewReader(csvFile)
	header, err := csvReader.Read()
	if err != nil {
		panic("第一行读取csv文件头出错")
	}
	var idt, ise int
	for i, v := range header {
		if v == "datetime" {
			idt = i
		}
		if v == "settlementprice" {
			ise = i
		}
	}

	rows, err := csvReader.ReadAll()
	if err != nil {
		panic("整块读取csv文件出错")
	}

	for _, row := range rows {
		BCM.FMTMDataMap[row[idt]] = make(map[string]float64)
		BCM.FMTMDataMap[row[idt]][instFID], err = strconv.ParseFloat(row[ise], 64)
		if err != nil {
			panic("解析csv数据出错")
		}
	}
}

func CompareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true

}

// func ListDir(dirPth string, suffix string) (files []string, err error) {
// 	files = make([]string, 0, 30)
// 	dir, err := ioutil.ReadDir(dirPth)
// 	if err != nil {
// 		return nil, err
// 	}
// 	PthSep := string(os.PathSeparator)
// 	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
// 	for _, fi := range dir {
// 		if fi.IsDir() { // 忽略目录
// 			continue
// 		}
// 		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
// 			files = append(files, dirPth+PthSep+fi.Name())
// 		}
// 	}
// 	return files, nil
// }

// use os.ReadDir(dirname)
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	dir, err := os.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, nil
}
