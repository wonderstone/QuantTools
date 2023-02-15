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

// MedianPrice is the Median Price indicator
type MedianPrice struct {
	Name      string
	ParSlice  []int    //period
	InfoSlice []string //   H L
	DQ        *tools.Queue
}

// NewMedianPrice returns a new MedianPrice indicator
func NewMedianPrice(Name string, ParSlice []int, infoslice []string) *MedianPrice {
	return &MedianPrice{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        tools.New(1),
	}
}

// GetName returns the name of the indicator
func (m *MedianPrice) GetName() string {
	return m.Name
}

// LoadData loads 1 tick info datas into the indicator
func (m *MedianPrice) LoadData(data map[string]float64) {
	m.DQ.Enqueue(data)
}

// Eval evaluates the indicator
func (m *MedianPrice) Eval() float64 {
	bar := m.DQ.Values()[0].(map[string]float64)
	return (bar[m.InfoSlice[0]] + bar[m.InfoSlice[1]]) / 2.0
}
