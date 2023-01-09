// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Maminghui (Digital Office Product Department #2)
// revisor:

package indicator

import (
	cb "github.com/wonderstone/QuantTools/indicator/tools"
)

type Qstick struct {
	Name string
	ParSlice                          []int
	closing, opening, period, sum, lv float64

	// info fields for indicator calculation
	InfoSlice []string
	MA        *MA
	DQ        *cb.Queue
}

// NewMA returns a new MA indicator
func NewQstick(Name string,ParSlice []int, infoslice []string) *Qstick {
	return &Qstick{
		Name: Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice, //[closing,opening]
		period:    float64(ParSlice[0]),
		DQ:        cb.New(ParSlice[0]),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (qs *Qstick) LoadData(data map[string]float64) {
	qs.closing = data[qs.InfoSlice[0]]
	qs.opening = data[qs.InfoSlice[1]]
	qs.sum += qs.closing - qs.opening
	if qs.DQ.Full() {
		qs.lv = qs.DQ.Vals[qs.DQ.End].(float64)
	}

	qs.DQ.Enqueue(qs.closing - qs.opening)
	if qs.DQ.Full() {
		qs.sum -= qs.lv
	}
}
func (qs *Qstick) Eval() float64 {

	return qs.sum / qs.period
}
func (qs *Qstick) GetName() string {
	return qs.Name
}