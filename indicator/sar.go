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

type SAR struct {
	Name                      string
	ParSlice                  []int
	InfoSlice                 []string
	af, ep, sar, max, min, lv float64
	upTrend                   bool
	highDQ, lowDQ             *cb.Queue
	zig                       *ZIG
}

func NewSAR(Name string, ParSlice []int, InfoSlice []string) *SAR {
	return &SAR{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: InfoSlice,
		highDQ:    cb.New(ParSlice[0]),
		lowDQ:     cb.New(ParSlice[0]),
		zig:       NewZIG("Name", []int{3, 1}, []string{"Open", "High", "Low", "Close"}),
		lv:        0,
		sar:       0,
		af:        0.02,
		max:       0,
		min:       math.Inf(0),
		upTrend:   true,
	}
}

func (s *SAR) LoadData(data map[string]float64) {
	// load data
	lMax := s.max
	lMin := s.min
	lTrend := s.upTrend

	s.min = data[s.InfoSlice[1]]
	s.max = data[s.InfoSlice[0]]

	s.zig.LoadData(data)
	s.highDQ.Enqueue(s.max)
	s.lowDQ.Enqueue(s.min)

	if s.highDQ.Size() < s.ParSlice[0] {
		s.max = math.Max(s.max, lMax)
		s.min = math.Min(s.min, lMin)
		s.lv = s.zig.Eval()
		return
	}

	// assess trends
	if s.lv <= s.zig.Eval() {
		s.upTrend = true
	} else {
		s.upTrend = false
	}
	s.lv = s.zig.Eval()
	// initialize
	if s.sar == 0 {
		if s.upTrend {
			s.sar = s.min
			s.ep = s.max
		} else {
			s.sar = s.max
			s.ep = s.min
		}
		return
	}
	// step
	s.sar += s.af * (s.ep - s.sar)
	if s.upTrend && s.sar > lMin {
		s.sar = lMin
	} else if !s.upTrend && s.sar < lMax {
		s.sar = lMax
	}
	// update ep
	if s.upTrend {
		s.ep = s.max
	} else {
		s.ep = s.min
	}
	// update af
	if s.upTrend != lTrend {
		s.af = 0.02
	} else if (s.upTrend && s.max > lMax) || (!s.upTrend && s.min < lMin) {
		s.af = math.Min(s.af+0.02, 0.2)
	}
}

func (s *SAR) Eval() float64 {
	if s.sar == 0 {
		return math.NaN()
	}
	return s.sar
}

func (s *SAR) GetName() string {
	return s.Name
}
