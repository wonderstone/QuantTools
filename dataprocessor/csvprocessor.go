package dataprocessor

import (
	"encoding/csv"
	"math"
	"os"
	"path/filepath"

	"github.com/wonderstone/QuantTools/indicator"

	"strconv"
	//"strings"
)

// function with dir path to get all files and CsvProcess them one by one and output to the targetdir
func CsvProcessDir(dirpath string, targetdir string, iis []indicator.IndiInfo) (ok bool, err error) {
	// get all the files in the dirpath
	files, err := filepath.Glob(dirpath + "/*.csv")
	if err != nil {
		panic("读取文件夹出错")
		// return false, fmt.Errorf("读取文件夹出错")
	}
	// iter the files and process them one by one
	for _, file := range files {
		// get the file name
		filename := filepath.Base(file)
		// get the target file name
		targetfile := filepath.Join(targetdir, filename)
		// process the file
		_, err := CsvProcess(file, targetfile, iis)
		if err != nil {
			panic("处理csv文件出错")
			// return false, fmt.Errorf("处理csv文件出错")
		}
	}
	return true, nil
}

// function to read csv file add some datas and write to a new csv file
func CsvProcess(fpath string, targetdir string, iis []indicator.IndiInfo) (ok bool, err error) {

	// open the csv file
	csvFile, err := os.Open(fpath)
	if err != nil {
		panic("建立csv文件handler出错")
		// return false, fmt.Errorf("建立csv文件handler出错")
	}
	defer csvFile.Close()
	// get the instID from the file name
	//instSID := strings.TrimSuffix(filepath.Base(filedir), filepath.Ext(filepath.Base(filedir)))
	// 逐行读取csv文件
	csvReader := csv.NewReader(csvFile)
	header, err := csvReader.Read()
	if err != nil {
		panic("第一行读取csv文件头出错")
		// return false, fmt.Errorf("第一行读取csv文件头出错")
	}
	// store the data
	rows, err := csvReader.ReadAll()
	if err != nil {
		panic("整块读取csv文件出错")
	}

	// new header and new rows
	newheader := header
	newrows := make([][]string, 0)
	// get the indicator slice for pointer to the indicator
	is := make([]*indicator.IIndicator, 0)
	for _, ii := range iis {
		indi := indicator.IndiFactory(ii)
		is = append(is, &indi)
		newheader = append(newheader, ii.Name)
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
				panic("解析csv数据出错")
				// return false, fmt.Errorf("解析csv数据出错")
			}
		}
		// iter the indicator slice and iis slice
		newrow := row
		for _, indi := range is {
			if !ContainNaN(tmpmap) {
				// load the data into the indicator
				(*indi).LoadData(tmpmap)
				// update the new row for the new csv file
				newrow = append(newrow, strconv.FormatFloat((*indi).Eval(), 'f', 6, 64))
			} else {
				// update the new row with NaN
				newrow = append(newrow, "NaN")
			}
		}
		newrows = append(newrows, newrow)

	}

	// get the dir from the targetdir
	dir := filepath.Dir(targetdir)
	// create the dir if not exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	// write the data out for another csv file
	// create a new csv file with the targetdir
	newcsvFile, err := os.Create(targetdir)
	if err != nil {
		panic("创建新csv文件出错")
		// return false, fmt.Errorf("创建新csv文件出错")
	}
	defer newcsvFile.Close()
	// create a new csv writer
	newcsvWriter := csv.NewWriter(newcsvFile)
	// write the header
	newcsvWriter.Write(newheader)
	// write the data
	for _, row := range newrows {
		newcsvWriter.Write(row)
	}
	// flush the data
	newcsvWriter.Flush()
	return true, nil

}

func ContainNaN(m map[string]float64) bool {
	for _, x := range m {
		if math.IsNaN(x) {
			return true
		}
	}
	return false
}
