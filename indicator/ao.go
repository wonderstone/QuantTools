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

type AO struct {
	ParSlice  []int
	InfoSlice []string //[Low,High]
	low, high float64
	DQM       *cb.Queue
	Ma        *MA
}

func NewAO(ParSlice []int, infoslice []string) *AO {
	return &AO{
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQM:       cb.New(ParSlice[0]),
		Ma:        NewMA(ParSlice, infoslice),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (a *AO) LoadData(data map[string]float64) {
	a.low = data[a.InfoSlice[0]]
	a.high = data[a.InfoSlice[1]]
	a.DQM.Enqueue((a.high + a.low) / 2)
}

// Eval evaluates the indicator
func (a *AO) Eval() float64 {
	a.Ma.DQ = a.DQM
	a.Ma.ParSlice[0] = 5
	Ma5 := a.Ma.Eval()
	a.Ma.DQ = a.DQM
	a.Ma.ParSlice[0] = 34
	Ma34 := a.Ma.Eval()
	ao := Ma5 - Ma34
	return ao
}
