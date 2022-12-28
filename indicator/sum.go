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

type Sum struct {
	ParSlice           []int
	p, sum, period, lv float64

	// info fields for indicator calculation
	InfoSlice []string
	DQ        *cb.Queue
}

func NewSum(ParSlice []int, infoslice []string) *Sum {
	return &Sum{
		ParSlice:  ParSlice, //period
		InfoSlice: infoslice,
		period:    float64(ParSlice[0]),
		DQ:        cb.New(ParSlice[0]),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (s *Sum) LoadData(data map[string]float64) {
	s.p = data[s.InfoSlice[0]]
	s.sum += data[s.InfoSlice[0]]
	if s.DQ.Full() {
		s.lv = s.DQ.Vals[s.DQ.End].(float64)
	}

	s.DQ.Enqueue(data[s.InfoSlice[0]])
	if s.DQ.Full() {
		s.sum -= s.lv
	}
}

func (s *Sum) Eval() float64 {
	result := s.sum
	return result
}
