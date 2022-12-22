// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Wonderstone (Digital Office Product Department #2)
// revisor:

package indicator

// use gods to generate the queue
import (
	cb "github.com/wonderstone/QuantTools/indicator/tools"
)

// Var 方差指标
type Var struct {
	// Period int
	ParSlice      []int
	sum, lv, mean float64

	InfoSlice []string

	DQ *cb.Queue
	MA *MA
}

// NewVar returns a new Variance indicator
func NewVar(ParSlice []int, infoslice []string) *Var {
	tmpma := NewMA(ParSlice, infoslice)
	return &Var{
		// Period: period,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        tmpma.DQ,
		MA:        tmpma,
	}
}

// LoadData loads 1 tick info datas into the indicator
func (v *Var) LoadDataOld(data map[string]float64) {
	v.MA.LoadData(data)
}

func (v *Var) LoadData(data map[string]float64) {
	tmp_endval, _ := v.DQ.Vals[v.DQ.End].(float64)

	if v.DQ.Full() {
		v.lv = tmp_endval * tmp_endval
		v.sum -= v.lv
	}
	v.MA.LoadData(data)
	v.mean = v.MA.Eval()

	v.sum += (data[v.InfoSlice[0]]) * (data[v.InfoSlice[0]])

}

// Eval evaluates the indicator
func (v *Var) EvalOld() float64 {
	var sum float64
	avg := v.MA.Eval()
	for _, val := range v.DQ.Values() {
		sum += (val.(float64) - avg) * (val.(float64) - avg)
	}
	return sum / float64(v.DQ.Size())
}
func (v *Var) Eval() float64 {
	return (v.sum - float64(v.DQ.Size())*v.mean*v.mean) / float64(v.DQ.Size())
}
