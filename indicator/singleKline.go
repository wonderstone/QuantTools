// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Zhangweixuan (Digital Office Product Department #2)
// revisor:

package indicator

type SingleKline struct {
	Name                   string
	ParSlice               []int
	InfoSlice              []string
	amplitude, para, large float64
	// public parameters. these four would be used in other indicator
	Open, Close, High, Low float64
}

func NewSingleKline(Name string, ParSlice []int, InfoSlice []string) *SingleKline {
	return &SingleKline{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: InfoSlice,
	}
}

func (s *SingleKline) LoadData(data map[string]float64) {
	s.Open = data[s.InfoSlice[2]]
	s.Close = data[s.InfoSlice[3]]
	s.High = data[s.InfoSlice[0]]
	s.Low = data[s.InfoSlice[1]]

	if s.Close >= s.Open {
		s.para = 1
	} else {
		s.para = -1
		s.Open = data[s.InfoSlice[3]]
		s.Close = data[s.InfoSlice[2]]
	}

}

func (s *SingleKline) Eval() float64 {

	s.amplitude = (s.Close - s.Open) / s.Open
	if s.amplitude == 0 {
		s.large = 0
	} else if s.amplitude < 0.006 {
		s.large = 1
	} else if s.amplitude < 0.016 {
		s.large = 2
	} else if s.amplitude < 0.035 {
		s.large = 3
	} else {
		s.large = 4
	}

	var status float64
	if s.High == s.Close && s.Low == s.Open {
		status = 2
	} else if s.High == s.Close && s.Low != s.Open {
		status = 3
	} else if s.High != s.Close && s.Low == s.Open {
		status = 4
	} else {
		status = 1
	}
	return s.para * (s.large*4 + status)
}

func (s *SingleKline) GetName() string {
	return s.Name
}
