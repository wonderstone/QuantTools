// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Yun Jinpeng (Digital Office Product Department #2)
// revisor:

package indicator

// use gods to generate the DQ
import (
	"github.com/wonderstone/QuantTools/indicator/tools"
)

// CCI 顺势指标/商品路径指标
type CCI struct {
	Name      string
	ParSlice  []int    //period
	InfoSlice []string // C H L
	DQ        *tools.Queue
}

// NewCCI  returns a new CCI indicator
func NewCCI(Name string, ParSlice []int, infoslice []string) *CCI {
	return &CCI{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        tools.New(ParSlice[0]),
	}
}

// GetName returns the name of the indicator
func (c *CCI) GetName() string {
	return c.Name
}
func (c *CCI) LoadData(data map[string]float64) {
	c.DQ.Enqueue(data)
}

// Eval evaluates the indicator
func (c *CCI) Eval() float64 {
	var typ []float64
	for _, v := range c.DQ.Values() {
		bar := v.(map[string]float64)
		typ = append(typ, (bar[c.InfoSlice[1]]+bar[c.InfoSlice[2]]+bar[c.InfoSlice[0]])/3.0)
	}
	ccl := (typ[len(typ)-1] - tools.Avg(typ)) * 1000 / (15 * tools.AvgDev(typ))
	return ccl
}
