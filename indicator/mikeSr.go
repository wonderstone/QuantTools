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

type MIKESR struct {
	Name         string
	ParSlice     []int
	InfoSlice    []string
	max, min     float64
	maxDQ, minDQ *cb.Queue
}

func NewMIKESR(Name string, ParSlice []int, InfoSlice []string) *MIKESR {
	return &MIKESR{
		Name:      Name,
		ParSlice:  ParSlice, //period
		InfoSlice: InfoSlice,
		maxDQ:     cb.New(ParSlice[0]),
		minDQ:     cb.New(ParSlice[0]),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (s *MIKESR) LoadData(data map[string]float64) {
	s.maxDQ.Enqueue(data[s.InfoSlice[0]])
	s.minDQ.Enqueue(data[s.InfoSlice[1]])

	s.max = 0
	for i := 0; i < s.maxDQ.Size(); i++ {
		s.max = math.Max(s.max, s.maxDQ.Vals[i].(float64))
	}

	s.min = math.Inf(0)
	for i := 0; i < s.minDQ.Size(); i++ {
		s.min = math.Min(s.min, s.minDQ.Vals[i].(float64))
	}
}

func (s *MIKESR) Eval() float64 {
	return 2*s.max - s.min
}

func (s *MIKESR) GetName() string {
	return s.Name
}
