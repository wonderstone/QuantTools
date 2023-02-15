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

type WVAD struct {
	Name      string
	ParSlice  []int
	InfoSlice []string
	sum       float64
	DQ        *cb.Queue
}

func NewWVAD(Name string, ParSlice []int, InfoSlice []string) *WVAD {
	return &WVAD{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: InfoSlice,
		sum:       0,
		DQ:        cb.New(ParSlice[0]),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (w *WVAD) LoadData(data map[string]float64) {
	_open := data[w.InfoSlice[0]]
	_close := data[w.InfoSlice[1]]
	_high := data[w.InfoSlice[2]]
	_low := data[w.InfoSlice[3]]
	_vol := data[w.InfoSlice[4]]
	if w.DQ.Full() {
		w.sum -= w.DQ.Vals[w.DQ.End].(float64)
	}
	wvad := (_close - _open) / (_high - _low) * _vol
	w.DQ.Enqueue(wvad)
	w.sum += wvad
}

func (w *WVAD) Eval() float64 {
	return w.sum
}

func (w *WVAD) GetName() string {
	return w.Name
}
