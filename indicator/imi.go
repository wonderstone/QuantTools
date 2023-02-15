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

// IMI is the Intraday Momentum Index indicator
// https://www.investopedia.com/terms/i/intraday-momentum-index-imi.asp
type IMI struct {
	Name      string
	ParSlice  []int    //period
	InfoSlice []string // O C
	DQ        *tools.Queue
}

// NewIMI returns a new EMA indicator
func NewIMI(Name string, ParSlice []int, infoslice []string) *IMI {
	return &IMI{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        tools.New(ParSlice[0]),
	}
}

// GetName returns the name of the indicator
func (i *IMI) GetName() string {
	return i.Name
}

// LoadData loads 1 tick info datas into the indicator
func (i *IMI) LoadData(data map[string]float64) {
	i.DQ.Enqueue(data)
}

// Eval evaluates the indicator
func (i *IMI) Eval() float64 {
	var gains, losses float64
	for _, v := range i.DQ.Values() {
		bar := v.(map[string]float64)
		if bar[i.InfoSlice[1]]-bar[i.InfoSlice[0]] >= 0 {
			gains += bar[i.InfoSlice[1]] - bar[i.InfoSlice[0]]
		} else {
			losses += bar[i.InfoSlice[0]] - bar[i.InfoSlice[1]]
		}

	}
	return gains / (gains + losses) * 100
}
