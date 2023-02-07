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

type PPO struct {
	Name       string
	ParSlice   []int //fastPeriod, slowPeriod, signalPeriod
	price, ppo float64
	// info fields for indicator calculation
	InfoSlice        []string //price
	DQ               *cb.Queue
	fastEma, slowEma *EMA
}

func NewPPO(Name string, ParSlice []int, infoslice []string) *PPO {
	return &PPO{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice, //high,low
		DQ:        cb.New(ParSlice[0]),
		fastEma:   NewEMA(Name, []int{2}, infoslice),
		slowEma:   NewEMA(Name, []int{3}, infoslice),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (p *PPO) LoadData(data map[string]float64) {
	p.price = data[p.InfoSlice[0]]
	p.fastEma.InfoSlice = []string{"Close"}
	p.fastEma.LoadData(map[string]float64{"Close": p.price})
	fastEma := p.fastEma.Eval()
	p.slowEma.InfoSlice = []string{"Close"}
	p.slowEma.LoadData(map[string]float64{"Close": p.price})
	slowEma := p.slowEma.Eval()
	p.ppo = ((fastEma - slowEma) / slowEma) * 100
}

// Eval evaluates the indicator
func (p *PPO) Eval() float64 {

	return p.ppo
}
func (p *PPO) GetName() string {
	return p.Name
}
