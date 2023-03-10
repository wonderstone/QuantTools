package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"path/filepath"
	"strings"

	//"csv"
	"os"
	"strconv"

	"github.com/wonderstone/QuantTools/configer"
	"github.com/wonderstone/QuantTools/indicator"
	"github.com/wonderstone/QuantTools/strategyModule"
)

var fifo = "./tmp/fifo"

func next_time(hour int, minute int) (new_time string) {
	if hour >= 15 || hour <= 9 || hour == 12 {
		panic("选择时间为非交易时间")
	} else if hour == 9 && minute < 30 {
		panic("选择时间为非交易时间")
	} else if hour == 11 && minute > 30 {
		panic("选择时间为非交易时间")
	}
	minute = minute + 1
	if minute == 60 {
		hour = hour + 1
		minute = 0
	}
	if minute == 31 && hour == 11 {
		hour = 13
		minute = 0
	}
	new_hour := fmt.Sprint(hour)
	new_minute := fmt.Sprint(minute)
	if hour == 9 {
		new_hour = "09"
	}
	if len(fmt.Sprint(minute)) == 1 {
		new_minute = "0" + fmt.Sprint(minute)
	}
	new_time = new_hour + ":" + new_minute
	return new_time
}

// * Normally NewManager from Config file
func in(target string, str_array []string) bool {
	for _, element := range str_array {
		if target == element {
			return true
		}
	}
	return false
}
func indiCsvProcessDir(dirpath string, targetdir string, iis []indicator.IndiInfo, needed_indicators []string, time_strings []string, stimecritic string) bool {
	// get all the files in the dirpath
	files, err := filepath.Glob(dirpath + "/*.csv")
	to_process_csv_num := len(files)
	if err != nil {
		panic("读取文件夹出错")
	}
	// if files is empty, panic
	if len(files) == 0 {
		panic("文件夹异常 请检查 @" + dirpath)
	}
	// iter the files and process them one by one
	f, err := os.OpenFile(fifo, os.O_WRONLY|os.O_TRUNC, 0666)
	defer f.Close()
	if err != nil {
		panic("openfile error: " + err.Error())
	}

	for _, file := range files {

		// get the file name
		filename := filepath.Base(file)
		// get the target file name
		targetfile := filepath.Join(targetdir, filename)
		// process the file
		isok := indiCsvProcess(file, targetfile, iis, needed_indicators, time_strings, stimecritic)
		if !isok {
			panic("处理csv文件出错")
		}
		target_files, err := filepath.Glob(targetdir + "/*.csv")
		if err != nil {
			panic("openfile error: " + err.Error())
		}
		processed_num := len(target_files)
		processed_ratio, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(processed_num)/float64(to_process_csv_num)), 64)
		_, fseek_err := f.Seek(0, 0)
		if fseek_err != nil {
			panic("seek error: " + err.Error())
		}
		_, fwrite_err := f.WriteString(fmt.Sprintf("%.2f", processed_ratio))
		if fwrite_err != nil {
			panic("writefile error: " + fwrite_err.Error())
		}
		//writer(f, fifo, processed_num, to_process_csv_num)

	}
	err = f.Close()
	return true
}
func ContainNaN(m map[string]float64) bool {
	for _, x := range m {
		if math.IsNaN(x) {
			return true
		}
	}
	return false
}
func indiCsvProcess(fdir string, targetdir string, iis []indicator.IndiInfo, needed_indicators []string, time_string []string, stimecritic string) bool {
	csvFile, err := os.Open(fdir)
	if err != nil {
		panic("建立csv文件handler出错")
	}
	defer csvFile.Close()
	// 读取csv文件
	csvReader := csv.NewReader(csvFile)
	header, err := csvReader.Read()
	if err != nil {
		panic("第一行读取csv文件头出错")
	}
	// store the data
	rows, err := csvReader.ReadAll()
	var is_used = make(map[string]bool, len(header))
	var row_needed = make([]bool, len(header))
	for key, value := range header {
		if in(value, needed_indicators) {
			row_needed[key] = true
			is_used[value] = true
		} else {
			row_needed[key] = false
			is_used[value] = false
		}
	}
	var temp_rows [][]string
	for _, row := range rows {
		if in(strategyModule.GetTimeValue(row[0]), time_string) {
			temp_rows = append(temp_rows, row)
		}
	}
	rows = temp_rows
	if err != nil {
		panic("整块读取csv文件出错")
	}
	time_iter := 0
	var newheader []string
	for iter, key := range header {
		if key == "Time" {
			time_iter = iter
		}
		if is_used[key] {
			newheader = append(newheader, key)
		}
	}
	newrows := make([][]string, 0)
	is := make([]*indicator.IIndicator, 0)
	for _, ii := range iis {
		indi := indicator.IndiFactory(ii)
		is = append(is, &indi)
		newheader = append(newheader, ii.Name)
	}
	need_indi_num := 0
	for _, i := range row_needed {
		if i {
			need_indi_num = need_indi_num + 1
		}
	}
	// iter the rows

	for _, row := range rows {
		tmpmap := make(map[string]float64, len(header))
		for i, j := len(header)-1, len(row)-1; i > 0; i, j = i-1, j-1 {
			tmpmap[header[i]], err = strconv.ParseFloat(row[j], 64)
			if err != nil {
				panic("解析csv数据出错")
			}
			//if(row[0])

		}
		temp_row := make([]string, need_indi_num)
		for iter, indi := range row {
			if row_needed[iter] {
				temp_row[iter] = indi
			}
		}
		// iter the indicator slice and iis slice
		newrow := temp_row
		for _, indi := range is {
			if !ContainNaN(tmpmap) && strategyModule.GetTimeValue(row[time_iter]) == stimecritic {
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
	return true
}
func data_pretreatment(yamlpath string, BTfilename string, STGfilename string) {
	c := configer.New(yamlpath + BTfilename)
	c_err := c.Load()
	if c_err != nil {
		panic(c_err)
	}
	c_err = c.Unmarshal()
	if c_err != nil {
		panic(c_err)
	}
	var indicator_indiinfo []indicator.IndiInfo
	var new_provided_indicators []string
	provided_indicators := c.GetStringSlice("default.scsvdatafields")
	needed_indicators := c.GetStringSlice("default.sindinames")
	s := configer.New(yamlpath + STGfilename)
	s_err := s.Load()
	if s_err != nil {
		panic(s_err)
	}
	s_err = s.Unmarshal()
	if s_err != nil {
		panic(s_err)
	}
	var time_strings []string
	//time_strings := []string{"14:50", "14:51", "15:00"}
	dirpath := c.GetString("default.stockdatadir")
	targetdir := c.GetString("default.stockdatadirfinal")
	stimecritic := s.GetString("default.stimecritic")
	Selected_time := strings.Split(stimecritic, ":")
	hour, err := strconv.Atoi(Selected_time[0])
	if err != nil {
		panic(s_err)
	}
	minute, err := strconv.Atoi(Selected_time[1])
	if err != nil {
		panic(s_err)
	}
	new_time := next_time(hour, minute)
	time_strings = append(time_strings, stimecritic)
	time_strings = append(time_strings, new_time)
	if new_time != "15:00" {
		time_strings = append(time_strings, "15:00")
	}
	//hour := strconv.ParseInt(Selected_time[0], 10, 64)
	//minute := Selected_time[1]
	for _, indicator := range provided_indicators {
		if in(indicator, needed_indicators) {
			new_provided_indicators = append(new_provided_indicators, indicator)
		}
	}
	all_indicators := indicator.GetIndiInfoSlice(yamlpath)
	for _, temp_indicator := range needed_indicators {
		for _, single_indicator := range all_indicators {
			if temp_indicator == single_indicator.Name {
				indicator_indiinfo = append(indicator_indiinfo, single_indicator)
			}
		}
	}
	needed_indicators = append(needed_indicators, "Time")
	indiCsvProcessDir(dirpath, targetdir, indicator_indiinfo, needed_indicators, time_strings, stimecritic)
}

func main() {
	yamldir := "./config/Manual/"
	BTfilename := "BackTest.yaml"
	STGfilename := "Strategy.yaml"
	data_pretreatment(yamldir, BTfilename, STGfilename)
	//var configdirPtr = flag.String("configdir", "./config/Manual/", "a string")
	//m := NewManagerfromConfig("default", "default", *configdirPtr)
}
