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

type MI struct {
	Name string
	ParSlice   []int
	high,low float64
	// info fields for indicator calculation
	InfoSlice []string
	DQ *cb.Queue
	Ema1,Ema2  *EMA
	Sum *Sum
}

func NewMI(Name string,ParSlice []int, infoslice []string) *MI {
	return &MI{
		Name: Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,  //high,low
		DQ:        cb.New(ParSlice[0]),
		Ema1:        NewEMA(Name,[]int{3}, infoslice),
		Ema2:        NewEMA(Name, []int{3}, infoslice),
		Sum:         NewSum(Name,[]int{25},infoslice),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (m *MI) LoadData(data map[string]float64) {
	m.high=data[m.InfoSlice[0]]
	m.low=data[m.InfoSlice[1]]
	//m.Ema1.ptoday = data[m.Ema1.InfoSlice[0]]-data[m.Ema1.InfoSlice[1]]
	m.Ema1.InfoSlice=[]string{"Close"}
	m.Ema1.LoadData(map[string]float64{"Close": m.high-m.low})
	ema1:=m.Ema1.Eval()
	m.Ema2.InfoSlice=[]string{"Close"}
	m.Ema2.LoadData(map[string]float64{"Close": ema1})
	ema2:=m.Ema2.Eval()
	ratio:=ema1/ema2
	m.DQ.Enqueue(ratio)
	m.Sum.LoadData(map[string]float64{"Close": ratio})
	m.Sum.InfoSlice=[]string{"Close"}
}

// Eval evaluates the indicator
func (m *MI) Eval() float64 {
	
	return m.Sum.Eval()
}
func (m *MI) GetName() string {
	return m.Name
}
