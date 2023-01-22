// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Maminghui (Digital Office Product Department #2)
// revisor:
package indicator

import (
	cb "github.com/wonderstone/QuantTools/indicator/tools"
)

type PVO struct {
	Name        string
	ParSlice    []int //fastPeriod, slowPeriod, signalPeriod
	volume, pvo float64
	// info fields for indicator calculation
	InfoSlice        []string //volume
	DQ               *cb.Queue
	fastEma, slowEma *EMA
}

func NewPVO(Name string, ParSlice []int, infoslice []string) *PVO {
	return &PVO{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice, //high,low
		DQ:        cb.New(ParSlice[0]),
		fastEma:   NewEMA(Name, []int{2}, infoslice),
		slowEma:   NewEMA(Name, []int{3}, infoslice),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (p *PVO) LoadData(data map[string]float64) {
	p.volume = data[p.InfoSlice[0]]
	p.fastEma.InfoSlice = []string{"Close"}
	p.fastEma.LoadData(map[string]float64{"Close": p.volume})
	fastEma := p.fastEma.Eval()
	p.slowEma.InfoSlice = []string{"Close"}
	p.slowEma.LoadData(map[string]float64{"Close": p.volume})
	slowEma := p.slowEma.Eval()
	p.pvo = ((fastEma - slowEma) / slowEma) * 100
}

// Eval evaluates the indicator
func (p *PVO) Eval() float64 {

	return p.pvo
}
func (p *PVO) GetName() string {
	return p.Name
}
