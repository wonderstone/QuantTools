package dataprocessor

import (
	"encoding/csv"
	"fmt"
	"os"

	//"path/filepath"
	"strconv"
	//"strings"
)

// function to read csv file add some datas and write to a new csv file
func CsvProcess(filedir string) (ok bool, err error) {
	// open the csv file
	csvFile, err := os.Open(filedir)
	if err != nil {
		return false, fmt.Errorf("建立csv文件handler出错")
	}
	defer csvFile.Close()
	// get the instID from the file name
	//instSID := strings.TrimSuffix(filepath.Base(filedir), filepath.Ext(filepath.Base(filedir)))
	// 逐行读取csv文件
	csvReader := csv.NewReader(csvFile)
	header, err := csvReader.Read()
	if err != nil {
		return false, fmt.Errorf("第一行读取csv文件头出错")
	}
	// store the data
	rows, err := csvReader.ReadAll()
	if err != nil {
		panic("整块读取csv文件出错")
	}

	// iter the rows
	for _, row := range rows {
		// 需要根据最终csv字段进行调整
		//dtstr := row[0]
		// iterate the header backwards and get the data in a temp map
		tmpmap := make(map[string]float64, len(header))
		for i, j := len(header)-1, len(row)-1; i > 0; i, j = i-1, j-1 {
			tmpmap[header[i]], err = strconv.ParseFloat(row[j], 64)
			if err != nil {
				return false, fmt.Errorf("解析csv数据出错")
			}
		}

	}

	// write the data out for another csv file
	// create a new csv file
	newcsvFile, err := os.Create("newcsv.csv")
	if err != nil {
		return false, fmt.Errorf("创建新csv文件出错")
	}
	defer newcsvFile.Close()
	// create a new csv writer
	newcsvWriter := csv.NewWriter(newcsvFile)
	// write the header
	newcsvWriter.Write(header)
	// write the data
	for _, row := range rows {
		newcsvWriter.Write(row)
	}
	// flush the data
	newcsvWriter.Flush()
	return true, nil

}
