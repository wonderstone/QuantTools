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
)

// MFI 资金流量指标
type MFI struct {
	Name      string
	ParSlice  []int    //period
	InfoSlice []string // C H L V
	DQ        *tools.Queue
}

// NewMFI  returns a new MFI indicator
func NewMFI(Name string, ParSlice []int, infoslice []string) *MFI {
	return &MFI{
		//period=N+1
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        tools.New(ParSlice[0]),
	}
}

// GetName returns the name of the indicator
func (m *MFI) GetName() string {
	return m.Name
}

// LoadData loads 1 tick info datas into the indicator
func (m *MFI) LoadData(data map[string]float64) {
	m.DQ.Enqueue(data)
}

// Eval evaluates the indicator
func (m *MFI) Eval() float64 {
	var pmf, nmf []float64
	prermf := 0.0
	for _, v := range m.DQ.Values() {
		bar := v.(map[string]float64)
		rmf := bar[m.InfoSlice[3]] * (bar[m.InfoSlice[1]] + bar[m.InfoSlice[2]] + bar[m.InfoSlice[0]]) / 3.0
		if rmf > prermf {
			pmf = append(pmf, rmf)
			nmf = append(nmf, 0)
		} else {
			pmf = append(pmf, 0)
			nmf = append(nmf, rmf)
		}
		prermf = rmf
	}
	mfr := tools.Sum1(pmf) / tools.Sum1(nmf)

	return 100 - 100/(1+mfr)
}
