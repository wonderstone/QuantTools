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

// Beta is the Beta indicator
type Beta struct {
	ParSlice []int
	// info fields for indicator calculation
	InfoSlice []string
	DQS       *cb.Queue
	DQI       *cb.Queue
	Conv      *Conv
	Var       *Var
}

// NewBeta returns a new BetaCoefficient indicator
func NewBeta(ParSlice []int, infoslice []string) *Beta {
	return &Beta{
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQS:       cb.New(ParSlice[0]),
		DQI:       cb.New(ParSlice[0]),
		Conv:      NewConv(ParSlice, infoslice),
		Var:       NewVar(ParSlice, infoslice),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (b *Beta) LoadData(data map[string]float64) {
	b.DQS.Enqueue(data[b.InfoSlice[0]])
	b.DQI.Enqueue(data[b.InfoSlice[1]])

}

// Eval evaluates the indicator
func (b *Beta) Eval() float64 {
	b.Conv.DQS = b.DQS
	b.Conv.DQI = b.DQI
	b.Var.DQ = b.DQI
	return b.Conv.Eval() / b.Var.Eval()
}
