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

type BullishHikkake struct {
	Name      string
	ParSlice  []int
	InfoSlice []string
	single    *SingleKline
	DQ        *cb.Queue
	// using indicator CYE to assess the trend
	cye *CYE
}

func NewBullishHikkake(Name string, ParSlice []int, InfoSlice []string) *BullishHikkake {
	return &BullishHikkake{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: InfoSlice,
		single:    NewSingleKline(Name, ParSlice, InfoSlice),
		cye:       NewCYE(Name, []int{5}, []string{"Close"}),
		DQ:        cb.New(5),
	}
}

func (h *BullishHikkake) LoadData(data map[string]float64) {
	h.cye.LoadData(data)
	h.single.LoadData(data)
	h.DQ.Enqueue(*h.single)
}

func (h *BullishHikkake) Eval() bool {
	if h.DQ.Size() < 5 {
		return false
	}

	start := h.DQ.Vals[h.DQ.Start].(SingleKline)
	mid1 := h.DQ.Vals[(h.DQ.Start+1)%5].(SingleKline)
	mid2 := h.DQ.Vals[(h.DQ.Start+2)%5].(SingleKline)
	mid3 := h.DQ.Vals[(h.DQ.Start+3)%5].(SingleKline)
	end := h.DQ.Vals[(h.DQ.Start+4)%5].(SingleKline)

	if h.cye.Eval() <= 0 && start.Eval() <= -13 && end.Eval() <= -13 &&
		mid1.Eval() < 13 && mid1.Eval() >= 5 && mid1.Close < start.Open &&
		mid2.Eval() < 13 && mid2.Eval() >= 5 && mid2.Close < start.Open &&
		mid3.Eval() < 13 && mid3.Eval() >= 5 && mid3.Close < start.Open {
		return true
	}
	return false

}

func (h *BullishHikkake) GetName() string {
	return h.Name
}
