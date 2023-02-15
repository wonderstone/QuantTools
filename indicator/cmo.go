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

// CMO 钱德动量摆动指标 Chande Momentum Oscillato
type CMO struct {
	Name      string
	ParSlice  []int    //period
	InfoSlice []string // C
	DQ        *tools.Queue
}

// NewCMO  returns a new cmf indicator
func NewCMO(Name string, ParSlice []int, infoslice []string) *CMO {
	return &CMO{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        tools.New(ParSlice[0]),
	}
}

// GetName returns the name of the indicator
func (c *CMO) GetName() string {
	return c.Name
}

func (c *CMO) LoadData(data map[string]float64) {
	c.DQ.Enqueue(data)
}

// Eval evaluates the indicator
// 该指标输入n个K线的序列,输出n-1个K线的CMO指标
func (c *CMO) Eval() float64 {
	var cz1 []float64
	var cz2 []float64
	preBar := map[string]float64{"Open": 0.0, "Close": 0.0, "High": 0.0, "Low": 0.0, "Vol": 0.0} //暂存上个周期的数据
	for _, v := range c.DQ.Values() {
		bar := v.(map[string]float64)
		cz1_ := 0.0
		if bar[c.InfoSlice[0]]-preBar[c.InfoSlice[0]] > 0 {
			cz1_ = bar[c.InfoSlice[0]] - preBar[c.InfoSlice[0]]
		} else {
			cz1_ = 0
		}
		cz1 = append(cz1, cz1_)

		cz2_ := 0.0
		if bar[c.InfoSlice[0]]-preBar[c.InfoSlice[0]] < 0 {
			cz2_ = math.Abs(bar[c.InfoSlice[0]] - preBar[c.InfoSlice[0]])
		} else {
			cz2_ = 0
		}
		cz2 = append(cz2, cz2_)

		preBar = v.(map[string]float64)
	}
	su := tools.Sum1(cz1[1:]) //cost 10ns
	sd := tools.Sum1(cz2[1:])
	//su := 1.1
	//sd := 2.1
	cmo := (su - sd) / (su + sd) * 100

	return cmo
}
