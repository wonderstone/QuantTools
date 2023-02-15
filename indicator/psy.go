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

type PSY struct {
	Name      string
	ParSlice  []int
	InfoSlice []string
	rise, lv  float64
	DQ        *cb.Queue
}

func NewPSY(Name string, ParSlice []int, InfoSlice []string) *PSY {
	return &PSY{
		Name:      Name,
		ParSlice:  ParSlice, //period
		InfoSlice: InfoSlice,
		rise:      0,
		DQ:        cb.New(ParSlice[0]),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (p *PSY) LoadData(data map[string]float64) {
	tmpVal := data[p.InfoSlice[0]]
	if p.DQ.Full() {
		p.lv = p.DQ.Vals[p.DQ.End].(float64)
	}

	tmpPeek, acc := p.DQ.Peek()
	p.DQ.Enqueue(tmpVal)
	if acc && tmpVal > tmpPeek.(float64) {
		p.rise++
	}

	if p.DQ.Full() {
		if p.lv < p.DQ.Vals[p.DQ.End].(float64) {
			p.rise--
		}
	}
}

func (p *PSY) Eval() float64 {
	if p.DQ.Size() < 2 {
		return math.NaN()
	}
	return p.rise / float64(p.DQ.Size()) * 100
}

func (p *PSY) GetName() string {
	return p.Name
}
