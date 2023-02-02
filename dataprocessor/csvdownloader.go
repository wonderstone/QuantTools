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
