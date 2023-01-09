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

type RSI struct {
	Name                 string
	ParSlice             []int
	lclose, up, down, lv float64
	InfoSlice            []string
	// store the difference of closing price between today and yesterday
	difDQ *cb.Queue
}

func NewRSI(Name string, ParSlice []int, InfoSlice []string) *RSI {
	return &RSI{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: InfoSlice,
		difDQ:     cb.New(ParSlice[0]),
		up:        0,
		down:      0,
	}
}

func (r *RSI) LoadData(data map[string]float64) {
	if r.difDQ.Empty() {
		r.lclose = data[r.InfoSlice[0]]
	}

	dif := data[r.InfoSlice[0]] - r.lclose
	r.up += math.Max(dif, 0)
	r.down += math.Abs(math.Min(dif, 0))

	if r.difDQ.Full() {
		r.lv = r.difDQ.Vals[r.difDQ.End].(float64)
	}

	r.difDQ.Enqueue(dif)
	r.lclose = data[r.InfoSlice[0]]

	if r.difDQ.Full() {
		r.up -= math.Max(r.lv, 0)
		r.down -= math.Abs(math.Min(r.lv, 0))
	}
}

func (r *RSI) Eval() float64 {
	if r.up+r.down == 0 {
		return 50
	}
	return 100 * r.up / (r.up + r.down)
}

func (r *RSI) GetName() string {
	return r.Name
}
