// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Yun Jinpeng (Digital Office Product Department #2)
// revisor:

package indicator

import (
	cb "github.com/wonderstone/QuantTools/indicator/tools"
	"math"
)

// AroonUp is the AroonUp Indicator
type AroonUp struct {
	Name      string
	ParSlice  []int
	InfoSlice []string // H L
	DQ        *cb.Queue
}

// NewAroonUp returns a new NewAROON indicator
func NewAroonUp(Name string, ParSlice []int, infoslice []string) *AroonUp {
	return &AroonUp{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        cb.New(ParSlice[0]),
	}
}

// GetName returns the name of the indicator
func (a *AroonUp) GetName() string {
	return a.Name
}

// LoadData loads 1 tick info datas into the indicator
func (a *AroonUp) LoadData(data map[string]float64) {
	a.DQ.Enqueue(data)
}

// Eval evaluates the indicator
func (a *AroonUp) Eval() float64 {
	var maxHigh float64 = -math.MaxFloat64
	var maxHighIndex int
	var minLow float64 = math.MaxFloat64

	for i, v := range a.DQ.Values() {
		if v.(map[string]float64)[a.InfoSlice[0]] >= maxHigh {
			maxHigh = v.(map[string]float64)[a.InfoSlice[0]]
			maxHighIndex = i
		}
		if v.(map[string]float64)[a.InfoSlice[1]] <= minLow {
			minLow = v.(map[string]float64)[a.InfoSlice[1]]
		}
	}

	var dayNum int = a.DQ.Size()
	AroonUp := float64(dayNum-(dayNum-1-maxHighIndex)) / float64(dayNum) * 100

	return AroonUp
}
