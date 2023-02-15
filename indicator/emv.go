// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Zhangweixuan (Digital Office Product Department #2)
// revisor:

package indicator

import cb "github.com/wonderstone/QuantTools/indicator/tools"

type EMV struct {
	Name      string
	ParSlice  []int
	InfoSlice []string
	lv, sum   float64
	DQ        *cb.Queue
}

func NewEMV(Name string, ParSlice []int, InfoSlice []string) *EMV {
	return &EMV{
		Name:      Name,
		ParSlice:  ParSlice, //period
		InfoSlice: InfoSlice,
		DQ:        cb.New(ParSlice[0]),
		lv:        0,
		sum:       0,
	}
}

// LoadData loads 1 tick info datas into the indicator
func (e *EMV) LoadData(data map[string]float64) {
	tmpH := data[e.InfoSlice[0]]
	tmpL := data[e.InfoSlice[1]]
	tmpV := data[e.InfoSlice[2]]

	A := (tmpH + tmpL) / 2
	B := e.lv
	C := tmpH - tmpL
	e.lv = A

	if e.lv == 0 {
		return
	}
	if e.DQ.Full() {
		e.sum -= e.DQ.Vals[e.DQ.End].(float64)
	}
	emv := (A - B) * C / tmpV
	e.DQ.Enqueue(emv)
	e.sum += emv
}

func (e *EMV) Eval() float64 {
	return e.sum
}

func (e *EMV) GetName() string {
	return e.Name
}
