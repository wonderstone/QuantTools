package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"strconv"

	// "io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/viper"
	"github.com/wonderstone/QuantTools/account/virtualaccount"
	"github.com/wonderstone/QuantTools/framework"

	"github.com/rs/zerolog/log"

	// "sync"

	"github.com/wonderstone/GEP-MOD/functions"
	"github.com/wonderstone/GEP-MOD/gene"
	"github.com/wonderstone/GEP-MOD/genome"
	"github.com/wonderstone/GEP-MOD/genomeset"
	"github.com/wonderstone/GEP-MOD/grammars"
	"github.com/wonderstone/GEP-MOD/model"
)

const debug = false

// init the random number generator
func init() {
	rand.Seed(time.Now().UnixNano())
}

// get GEP configuration from file
func getFuncWeight(sec string, dir string) (res []gene.FuncWeight, resMap map[string]interface{}, ExpectFitness float64, linkFunc string) {
	// read GEP configuration from file  viper is not thread safe
	viper.SetConfigName(sec)
	viper.AddConfigPath(dir)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	tmpMap := viper.GetStringMap("GEP")
	// fmt.Println("tmpMap: ", tmpMap)

	res = make([]gene.FuncWeight, len(tmpMap["funcweight"].([]interface{})))

	for k, v := range tmpMap["funcweight"].([]interface{}) {
		res[k].Symbol = v.([]interface{})[0].(string)
		res[k].Weight = v.([]interface{})[1].(int)
	}
	return res, tmpMap, viper.GetFloat64("GEP.expectFitness"), viper.GetString("GEP.linkFunc")

}

// declare the manager struct used for aggregating the backtest and strategy module
// backtest has the parameters and the market data
// strategy interface relates to the strategy module
type manager struct {
	BT     *framework.BackTest // BackTest framework component
	secBT  string
	secSTG string
	dir    string
}

// NewManager creates a new manager instance
// * Normally NewManager from Config file
func NewManagerfromConfig(secBT string, secSTG string, dir string) *manager {
	BT := framework.NewBackTestConfig(dir, "BackTest.yaml", secBT)
	return &manager{
		BT:     &BT,
		secBT:  secBT,
		secSTG: secSTG,
		dir:    dir,
	}
}

// validateFunc is used to provide the fitness function for the GEP
// it gets needed data and strategy action from the manager struct	and return the fitness
func (m *manager) validateFunc(g *genome.Genome) (result float64) {
	// var mutex sync.Mutex
	// new a strategy from backtest
	pstg := m.BT.GetStrategy(m.dir, "BackTest.yaml", m.secSTG, "Strategy.yaml", "DMT")
	// new virtual account
	va := virtualaccount.NewVirtualAccount(m.BT.BeginDate, m.BT.StockInitValue, m.BT.FuturesInitValue)
	log.Info().Str("Account UUID", va.SAcct.UUID).Float64("AccountVal", va.SAcct.MktVal).Msg("Virtual Account Created!")

	// IterData 这个写法感觉不太好  有机会调整一下
	m.BT.IterData(&va, m.BT.BCM, pstg, m.BT.CPMap, func(gepin []float64) []float64 { return []float64{g.EvalMath(gepin)} }, "GEP")
	einfo := make(map[string]interface{})
	einfo["tag"] = m.BT.PAType
	einfo["par"] = 0.0
	tmp := m.BT.EvalPerformance(va.SAcct.MarketValueSlice, einfo)

	log.Info().Str("Account UUID", va.SAcct.UUID).Float64("AccountVal", va.SAcct.MktVal).Msg("VA EvalPerformance Finished!")
	return tmp
}

// output the VirtualAccount value to csv file
func (m *manager) outputVAvalues(g *genome.Genome, gr *grammars.Grammar) {

	// new a strategy from backtest
	pstg := m.BT.GetStrategy(m.dir, "BackTest.yaml", m.secSTG, "Strategy.yaml", "DMT")
	// new virtual account
	va := virtualaccount.NewVirtualAccount(m.BT.BeginDate, m.BT.StockInitValue, m.BT.FuturesInitValue)
	// IterData 这个写法感觉不太好  有机会调整一下
	m.BT.IterData(&va, m.BT.BCM, pstg, m.BT.CPMap, func(gepin []float64) []float64 { return []float64{g.EvalMath(gepin)} }, "GEP")
	// output the virtual account value to csv file
	file, err := os.Create("./records.csv")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()
	// Using Write
	for _, record := range va.SAcct.MarketValueSlice {
		row := []string{record.Time, strconv.FormatFloat(record.MktVal, 'f', 2, 64)}
		if err := w.Write(row); err != nil {
			log.Fatal().Msg(err.Error())
		}
	}
	// fmt.Println("--------------------------------")
	// fmt.Println(gr)
	// fmt.Printf("simple solution could be : %v\n", g)
	g.WriteExps(os.Stdout, gr, m.BT.SIndiNames)
	// yield the additional info from g
	res, err := g.SimplifyGenome(gr, m.BT.SIndiNames)
	if err != nil {
		panic("err")
	}
	// fmt.Println("--------------------------------")
	// fmt.Println(gr)
	// fmt.Printf("simple solution could be : %v\n", g)
	g.WriteExps(os.Stdout, gr, res)
	// fmt.Println(res)
	va.SAcct.ResetMVSlice()
	tmp := make(map[string]interface{})
	tmp["SIndiNmsAfter"] = res
	// exporter.ExportRealtimeYaml("./config/Manual", "Default", va, tmp)
	// make a string slice
	// The zero value of a slice type is nil
	var KES []string //zero value of a slice is nil
	for _, v := range g.Genes {
		KES = append(KES, v.String())
	}

	// exporter.ExportSKE("./config/Manual", "GEP", KES)
	// use exporter to output the RTyaml file
	// exporter.ExportRealtimeYaml("../config/Manual", "Default", va, tmp)
}

