package exporter

import (
	"io/ioutil"

	"github.com/rs/zerolog/log"
	"github.com/wonderstone/QuantTools/account/virtualaccount"
	"github.com/wonderstone/QuantTools/configer"
	"gopkg.in/yaml.v3"
)

// export realtime yaml file in the same dir as the executable file, info would be read from backtest.yaml
// configDir: the directory of the configuration files, let's say BackTest.yaml
// filename: the name of the file to be read from, BackTest it is
// sec: the section name in the BackTest.yaml, most likely "Default"
// va: the virtual account info to be added
// AInfo: additional info to be added
func ExportRealtimeYaml(configDir string, filename string, sec string, va virtualaccount.VAcct) {
	// read BackTest configuration from file
	c := configer.New(configDir + "/" + filename)
	err := c.Load()
	if err != nil {
		panic(err)
	}
	err = c.Unmarshal()
	if err != nil {
		panic(err)
	}

	m := make(map[string]interface{})
	// Add Virtual account fields
	m["va"] = va

	// Add Data fields
	tdm := make(map[string]interface{})
	// fmt.Println(viper.GetString("SMName"))
	tmpMap := c.GetStringMap(sec)
	var sinstrnames []string
	for _, v := range tmpMap["sinstrnames"].([]interface{}) {
		sinstrnames = append(sinstrnames, v.(string))
	}
	tdm["sinstrnames"] = sinstrnames
	var sindinames []string
	for _, v := range tmpMap["sindinames"].([]interface{}) {
		sindinames = append(sindinames, v.(string))
	}
	tdm["sindinames"] = sindinames
	var scsvdatafields []string
	for _, v := range tmpMap["scsvdatafields"].([]interface{}) {
		scsvdatafields = append(scsvdatafields, v.(string))
	}
	tdm["scsvdatafields"] = scsvdatafields
	var sadfields []string
	for _, v := range tmpMap["sadfields"].([]interface{}) {
		sadfields = append(sadfields, v.(string))
	}
	tdm["sadfields"] = sadfields

	var finstrnames []string
	for _, v := range tmpMap["finstrnames"].([]interface{}) {
		finstrnames = append(finstrnames, v.(string))
	}
	tdm["finstrnames"] = finstrnames
	var findinames []string
	for _, v := range tmpMap["findinames"].([]interface{}) {
		findinames = append(findinames, v.(string))
	}
	tdm["findinames"] = findinames
	var fcsvdatafields []string
	for _, v := range tmpMap["fcsvdatafields"].([]interface{}) {
		fcsvdatafields = append(fcsvdatafields, v.(string))
	}
	tdm["fcsvdatafields"] = fcsvdatafields
	var fadfields []string
	for _, v := range tmpMap["fadfields"].([]interface{}) {
		fadfields = append(fadfields, v.(string))
	}
	tdm["fadfields"] = fadfields
	m["datafields"] = tdm

	// #  Section for ContractProp
	tCPm := make(map[string]interface{})
	tCPm["confname"] = tmpMap["confname"]
	tCPm["cpdatadir"] = tmpMap["cpdatadir"]
	m["contractprop"] = tCPm
	// #  Section for Matcher parameter
	tMPm := make(map[string]interface{})
	tMPm["matcherslippage4s"] = tmpMap["matcherslippage4s"]
	tMPm["matcherslippage4f"] = tmpMap["matcherslippage4f"]
	m["matcherparam"] = tMPm
	// #  Section for Performance Analytics Parameter
	// tPAm := make(map[string]interface{})
	// tPAm["riskfreerate"] = tmpMap["riskfreerate"]
	// tPAm["patype"] = tmpMap["patype"]
	// m["pa"] = tPAm
	// #  Section for Strategy Module Selection
	tSMm := make(map[string]interface{})
	tSMm["strategymodule"] = tmpMap["strategymodule"]
	tSMm["smgeptype"] = tmpMap["smgeptype"]
	tSMm["smname"] = tmpMap["smname"]
	tSMm["smdatadir"] = tmpMap["smdatadir"]
	m["stgmodel"] = tSMm
	// export yaml file with yaml.v3
	data, err := yaml.Marshal(&m)
	if err != nil {
		log.Fatal().Msg(err.Error())

	}
	err2 := ioutil.WriteFile("./realtime.yaml", data, 0777)
	if err2 != nil {
		log.Fatal().Msg(err2.Error())
	}
	// fmt.Println("data written")
}

// for realtime job, when exit the process, update the va info by replacing the va field for realtime.yaml
// the file would in the same dir as the executable file
// configDir: the directory of the configuration files, let's say realtime.yaml
// filename: the name of the file to be updated, realtime.yaml it is
// va: the virtual account info to be updated
func ReplaceVA(configDir string, filename string, va virtualaccount.VAcct) {
	conf := configer.New(configDir + filename)
	err := conf.Load()
	if err != nil {
		panic(err)
	}
	err = conf.Unmarshal()
	if err != nil {
		panic(err)
	}

	c := conf.GetContent()
	// replace the va field
	c["va"] = va
	// export yaml file with yaml.v3
	data, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatal().Msg(err.Error())

	}
	err2 := ioutil.WriteFile("./realtime.yaml", data, 0777)
	if err2 != nil {
		log.Fatal().Msg(err2.Error())
	}
}

/*
	export the simplified Karva expression to refactor the realtime expression trees(ETs)

the file would in the same dir as the executable file
configDir: the directory of the configuration files, let's say GEP.yaml, read info from here
filename : the name of the yaml file, GEP.yaml would be GEP
secname :  the yaml file section name, e.g. "GEP"
kes : the Karva expression string slice, but put a interface{} here
*/
func ExportSKE(configDir string, filename string, secname string, KES interface{}) {
	// read GEP configuration from file.
	conf := configer.New(configDir + filename)
	err := conf.Load()
	if err != nil {
		panic(err)
	}
	err = conf.Unmarshal()
	if err != nil {
		panic(err)
	}

	tmpMap := conf.GetStringMap(secname)
	m := make(map[string]interface{})
	// Add Data fields
	// make a slice to store all the function names
	var funcnames []string
	// fmt.Println(tmpMap["funcweight"])
	for _, v := range tmpMap["funcweight"].([]interface{}) {
		tmp := v.([]interface{})
		// fmt.Println(tmp[0])
		funcnames = append(funcnames, tmp[0].(string))
	}
	m["FuncNames"] = funcnames
	// data fields
	m["HeadSize"] = tmpMap["headsize"]
	m["numConstants"] = tmpMap["numconstants"]
	m["linkFunc"] = tmpMap["linkfunc"]
	m["Mode"] = tmpMap["mode"]
	// KES field
	m["KES"] = KES
	// export yaml file with yaml.v3
	data, err := yaml.Marshal(&m)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	err2 := ioutil.WriteFile("./KarvaExp.yaml", data, 0777)
	if err2 != nil {
		log.Fatal().Msg(err2.Error())
	}
}
