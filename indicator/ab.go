// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Maminghui (Digital Office Product Department #2)
// revisor:
// Acceleration Bands. Plots upper and lower envelope bands
// around a simple moving average.
//
// Upper Band = SMA(High * (1 + 4 * (High - Low) / (High + Low)))
// Middle Band = SMA(Closing)
// Lower Band = SMA(Low * (1 - 4 * (High - Low) / (High + Low)))
//
// Returns upper band, middle band, lower band.
package indicator

import (
	cb "github.com/wonderstone/QuantTools/indicator/tools"
)

type AB struct {
	Name                  string
	ParSlice              []int
	InfoSlice             []string //[high,low,closing]
	high, low, closing, k float64
	DQH, DQL, DQM         *cb.Queue
	Ma                    *MA
}

func NewAB(Name string, ParSlice []int, infoslice []string) *AB {
	return &AB{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQH:       cb.New(ParSlice[0]),
		DQL:       cb.New(ParSlice[0]),
		DQM:       cb.New(ParSlice[0]),
		Ma:        NewMA("tempMA", ParSlice, infoslice),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (a *AB) LoadData(data map[string]float64) {
	a.high = data[a.InfoSlice[0]]
	a.low = data[a.InfoSlice[1]]
	a.closing = data[a.InfoSlice[2]]
	a.k = (a.high - a.low) / (a.high + a.low)
	a.DQH.Enqueue((4*a.k + 1) * a.high)
	a.DQL.Enqueue((1 - 4*a.k) * a.low)
	a.DQL.Enqueue(a.closing)
}

func (a *AB) Eval() float64 {

	upperBand := a.DQH.Vals[a.DQH.End-1].(float64)
	// lowerband middleband

	return upperBand
}
func (a *AB) GetName() string {
	return a.Name
}
