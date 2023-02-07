// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Maminghui (Digital Office Product Department #2)
// revisor:
// TR = Max((High - Low), (High - Closing), (Closing - Low))
// ATR = SMA TR
package indicator

import (
	"math"
	cb "github.com/wonderstone/QuantTools/indicator/tools"
)

type Atr struct {
	Name                        string
	ParSlice                    []int //period
	period                      int
	high, low, closing, sum, lv float64
	// info fields for indicator calculation
	InfoSlice []string //high,low,closing
	DQ        *cb.Queue
}

func NewAtr(Name string, ParSlice []int, infoslice []string) *Atr {
	return &Atr{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        cb.New(ParSlice[0]),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (a *Atr) LoadData(data map[string]float64) {
	a.period = a.ParSlice[0]
	a.high = data[a.InfoSlice[0]]
	a.low = data[a.InfoSlice[1]]
	a.closing = data[a.InfoSlice[2]]
	a.sum += math.Max(a.high-a.low, math.Max(a.high-a.low, a.closing-a.low))
	if a.DQ.Full() {
		a.lv = a.DQ.Vals[a.DQ.End].(float64)
	}
	a.DQ.Enqueue(math.Max(a.high-a.low, math.Max(a.high-a.low, a.closing-a.low)))
	if a.DQ.Full() {
		a.sum -= a.lv
	}

}

// Eval evaluates the indicator
func (a *Atr) Eval() float64 {

	return a.sum / float64(a.DQ.Size())
}
func (a *Atr) GetName() string {
	return a.Name
}
