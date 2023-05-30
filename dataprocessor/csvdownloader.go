// this file is for downloading csv file from the VDS
// 1. read the info needed from the config file BackTest.yaml
//    data range(begindate enddate)、targets(sinstrnames)、data on VDS(scsvdatafields)

// 2. iterate the targets
//    download the data and save csv files to target dir(stockdatadir)
//    with VDS data format see in tmpdata/stockdata/test/sh510050.csv
// *  Time,Open,Close,High,Low,Volume,Amount
// *  2023.01.18T09:31:00.000,51.5,51.38,51.5,51.36,239600,12333986.61

// author:  CheYang (Digital Office Product Department #2)
package dataprocessor

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"io/ioutil"
	"net/http"
	"strings"
)

type HTTPRspBody struct {
	Result Results `json:"Result"`
}
type Results struct {
	RequestID     string   `json:"Result"`
	HasError      bool     `json:"HasError"`
	ResponseItems ErrorMsg `json:"ResponseItems"`
}
type ErrorMsg struct {
	ErrorMsg string `json:"ErrorMsg"`
}
type Login struct {
	Uname string `json:"uname" validate:"required"`
	Upwd  string `json:"upwd" validate:"required"`
}
type historyData_request struct {
	Symbol     string `json:"symbol" validate:"required"`
	StartDt    string `json:"start_dt" validate:"required"`
	EndDt      string `json:"end_dt" validate:"required"`
	Count      int    `json:"count" validate:"required"`
	Field      string `json:"fields" validate:"required"`
	CandleType string `json:"candle_type" validate:"required"`
}
type Response_body struct {
	Code     int    `json:"code"`
	Text     string `json:"text"`
	Data     string `json:"data"`
	Userdata string `json:"userdata"`
}

func acquire_token(lg *Login, url string) string {
	user_login_json, err := json.Marshal(lg)
	if err != nil {
		panic(err)
	}
	loginReqBody := strings.NewReader(string(user_login_json))
	loginHttpReq, err := http.NewRequest("POST", url, loginReqBody)
	if err != nil {
		panic(err)
	}
	loginHttpReq.Header.Add("Content-Type", "application/json")
	loginHttpRsp, err := http.DefaultClient.Do(loginHttpReq)
	if err != nil {
		panic(err)
	}
	defer loginHttpRsp.Body.Close()
	loginRspBody, err := ioutil.ReadAll(loginHttpRsp.Body)
	if err != nil {
		panic(err)
	}
	var result Response_body
	if err = json.Unmarshal(loginRspBody, &result); err != nil {
		panic(err)
	}
	return result.Data

}
func myReverse(l [][]interface{}) [][]interface{} {
	new_l := make([][]interface{}, len(l))
	for i := len(l) - 1; i >= 0; i-- {
		new_l[len(l)-i-1] = l[i]
	}
	return new_l
}

type VDSData struct {
	Fields []string        `json:"fields"`
	Items  [][]interface{} `json:"items"`
}

func historydata_download(hr *historyData_request, uname string, upwd string) [][]interface{} {
	url := "http://10.1.90.91:9002/xbzq/vds/v1/hq/candle"
	jsonData, err := json.Marshal(hr)
	if err != nil {
		panic(err)
	}
	var user_login Login
	user_login.Uname = uname
	user_login.Upwd = upwd
	token := acquire_token(&user_login, url)
	reqBody := strings.NewReader(string(jsonData))
	httpReq, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		panic(err)
	}
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", token)
	httpRsp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		panic(err)
	}
	defer httpRsp.Body.Close()
	RspBody, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		panic(err)
	}
	var result Response_body
	if err = json.Unmarshal(RspBody, &result); err != nil {
		panic(err)
	}
	var vdsdata VDSData
	if err = json.Unmarshal([]byte(result.Data), &vdsdata); err != nil {
		panic(err)
	}
	reversedata := myReverse(vdsdata.Items)
	for _, data := range reversedata {
		fmt.Println(data)
	}
	return reversedata
}
func exch_recog(vdsdata [][]interface{}) string {
	switch vdsdata[0][0] {
	case 1:
		return "sh"
	case 2:
		return "sz"
	case 3:
		return "bj"
	default:
		return ""
	}

}
func csv_download(vdsdata [][]interface{}, fdir string, fname string) {
	csvFile, err := os.Open(fdir + fname)
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
	fmt.Println(header)
	var temp_rows [][]interface{}

	for i := 0; i < len(vdsdata); i++ {
		temp_rows = append(temp_rows, vdsdata[i][2:])
	}
	exch := exch_recog(temp_rows)
	symbol := vdsdata[0][1].(string)
	targetdir := exch + symbol + ".csv"
	dir := filepath.Dir(targetdir)
	// create the dir if not exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	newcsvFile, err := os.Create(targetdir)
	if err != nil {
		panic("创建新csv文件出错")
	}
	defer newcsvFile.Close()
	// create a new csv writer
	newcsvWriter := csv.NewWriter(newcsvFile)
	// write the header
	newcsvWriter.Write(header)
	string_rows := make([][]string, len(temp_rows))
	for arr := range string_rows {
		string_rows[arr] = make([]string, len(temp_rows[0]))
	}
	for i, arg := range temp_rows {
		for j, arg_string := range arg {
			string_rows[i][j] = fmt.Sprintf("%v", arg_string)
		}
	}
	for _, row := range string_rows {
		newcsvWriter.Write(row)
	}
	newcsvWriter.Flush()
}
