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

// CMF Chaikin Money Flow
// https://www.fidelity.com/learning-center/trading-investing/technical-analysis/technical-indicator-guide/cmf
type CMF struct {
	Name      string
	ParSlice  []int    //period
	InfoSlice []string // C H L V
	DQ        *tools.Queue
}

// NewCMF  returns a new CMF indicator
func NewCMF(Name string, ParSlice []int, infoslice []string) *CMF {
	return &CMF{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        tools.New(ParSlice[0]),
	}
}

// GetName returns the name of the indicator
func (c *CMF) GetName() string {
	return c.Name
}

// LoadData loads 1 tick info datas into the indicator
func (c *CMF) LoadData(data map[string]float64) {
	c.DQ.Enqueue(data)
}

// Eval evaluates the indicator
func (c *CMF) Eval() float64 {
	var cmf []float64
	var vol []float64
	for _, v := range c.DQ.Values() {
		bar := v.(map[string]float64)
		cmf = append(cmf, (((bar[c.InfoSlice[0]]-bar[c.InfoSlice[2]])-(bar[c.InfoSlice[1]]-bar[c.InfoSlice[0]]))/(bar[c.InfoSlice[1]]-bar[c.InfoSlice[2]]))*bar[c.InfoSlice[2]])
		vol = append(vol, bar[c.InfoSlice[3]])
	}

	return tools.Sum1(cmf) / tools.Sum1(vol)
}
