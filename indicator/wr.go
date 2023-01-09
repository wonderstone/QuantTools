// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Zhangweixuan (Digital Office Product Department #2)
// revisor:
package indicator

import (
	cb "github.com/wonderstone/QuantTools/indicator/tools"
	"math"
)

type WR struct {
	Name         string
	ParSlice     []int
	InfoSlice    []string
	close        float64
	maxDQ, minDQ *cb.Queue
}

func NewWR(Name string, ParSlice []int, infoslice []string) *WR {
	return &WR{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		maxDQ:     cb.New(ParSlice[0]),
		minDQ:     cb.New(ParSlice[0]),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (w *WR) LoadData(data map[string]float64) {
	w.close = data[w.InfoSlice[0]]
	w.minDQ.Enqueue(data[w.InfoSlice[1]])
	w.maxDQ.Enqueue(data[w.InfoSlice[2]])
}

// Eval evaluates the indicator
func (w *WR) Eval() float64 {
	var max float64 = 0
	var min float64 = math.Inf(0)
	for i := 0; i < w.maxDQ.Size(); i++ {
		max = math.Max(max, w.maxDQ.Vals[i].(float64))
	}
	for i := 0; i < w.minDQ.Size(); i++ {
		min = math.Min(min, w.minDQ.Vals[i].(float64))
	}
	return (max - w.close) / (max - min) * 100
}

func (w *WR) GetName() string {
	return w.Name
}
