// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Zhangweixuan (Digital Office Product Department #2)
// revisor:
package indicator

import (
	cb "github.com/wonderstone/QuantTools/indicator/tools"
)

type VR struct {
	Name               string
	ParSlice           []int
	InfoSlice          []string
	lv, upSum, downSum float64
	statusDQ, DQ       *cb.Queue
}

func NewVR(Name string, ParSlice []int, infoslice []string) *VR {
	return &VR{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        cb.New(ParSlice[0]),
		statusDQ:  cb.New(ParSlice[0]),
		lv:        0,
		upSum:     0,
		downSum:   0,
	}
}

// LoadData loads 1 tick info datas into the indicator
func (v *VR) LoadData(data map[string]float64) {
	if v.DQ.Full() {
		status := v.statusDQ.Vals[v.statusDQ.End].(int)
		if status == 2 {
			v.upSum -= v.DQ.Vals[v.DQ.End].(float64)
		} else if status == 0 {
			v.downSum -= v.DQ.Vals[v.DQ.End].(float64)
		} else {
			v.upSum -= v.DQ.Vals[v.DQ.End].(float64) / 2
			v.downSum -= v.DQ.Vals[v.DQ.End].(float64) / 2
		}
	}

	if data[v.InfoSlice[0]] > v.lv {
		v.statusDQ.Enqueue(2)
		v.upSum += data[v.InfoSlice[1]]
	} else if data[v.InfoSlice[0]] < v.lv {
		v.statusDQ.Enqueue(0)
		v.downSum += data[v.InfoSlice[1]]
	} else {
		v.statusDQ.Enqueue(1)
		v.downSum += data[v.InfoSlice[1]] / 2
		v.upSum += data[v.InfoSlice[1]] / 2
	}

	v.DQ.Enqueue(data[v.InfoSlice[1]])

	v.lv = data[v.InfoSlice[0]]
}

// Eval evaluates the indicator
func (v *VR) Eval() float64 {
	return v.upSum / v.downSum
}

func (v *VR) GetName() string {
	return v.Name
}
