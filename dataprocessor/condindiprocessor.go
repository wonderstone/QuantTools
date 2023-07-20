package dataprocessor

import (
	"errors"

	"github.com/wonderstone/QuantTools/indicator"
)

// Condindi Chain for realtime data
type CondindiChain struct {
	CICMap map[string][]indicator.CondIndiInfo
}

// new a CondindiChain
func NewCondindiChain(CICMap map[string][]indicator.CondIndiInfo) CondindiChain {
	return CondindiChain{CICMap}
}

// new a CondindiChain with yaml file
func NewCondindiChainWithYaml(path string, instIDS []string, infomap map[string]interface{}) CondindiChain {
	// new a CICMap and iter the instIDS to add the CondindiChain
	CICMap := make(map[string][]indicator.CondIndiInfo)
	for _, v := range instIDS {
		CICMap[v] = indicator.GetCondIndiInfoSlice(path, infomap)
	}
	return CondindiChain{CICMap}
}

// method to run the Condindi Chain with BarC in and out
func (c *CondindiChain) Run(bar *BarC) (bool, error) {
	// iter the bar stockdata to update the CondindiChain
	for k := range bar.Stockdata {
		// iter the CondindiChain to get the CondIndiInfo
		for k1 := range c.CICMap[k] {
			// check the CondIndiInfo's CL
			if !c.CICMap[k][k1].CL.Judge(bar.Stockdata[k].BarTime) {
				c.CICMap[k][k1].II.LoadData(bar.Stockdata[k].IndiDataMap)
				tmpn := c.CICMap[k][k1].II.GetName()
				tmpv := c.CICMap[k][k1].II.Eval()
				bar.Stockdata[k].IndiDataMap[tmpn] = tmpv
				return true, nil
			}
		}

	}
	return false, errors.New("no CondIndiInfo is updated")
}