func (m *manager) validateFuncGS(g *genomeset.GenomeSet) (result float64) {
	// var mutex sync.Mutex
	// new a strategy from backtest
	pstg := m.BT.GetStrategy(m.dir, "BackTest.yaml", m.secSTG, "Strategy.yaml", "DMT")
	// new virtual account
	va := virtualaccount.NewVirtualAccount(m.BT.BeginDate, m.BT.StockInitValue, m.BT.FuturesInitValue)
	// IterData
	// mutex.Lock()
	// 这个写法感觉不太好  有机会调整一下
	m.BT.IterData(&va, m.BT.BCM, pstg, m.BT.CPMap,
		func(gepin []float64) []float64 {
			res := []float64{}
			for _, v := range g.Genomes {
				res = append(res, v.EvalMath(gepin))
			}
			return res
		}, "GEP")
	// mutex.Unlock()
	einfo := make(map[string]interface{})
	einfo["tag"] = m.BT.PAType
	tmp := m.BT.EvalPerformance(va.SAcct.MarketValueSlice, einfo)

	return tmp
}

func main() {
	var configdirPtr = flag.String("configdir", "./config/Manual/", "a string")
	// create a manager instance:
	m := NewManagerfromConfig("default", "default", *configdirPtr)
	// manager prepares the market data
	m.BT.PrepareData("VDS")
	if debug {
		log.Info().Msg("Data Prepared!")
	}
	// read GEP configuration from file
	funcs, GEPmap, expectFitness, lincfunc := getFuncWeight("GEP", "./config/Manual/")
	// numTerminals is inferred from the data provided by the backtest in manager
	var numTerminals int
	if len(m.BT.FInstrNames) != 0 {
		numTerminals = len(m.BT.FInstrNames)
		log.Info().Msg("GEP USE Futures Variables!")
	} else if len(m.BT.SInstrNames) != 0 {
		numTerminals = len(m.BT.SIndiNames)
		log.Info().Msg("GEP USE Stock Variables!")
	}
	NumGenomes := GEPmap["numgenomes"].(int)
	NumGenomeSet := GEPmap["numgenomeset"].(int)
	HeadSize := GEPmap["headsize"].(int)
	numGenesPerGenome := GEPmap["numgenespergenome"].(int)
	numGenomesPerGenomeSet := GEPmap["numgenomespergenomeset"].(int)
	numConstants := GEPmap["numconstants"].(int)
	Iteration := GEPmap["iteration"].(int)
	pm := GEPmap["pmutate"].(float64)
	pis := GEPmap["pis"].(float64)
	glis := GEPmap["glis"].(int)
	pris := GEPmap["pris"].(float64)
	glris := GEPmap["glris"].(int)
	pgene := GEPmap["pgene"].(float64)
	p1p := GEPmap["p1p"].(float64)
	p2p := GEPmap["p2p"].(float64)
	pr := GEPmap["pr"].(float64)
	switch m.BT.SMGEPType {
	case "Genome":
		e := model.New(funcs, functions.Float64, NumGenomes, HeadSize, numGenesPerGenome, numTerminals, numConstants, lincfunc, m.validateFunc)
		s := e.Evolve(Iteration, expectFitness, pm, pis, glis, pris, glris, pgene, p1p, p2p, pr)
		gr, err := grammars.LoadGoMathGrammar()
		// log.Fatal().Str("Grammar", gr.).Msg("Check NOW!!!!!!")
		if err != nil {
			fmt.Printf("unable to load grammar: %v", err)
		}
		// output the va to csv file
		m.outputVAvalues(s, gr)

	case "GenomeSet":
		e := model.NewGS(funcs, functions.Float64, NumGenomeSet, HeadSize, numGenomesPerGenomeSet, numGenesPerGenome, numTerminals, numConstants, lincfunc, m.validateFuncGS)
		s := e.EvolveGS(Iteration, expectFitness, pm, pis, glis, pris, glris, pgene, p1p, p2p, pr)
		gr, err := grammars.LoadGoMathGrammar()
		if err != nil {
			fmt.Printf("unable to load grammar: %v", err)
		}
		for _, v := range s.Genomes {
			v.WriteExps(os.Stdout, gr, m.BT.SIndiNames)
		}
	}
}
