// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Zhangweixuan (Digital Office Product Department #2)
// revisor:
package indicator

import "math"

type TRIX struct {
	Name             string
	ParSlice         []int
	InfoSlice        []string
	lv, v            float64
	cnt              int
	EMA1, EMA2, EMA3 *EMA
}

func NewTRIX(Name string, ParSlice []int, Infoslice []string) *TRIX {
	return &TRIX{
		Name:      Name,
		ParSlice:  ParSlice, //period
		InfoSlice: Infoslice,
		EMA1:      NewEMA(Name, ParSlice, Infoslice),
		EMA2:      NewEMA(Name, ParSlice, Infoslice),
		EMA3:      NewEMA(Name, ParSlice, Infoslice),
		v:         0,
		cnt:       0,
	}
}

// LoadData loads 1 tick info datas into the indicator
func (t *TRIX) LoadData(data map[string]float64) {
	t.EMA1.LoadData(data)
	if t.EMA1.DQ.Size() >= t.ParSlice[0] {
		t.EMA2.LoadData(map[string]float64{t.InfoSlice[0]: t.EMA1.Eval()})
	}
	if t.EMA2.DQ.Size() >= t.ParSlice[0] {
		t.EMA3.LoadData(map[string]float64{t.InfoSlice[0]: t.EMA2.Eval()})
	}
	t.lv = t.v
	t.v = t.EMA3.Eval()
	t.cnt++
}

// Eval is invalid until the 7th day
func (t *TRIX) Eval() float64 {
	if t.cnt <= t.ParSlice[0]*3-2 {
		return math.NaN()
	}
	return (t.v - t.lv) / t.v * 100
}

func (t *TRIX) GetName() string {
	return t.Name
}
