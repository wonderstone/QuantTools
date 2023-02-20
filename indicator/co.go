// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Maminghui (Digital Office Product Department #2)
// revisor:
//CO = Ema(fastPeriod, AD) - Ema(slowPeriod, AD)
package indicator

import (
	cb "github.com/wonderstone/QuantTools/indicator/tools"
)

type CO struct {
	Name                        string
	ParSlice                    []int //fastperiod,slowperiod
	fastPeriod,slowPeriod       int
	high, low, closing, volume,ad,co float64
	// info fields for indicator calculation
	InfoSlice []string //high,low,closing,volume
	DQ        *cb.Queue
	AD        *AD
	Ema1,Ema2 *EMA
}

func NewCO(Name string, ParSlice []int, infoslice []string) *CO {
	return &CO{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        cb.New(ParSlice[0]),
		AD: 	NewAD(Name,ParSlice,infoslice),//Name
		Ema1: NewEMA(Name, []int{2}, infoslice),
		Ema2: NewEMA(Name, []int{3}, infoslice),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (c *CO) LoadData(data map[string]float64) {
	c.fastPeriod=c.ParSlice[0]
	c.slowPeriod=c.ParSlice[1]
	c.high=data[c.InfoSlice[0]]
	c.low=data[c.InfoSlice[1]]
	c.closing=data[c.InfoSlice[2]]
	c.volume=data[c.InfoSlice[3]]
	c.AD.ParSlice[0]=c.fastPeriod
	c.AD.InfoSlice = []string{"high", "low", "closing","volume"}
	c.AD.LoadData(map[string]float64{"high": c.high, "low": c.low, "closing": c.closing,"volume":c.volume})
	c.ad=c.AD.Eval()
	c.Ema1.ParSlice[0] = c.ParSlice[0]
	c.Ema1.InfoSlice = []string{"values"}
	c.Ema1.LoadData(map[string]float64{"values": c.ad})
	ema1 := c.Ema1.Eval()
	c.Ema2.ParSlice[0] = c.ParSlice[1]
	c.Ema2.InfoSlice = []string{"values"}
	c.Ema2.LoadData(map[string]float64{"values": c.ad})
	ema2 := c.Ema2.Eval()
	c.co=ema1-ema2

}

// Eval evaluates the indicator
func (c *CO) Eval() float64 {
	
	return c.co
}
func (c *CO) GetName() string {
	return c.Name
}
