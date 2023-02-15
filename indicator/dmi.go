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

// DMI 趋向指标
type DMI struct {
	Name      string
	ParSlice  []int    // len 0; N 1 ;M 2    满足:M<N,len>N
	InfoSlice []string //O C H L
	DQ        *tools.Queue
}

// NewDMI  returns a new DMI indicator
func NewDMI(Name string, ParSlice []int, infoslice []string) *DMI {
	return &DMI{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        tools.New(ParSlice[0]),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (m *DMI) LoadData(data map[string]float64) {
	m.DQ.Enqueue(data[m.InfoSlice[0]])
}

// Eval evaluates the indicator
func (m *DMI) Eval() float64 {
	//preBar := bars.Bar{} //暂存上个周期的数据
	//var cmf []float64
	//var vol []float64
	if float64(m.DQ.Size()) < math.Max(float64(m.ParSlice[1]), float64(m.ParSlice[2])) {
		return 0
	}
	_, C, H, L, _ := tools.GetArrayFromMapList(m.DQ)

	t1 := tools.VecMax(tools.VecMax(tools.VecSub(H, L), tools.VecAbs(tools.VecSub(H, tools.VecRef1(C)))), tools.VecAbs(tools.VecSub(tools.VecRef1(C), L)))
	mtr := tools.SumN(t1, m.ParSlice[1])
	hd := tools.VecSub(H, tools.VecRef1(H))
	ld := tools.VecSub(tools.VecRef1(L), L)
	dmp := hd
	dmm := ld
	//IF(HD>0&&HD>LD,HD,0)
	for i := 0; i < len(hd); i++ {
		if hd[i] > 0 && hd[i] > ld[i] {
			dmp[i] = hd[i]
		} else {
			dmp[i] = 0
		}
	}
	dmp = tools.SumN(dmp, m.ParSlice[1])
	for i := 0; i < len(ld); i++ {
		if ld[i] > 0 && ld[i] > hd[i] {
			dmm[i] = ld[i]
		} else {
			dmm[i] = 0
		}
	}
	dmm = tools.SumN(dmm, m.ParSlice[1])
	pdi := tools.VecDiv(dmp, mtr[len(mtr)-len(dmp):]) //
	mdi := tools.VecDiv(dmm, mtr[len(mtr)-len(dmp):]) //
	var adx []float64
	adx = append(adx, pdi...)
	for i := 0; i < len(pdi); i++ {
		adx[i] = math.Abs(mdi[i]-pdi[i]) / (mdi[i] + pdi[i]) * 100
	}
	r := tools.SumN(adx, m.ParSlice[2])
	return r[len(r)-1] / float64(m.ParSlice[2])

}
