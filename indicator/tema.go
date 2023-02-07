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

type Tema struct {
	Name        string
	ParSlice    []int //period
	period,values,tema float64
	// info fields for indicator calculation
	InfoSlice        []string //values
	DQ               *cb.Queue
	EMA1,EMA2,EMA3        *EMA
}

func NewTema(Name string, ParSlice []int, infoslice []string) *Tema {
	return &Tema{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice, 
		DQ:        cb.New(ParSlice[0]),
		EMA1:   NewEMA(Name, ParSlice, infoslice),
		EMA2:   NewEMA(Name, ParSlice, infoslice),
		EMA3:   NewEMA(Name, ParSlice, infoslice),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (t *Tema) LoadData(data map[string]float64) {
	t.period=float64(t.ParSlice[0])
	t.values=data[t.InfoSlice[0]]

	t.EMA1.InfoSlice = []string{"values"}
	t.EMA1.ParSlice[0]=int(t.period)
	t.EMA1.LoadData(map[string]float64{"values": t.values})
	ema1 := t.EMA1.Eval()
	t.EMA2.InfoSlice = []string{"values"}
	t.EMA2.ParSlice[0]=int(t.period)
	t.EMA2.LoadData(map[string]float64{"values": ema1})
	ema2 := t.EMA2.Eval()
	t.EMA3.InfoSlice = []string{"values"}
	t.EMA3.ParSlice[0]=int(t.period)
	t.EMA3.LoadData(map[string]float64{"values": ema2})
	ema3 := t.EMA3.Eval()
	t.tema = (ema1*3-ema2*3)+ema3
}

// Eval evaluates the indicator
func (t *Tema) Eval() float64 {

	return t.tema
}
func (t *Tema) GetName() string {
	return t.Name
}
