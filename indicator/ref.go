// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Wonderstone (Digital Office Product Department #2)
// revisor:

package indicator

import (
	"math"

	cb "github.com/wonderstone/QuantTools/indicator/tools"
)

// EMA is the moving average indicator
type Ref struct {
	Name     string
	ParSlice []int
	// info fields for indicator calculation
	InfoSlice []string
	DQ        *cb.Queue
}

// NewMA returns a new MA indicator
func NewRef(Name string, ParSlice []int, infoslice []string) *Ref {
	return &Ref{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		DQ:        cb.New(ParSlice[0]),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (r *Ref) LoadData(data map[string]float64) {
	r.DQ.Enqueue(data[r.InfoSlice[0]])
}

// Eval evaluates the indicator
func (r *Ref) Eval() float64 {
	res := r.DQ.Vals[r.DQ.End]
	// if res is nil return float64.nan
	if res == nil {
		return math.NaN()
	}
	return res.(float64)
}

// GetName returns the name of the indicator
func (r *Ref) GetName() string {
	return r.Name
}
