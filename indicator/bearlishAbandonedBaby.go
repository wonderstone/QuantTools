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

type BearishAbandonedBaby struct { 
	Name      string
	ParSlice  []int
	InfoSlice []string
	single    *SingleKline
	DQ        *cb.Queue
	// using indicator CYE to assess the trend
	cye *CYE
}

func NewBearishAbandonedBaby(Name string, ParSlice []int, InfoSlice []string) *BearishAbandonedBaby {
	return &BearishAbandonedBaby{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: InfoSlice,
		single:    NewSingleKline(Name, ParSlice, InfoSlice),
		cye:       NewCYE(Name, []int{5}, []string{"Close"}),
		DQ:        cb.New(3),
	}
}

func (h *BearishAbandonedBaby) LoadData(data map[string]float64) {
	h.cye.LoadData(data)
	h.single.LoadData(data)
	h.DQ.Enqueue(*h.single)
}

func (h *BearishAbandonedBaby) Eval() bool {
	if h.DQ.Size() < 3 {
		return false
	}

	start := h.DQ.Vals[h.DQ.Start].(SingleKline)
	mid := h.DQ.Vals[(h.DQ.Start+1)%3].(SingleKline)
	end := h.DQ.Vals[(h.DQ.Start+2)%3].(SingleKline)

	if h.cye.Eval() >= 0 && start.Eval() >= 9 &&
		math.Abs(mid.Eval()) <= 8 &&
		end.Eval() <= -9 && end.Close < start.Close &&
		math.Min(mid.Open, mid.Close) > math.Max(start.Close, end.Open) {
		return true
	}
	return false

}

func (h *BearishAbandonedBaby) GetName() string {
	return h.Name
}
