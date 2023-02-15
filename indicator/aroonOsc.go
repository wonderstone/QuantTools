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

// AroonOsc is the AroonOsc Indicator
type AroonOsc struct {
	Name       string
	ReturnName []string
	ParSlice   []int
	InfoSlice  []string // H L
	DQ         *cb.Queue
}

// NewAroonOsc returns a new NewAROON indicator
func NewAroonOsc(Name string, ParSlice []int, infoslice []string) *AroonOsc {
	return &AroonOsc{
		Name:       Name,
		ReturnName: []string{"AroonUp", "AroonDown", "AroonOsc"},
		ParSlice:   ParSlice,
		InfoSlice:  infoslice,
		DQ:         cb.New(ParSlice[0]),
	}
}

// GetName returns the name of the indicator
func (a *AroonOsc) GetName() string {
	return a.Name
}

func (a *AroonOsc) GetReturnName() []string {
	return a.ReturnName
}

func (a *AroonOsc) SetReturnName(ReturnNames []string) {
	a.ReturnName[0] = ReturnNames[0]
	a.ReturnName[1] = ReturnNames[1]
	a.ReturnName[2] = ReturnNames[2]

}

// LoadData loads 1 tick info datas into the indicator
func (a *AroonOsc) LoadData(data map[string]float64) {
	a.DQ.Enqueue(data)
}

// Eval evaluates the indicator
func (a *AroonOsc) Eval() float64 {
	var maxHigh float64 = -math.MaxFloat64
	var maxHighIndex int
	var minLow float64 = math.MaxFloat64
	var minLowIndex int

	for i, v := range a.DQ.Values() {
		if v.(map[string]float64)[a.InfoSlice[0]] >= maxHigh {
			maxHigh = v.(map[string]float64)[a.InfoSlice[0]]
			maxHighIndex = i
		}
		if v.(map[string]float64)[a.InfoSlice[1]] <= minLow {
			minLow = v.(map[string]float64)[a.InfoSlice[1]]
			minLowIndex = i
		}
	}

	var dayNum int = a.DQ.Size()
	AroonUp := float64(dayNum-(dayNum-1-maxHighIndex)) / float64(dayNum) * 100
	AroonDown := float64(dayNum-(dayNum-1-minLowIndex)) / float64(dayNum) * 100
	AroonOsc := AroonUp - AroonDown
	return AroonOsc
}
