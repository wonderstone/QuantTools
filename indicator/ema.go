// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Wonderstone (Digital Office Product Department #2)
// revisor:

package indicator

// EMA is the moving average indicator
type EMA struct {
	ParSlice   []int
	tmp, lv, w float64
	// info fields for indicator calculation
	InfoSlice []string
}

// NewMA returns a new MA indicator
func NewEMA(ParSlice []int, infoslice []string) *EMA {
	return &EMA{
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		w:         2.0 / (float64(ParSlice[0]) + 1.0),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (e *EMA) LoadData(data map[string]float64) {
	e.tmp = data[e.InfoSlice[0]]
}

// Eval evaluates the indicator
func (e *EMA) Eval() float64 {
	res := e.w*e.tmp + (1-e.w)*e.lv
	e.lv = res
	return res
}
