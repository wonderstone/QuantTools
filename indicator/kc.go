// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Maminghui (Digital Office Product Department #2)
// revisor:
// Middle Line = EMA(period, closings)
// Upper Band = EMA(period, closings) + 2 * ATR(period, highs, lows, closings)
// Lower Band = EMA(period, closings) - 2 * ATR(period, highs, lows, closings)
package indicator

import (
	cb "github.com/wonderstone/QuantTools/indicator/tools"
)

type KC struct {
	Name                             string
	ParSlice                         []int //period
	period                           int
	high, low, closing, atr2         float64
	middleLine, upperBand, lowerBand float64
	// info fields for indicator calculation
	InfoSlice []string //high,low,closing
	DQ        *cb.Queue
	atr       *Atr
	Ema       *EMA
}

func NewKC(Name string, ParSlice []int, infoslice []string) *KC {
	return &KC{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        cb.New(ParSlice[0]),
		atr:       NewAtr(Name, ParSlice, infoslice),
		Ema:       NewEMA(Name, ParSlice, infoslice),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (k *KC) LoadData(data map[string]float64) {
	k.period = k.ParSlice[0]
	k.high = data[k.InfoSlice[0]]
	k.low = data[k.InfoSlice[1]]
	k.closing = data[k.InfoSlice[2]]
	k.atr.ParSlice[0] = k.period
	k.atr.InfoSlice = []string{"high", "low", "closing"}
	k.atr.LoadData(map[string]float64{"high": k.high, "low": k.low, "closing": k.closing})
	k.atr2 = k.atr.Eval() * 2
	k.Ema.ParSlice[0] = k.ParSlice[0]
	k.Ema.InfoSlice = []string{"closing"}
	k.Ema.LoadData(map[string]float64{"closing": k.atr.closing})
	k.middleLine = k.Ema.Eval()
	k.upperBand = k.middleLine + k.atr2
	k.lowerBand = k.middleLine - k.atr2

}

// Eval evaluates the indicator
func (k *KC) Eval() (float64, float64, float64) {

	return k.middleLine, k.upperBand, k.lowerBand
}
func (k *KC) GetName() string {
	return k.Name
}
