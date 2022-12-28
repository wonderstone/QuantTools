// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Maminghui (Digital Office Product Department #2)
// revisor:
package indicator

import (
	cb "github.com/wonderstone/QuantTools/indicator/tools"
	//"fmt"
)

// EMA is the moving average indicator
type EMA struct {
	ParSlice           []int
	ptoday, lv, k, sum float64

	// info fields for indicator calculation
	InfoSlice []string
	DQ        *cb.Queue
}

func NewEMA(ParSlice []int, infoslice []string) *EMA {
	return &EMA{
		ParSlice:  ParSlice, //period
		InfoSlice: infoslice,
		k:         2.0 / (float64(ParSlice[0]) + 1.0),
		DQ:        cb.New(ParSlice[0]),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (e *EMA) LoadData(data map[string]float64) {
	e.ptoday = data[e.InfoSlice[0]] //输入的最后一个值
	e.sum += data[e.InfoSlice[0]]
	if e.DQ.Full() {
		e.lv = e.DQ.Vals[e.DQ.End].(float64)
	}

	e.DQ.Enqueue(data[e.InfoSlice[0]])
	if e.DQ.Full() {
		e.sum -= e.lv
	}
}

func (e *EMA) Eval() float64 {
	result := make([]float64, 2)
	result[0] = e.lv
	if result[0] == 0 {
		result[0] = e.ptoday
		result[1] = (e.ptoday * e.k) + (result[0] * float64(1-e.k))
		e.lv = result[1]
	} else {
		result[1] = (e.ptoday * e.k) + (result[0] * float64(1-e.k))
		e.lv = result[1]
	}
	return result[1]
}
