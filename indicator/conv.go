// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Wonderstone (Digital Office Product Department #2)
// revisor:

package indicator

// use gods to generate the queue
import (
	cb "github.com/wonderstone/QuantTools/indicator/tools"
)

// Conv is the Conv indicator
type Conv struct {
	ParSlice  []int
	InfoSlice []string
	DQS       *cb.Queue
	DQI       *cb.Queue
	Ma        *MA
}

// NewConv returns a new Conv indicator
func NewConv(ParSlice []int, infoslice []string) *Conv {
	return &Conv{
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQS:       cb.New(ParSlice[0]),
		DQI:       cb.New(ParSlice[0]),
		Ma:        NewMA(ParSlice, infoslice),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (c *Conv) LoadData(data map[string]float64) {
	c.DQS.Enqueue(data[c.InfoSlice[0]])
	c.DQI.Enqueue(data[c.InfoSlice[1]])
}

// Eval evaluates the indicator
func (c *Conv) Eval() float64 {
	var sum float64
	c.Ma.DQ = c.DQS
	avgStock := c.Ma.Eval()
	c.Ma.DQ = c.DQI
	avgIndex := c.Ma.Eval()
	for i := range c.DQS.Values() {
		sum += (c.DQI.Values()[i].(float64) - avgStock) * (c.DQI.Values()[i].(float64) - avgIndex)
	}
	return sum / float64(c.DQI.Size())
}
