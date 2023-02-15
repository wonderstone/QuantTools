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

// KVO Klinger Oscillator
type KVO struct {
	Name      string
	ParSlice  []int    //period
	InfoSlice []string // C H L V
	DQ        *tools.Queue
}

// NewKVO  returns a new EOM indicator
func NewKVO(Name string, ParSlice []int, infoslice []string) *KVO {
	return &KVO{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        tools.New(ParSlice[0]),
	}
}

// GetName returns the name of the indicator
func (k *KVO) GetName() string {
	return k.Name
}

// LoadData loads 1 tick info datas into the indicator
func (k *KVO) LoadData(data map[string]float64) {
	k.DQ.Enqueue(data)
}

// Eval evaluates the indicator
func (k *KVO) Eval() float64 {
	var vf []float64
	preT := -1
	preCm := 0.0
	preBar := k.DQ.Values()[0].(map[string]float64)
	for i, v := range k.DQ.Values() {
		bar := v.(map[string]float64)
		T := -1
		if bar[k.InfoSlice[1]]+bar[k.InfoSlice[2]]+bar[k.InfoSlice[0]] > preBar[k.InfoSlice[1]]+preBar[k.InfoSlice[2]]+preBar[k.InfoSlice[0]] {
			T = 1
		}
		dm := bar[k.InfoSlice[1]] - bar[k.InfoSlice[2]]
		preDm := preBar[k.InfoSlice[1]] - preBar[k.InfoSlice[2]]
		cm := preCm
		if i != 0 {
			if T == preT {
				cm = preCm + dm
			} else {
				cm = preDm + dm
			}
		} else {
			cm = dm
		}
		preCm = cm
		vfElem := bar[k.InfoSlice[3]] * 2 * ((dm / cm) - 1) * float64(T) * 100
		vf = append(vf, vfElem)
		preBar = bar
	}

	return tools.Ema(vf, 34) - tools.Ema(vf, 55)
}
