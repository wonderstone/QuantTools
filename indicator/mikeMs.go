// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Zhangweixuan (Digital Office Product Department #2)
// revisor:

package indicator

import (
	cb "github.com/wonderstone/QuantTools/indicator/tools"

	"math"
)

type MIKEMS struct {
	Name         string
	ParSlice     []int
	InfoSlice    []string
	max, min     float64
	typ          *MIKETYP
	maxDQ, minDQ *cb.Queue
}

func NewMIKEMS(Name string, ParSlice []int, InfoSlice []string) *MIKEMS {
	return &MIKEMS{
		Name:      Name,
		ParSlice:  ParSlice, //period
		InfoSlice: InfoSlice,
		typ:       NewMIKETYP(Name, []int{}, InfoSlice),
		maxDQ:     cb.New(ParSlice[0]),
		minDQ:     cb.New(ParSlice[0]),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (m *MIKEMS) LoadData(data map[string]float64) {
	m.typ.LoadData(data)
	m.maxDQ.Enqueue(data[m.InfoSlice[0]])
	m.minDQ.Enqueue(data[m.InfoSlice[1]])

	m.max = 0
	for i := 0; i < m.maxDQ.Size(); i++ {
		m.max = math.Max(m.max, m.maxDQ.Vals[i].(float64))
	}

	m.min = math.Inf(0)
	for i := 0; i < m.minDQ.Size(); i++ {
		m.min = math.Min(m.min, m.minDQ.Vals[i].(float64))
	}
}

func (m *MIKEMS) Eval() float64 {
	return m.typ.Eval() - m.max + m.min
}

func (m *MIKEMS) GetName() string {
	return m.Name
}
