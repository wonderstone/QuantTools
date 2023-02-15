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

type BR struct {
	Name                   string
	ParSlice               []int
	InfoSlice              []string
	lv1, lv2, sumP1, sumP2 float64
	P1, P2                 *tools.Queue
}

func NewBR(Name string, ParSlice []int, InfoSlice []string) *BR {
	return &BR{
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
func (b *BR) LoadData(data map[string]float64) {
	H := data[b.InfoSlice[0]]
	L := data[b.InfoSlice[1]]
	CY := data[b.InfoSlice[2]]

	if b.P1.Full() {
		b.lv1 = b.P1.Vals[b.P1.End].(float64)
		b.lv2 = b.P2.Vals[b.P2.End].(float64)
		b.sumP1 -= b.lv1
		b.sumP2 -= b.lv2
	}

	b.P1.Enqueue(H - CY)
	b.P2.Enqueue(CY - L)

	b.sumP1 += H - CY
	b.sumP2 += CY - L
}

func (b *BR) Eval() float64 {
	return b.sumP1 / b.sumP2 * 100
}

func (b *BR) GetName() string {
	return b.Name
}
