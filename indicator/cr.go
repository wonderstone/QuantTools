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

// CR is the moving average indicator
type CR struct {
	Name                       string
	ParSlice                   []int
	InfoSlice                  []string
	lv1, lv2, ym, sumP1, sumP2 float64
	cnt                        int
	P1, P2                     *cb.Queue
}

func NewCR(Name string, ParSlice []int, InfoSlice []string) *CR {
	return &CR{
		Name:      Name,
		ParSlice:  ParSlice, //period
		InfoSlice: InfoSlice,
		P1:        cb.New(ParSlice[0]),
		P2:        cb.New(ParSlice[0]),
		ym:        0,
		sumP1:     0,
		sumP2:     0,
		cnt:       0,
	}
}

// LoadData loads 1 tick info datas into the indicator
func (c *CR) LoadData(data map[string]float64) {
	high := data[c.InfoSlice[0]]
	low := data[c.InfoSlice[1]]
	_close := data[c.InfoSlice[2]]

	if c.P1.Full() {
		c.sumP1 -= c.P1.Vals[c.P1.End].(float64)
		c.sumP2 -= c.P2.Vals[c.P2.End].(float64)
	}

	if c.cnt != 0 {
		c.sumP1 += math.Max(0, high-c.ym)
		c.sumP2 += math.Max(0, c.ym-low)
		c.P1.Enqueue(math.Max(0, high-c.ym))
		c.P2.Enqueue(math.Max(0, c.ym-low))
	}

	c.ym = (high + low + _close) / 3

	c.cnt++
}

func (c *CR) Eval() float64 {
	if c.cnt <= c.ParSlice[0] {
		return math.NaN()
	}
	return c.sumP1 / c.sumP2 * 100
}

func (c *CR) GetName() string {
	return c.Name
}
