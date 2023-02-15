// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Yun Jinpeng (Digital Office Product Department #2)
// revisor:

package indicator

// use gods to generate the DQ
import (
	"github.com/wonderstone/QuantTools/indicator/tools"

	"math"
)

// ASI 振动升降指标  同花顺
type ASI struct {
	Name      string
	ParSlice  []int    //period1 period2
	InfoSlice []string //O C H L
	DQ        *tools.Queue
}

// NewASI  returns a new ASI indicator
func NewASI(Name string, ParSlice []int, infoslice []string) *ASI {
	return &ASI{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        tools.New(ParSlice[0] + ParSlice[1] - 1),
	}
}

// GetName returns the name of the indicator
func (a *ASI) GetName() string {
	return a.Name
}

// LoadData loads 1 tick info datas into the indicator
func (a *ASI) LoadData(data map[string]float64) {
	a.DQ.Enqueue(data)
}

// Eval evaluates the indicator
func (a *ASI) Eval() float64 {
	preBar := map[string]float64{"Open": 0.0, "Close": 0.0, "High": 0.0, "Low": 0.0, "Vol": 0.0} //暂存上个周期的数据
	var si []float64
	for _, v := range a.DQ.Values() {
		lc := preBar[a.InfoSlice[1]]
		aa := math.Abs(v.(map[string]float64)[a.InfoSlice[2]] - lc)
		bb := math.Abs(v.(map[string]float64)[a.InfoSlice[3]] - lc)
		cc := math.Abs(v.(map[string]float64)[a.InfoSlice[2]] - preBar[a.InfoSlice[3]])
		dd := math.Abs(lc - preBar["Open"])
		preBar := v.(map[string]float64)
		var r float64
		if aa > bb && aa > cc {
			r = aa + bb/2 + dd/4
		} else {
			if bb > cc && bb > aa {
				r = bb + aa/2 + dd/4
			} else {
				r = cc + dd/4
			}
		}
		x := v.(map[string]float64)[a.InfoSlice[1]] - lc + (v.(map[string]float64)[a.InfoSlice[1]]-v.(map[string]float64)[a.InfoSlice[0]])/2 + lc - preBar["Open"]
		si = append(si, 16*x/r*math.Max(aa, bb))
	}
	asi := tools.SumN(si, a.ParSlice[0])
	astt := tools.AvgN(asi, a.ParSlice[1])
	return astt[len(astt)-1]
}
