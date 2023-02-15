// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Zhangweixuan (Digital Office Product Department #2)
// revisor:

package indicator

import (
	"github.com/wonderstone/QuantTools/indicator/tools"
)

type AR struct {
	Name                   string
	ParSlice               []int
	InfoSlice              []string
	lv1, lv2, sumP1, sumP2 float64
	P1, P2                 *tools.Queue
}

func NewAR(Name string, ParSlice []int, InfoSlice []string) *AR {
	return &AR{
		Name:      Name,
		ParSlice:  ParSlice, //period
		InfoSlice: InfoSlice,
		P1:        tools.New(ParSlice[0]),
		P2:        tools.New(ParSlice[0]),
		sumP1:     0,
		sumP2:     0,
	}
}

// LoadData loads 1 tick info datas into the indicator
func (a *AR) LoadData(data map[string]float64) {
	H := data[a.InfoSlice[0]]
	L := data[a.InfoSlice[1]]
	CY := data[a.InfoSlice[2]]

	if a.P1.Full() {
		a.lv1 = a.P1.Vals[a.P1.End].(float64)
		a.lv2 = a.P2.Vals[a.P2.End].(float64)
		a.sumP1 -= a.lv1
		a.sumP2 -= a.lv2
	}

	a.P1.Enqueue(H - CY)
	a.P2.Enqueue(CY - L)

	a.sumP1 += H - CY
	a.sumP2 += CY - L
}

func (a *AR) Eval() float64 {
	return a.sumP1 / a.sumP2 * 100
}

func (a *AR) GetName() string {
	return a.Name
}
