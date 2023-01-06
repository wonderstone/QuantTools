// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Maminghui (Digital Office Product Department #2)
// revisor:
// Bollinger Band Width. It measures the percentage difference between the
// upper band and the lower band. It decreases as Bollinger Bands narrows
// and increases as Bollinger Bands widens
//
// During a period of rising price volatity the band width widens, and
// during a period of low market volatity band width contracts.
//
// Band Width = (Upper Band - Lower Band) / Middle Band
package indicator

import (
	cb "github.com/wonderstone/QuantTools/indicator/tools"
)

type BBW struct {
	ParSlice  []int
	InfoSlice []string //[upperBand, middleBand, lowerBand]
	upperBand, middleBand, lowerBand float64
	DQ *cb.Queue
}

func NewBBW(ParSlice []int, infoslice []string) *BBW {
	return &BBW{
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:       cb.New(ParSlice[0]),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (b *BBW) LoadData(data map[string]float64) {
	b.upperBand=data[b.InfoSlice[0]]
	b.middleBand=data[b.InfoSlice[1]]
	b.lowerBand=data[b.InfoSlice[2]]
	b.DQ.Enqueue((b.upperBand - b.lowerBand) / b.middleBand)
}

func (b *BBW) Eval() float64 {
	bandWidth:=b.DQ.End
	return float64(bandWidth)
}
