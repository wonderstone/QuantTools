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

type BearishBeltHold struct {
	Name      string
	ParSlice  []int
	InfoSlice []string
	single    *SingleKline
	DQ        *cb.Queue
	// using indicator CYE to assess the trend
	cye *CYE
}

func NewBearishBeltHold(Name string, ParSlice []int, InfoSlice []string) *BearishBeltHold {
	return &BearishBeltHold{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: InfoSlice,
		single:    NewSingleKline(Name, ParSlice, InfoSlice),
		cye:       NewCYE(Name, []int{5}, []string{"Close"}),
		DQ:        cb.New(2),
	}
}

func (h *BearishBeltHold) LoadData(data map[string]float64) {
	h.cye.LoadData(data)
	h.single.LoadData(data)
	h.DQ.Enqueue(*h.single)
}

func (h *BearishBeltHold) Eval() bool {
	if h.DQ.Size() < 2 {
		return false
	}

	start := h.DQ.Vals[h.DQ.Start].(SingleKline)
	end := h.DQ.Vals[(h.DQ.Start+1)%2].(SingleKline)

	if h.cye.Eval() > 0 && start.Eval() >= 9 &&
		end.Eval() <= -13 && end.Open > start.Close &&
		end.High == end.Open && end.Low < end.High {
		return true
	}
	return false

}

func (h *BearishBeltHold) GetName() string {
	return h.Name
}
