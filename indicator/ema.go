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

// EMA is the moving average indicator
type EMA struct {
	Name                       string
	ParSlice                   []int
	ptoday, lv, k, result, res float64
	// info fields for indicator calculation
	InfoSlice []string
	DQ        *cb.Queue
}

// NewMA returns a new MA indicator
func NewEMA(Name string, ParSlice []int, infoslice []string) *EMA {
	return &EMA{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		k:         2.0 / (float64(ParSlice[0]) + 1.0),
		DQ:        cb.New(ParSlice[0]),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (e *EMA) LoadData(data map[string]float64) {
	e.ptoday = data[e.InfoSlice[0]]
	e.DQ.Enqueue(data[e.InfoSlice[0]])
	if e.DQ.Size() == 1 {
		e.lv = e.ptoday
	}
	if e.DQ.Size() > 1 {
		e.result = e.k*e.ptoday + (1-e.k)*e.lv
		e.lv = e.result
	}
	e.res = e.result
}

// Eval evaluates the indicator
func (e *EMA) Eval() float64 {
	return e.res
}
func (e *EMA) GetName() string {
	return e.Name
}
