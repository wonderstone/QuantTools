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

// BollingUpper 真实波幅  同花顺
type BollingUpper struct {
	Name      string
	ParSlice  []int    //period P
	InfoSlice []string // C
	DQ        *tools.Queue
}

// NewBollingUpper  returns a new BollingUpper indicator
func NewBollingUpper(Name string, ParSlice []int, infoslice []string) *BollingUpper {
	return &BollingUpper{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        tools.New(ParSlice[0]),
	}
}

// GetName returns the name of the indicator
func (m *BollingUpper) GetName() string {
	return m.Name
}

// LoadData loads 1 tick info datas into the indicator
func (m *BollingUpper) LoadData(data map[string]float64) {
	m.DQ.Enqueue(data)
}

// Eval evaluates the indicator
func (m *BollingUpper) Eval() float64 {
	var close []float64
	for _, v := range m.DQ.Values() {
		bar := v.(map[string]float64)
		close = append(close, bar["Close"])
	}
	sum := tools.Sum1(close)
	mid := sum / float64(m.ParSlice[0])
	std := tools.Std(close)
	upper := mid + float64(m.ParSlice[1])*std
	return upper
}
