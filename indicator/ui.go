// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Maminghui (Digital Office Product Department #2)
// revisor:
// The Ulcer Index (UI) measures downside risk. The index increases in value
// as the price moves farther away from a recent high and falls as the price
// rises to new highs.
//
// High Closings = Max(period, Closings)
// Percentage Drawdown = 100 * ((Closings - High Closings) / High Closings)
// Squared Average = Sma(period, Percent Drawdown * Percent Drawdown)
// Ulcer Index = Sqrt(Squared Average)
//
// Returns ui.
package indicator

import (
	"math"

	cb "github.com/wonderstone/QuantTools/indicator/tools"
)

type UI struct {
	Name                                              string
	ParSlice                                          []int //period
	period                                            int
	closing, highClosing, lv, percentageDrawdown, sum float64
	uires, squaredAverage                             []float64
	// info fields for indicator calculation
	InfoSlice []string //closing
	DQ        *cb.Queue
	MA        *MA
}

func NewUI(Name string, ParSlice []int, infoslice []string) *UI {
	return &UI{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        cb.New(ParSlice[0]),
		MA:        NewMA(Name, ParSlice, infoslice),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (ui *UI) LoadData(data map[string]float64) {
	ui.period = ui.ParSlice[0]
	ui.closing = data[ui.InfoSlice[0]]
	ui.highClosing = math.Max(ui.closing, ui.highClosing)

	ui.percentageDrawdown = ((ui.closing - ui.highClosing) / ui.highClosing) * 100
	if ui.DQ.Full() {
		ui.lv = ui.DQ.Vals[ui.DQ.End].(float64)
	}

	ui.DQ.Enqueue(ui.percentageDrawdown)
	if ui.DQ.Full() {
		ui.sum -= ui.lv
	}
	ui.MA.ParSlice[0] = ui.period
	ui.MA.InfoSlice = []string{"Closing"}
	ui.MA.LoadData(map[string]float64{"Closing": ui.percentageDrawdown * ui.percentageDrawdown})

}

// Eval evaluates the indicator
func (ui *UI) Eval() float64 {

	ui.squaredAverage = append(ui.squaredAverage, ui.MA.Eval())
	ui.uires = sqrt(ui.squaredAverage)
	return ui.uires[len(ui.uires)-1]
}
func (ui *UI) GetName() string {
	return ui.Name
}
func sqrt(values []float64) []float64 {
	result := make([]float64, len(values))

	for i := 0; i < len(values); i++ {
		result[i] = math.Sqrt(values[i])
	}

	return result
}
