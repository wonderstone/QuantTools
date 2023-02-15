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

type MIKEWS struct {
	Name      string
	ParSlice  []int
	InfoSlice []string
	max       float64
	typ       *MIKETYP
	DQ        *cb.Queue
}

func NewMIKEWS(Name string, ParSlice []int, InfoSlice []string) *MIKEWS {
	return &MIKEWS{
		Name:      Name,
		ParSlice:  ParSlice, //period
		InfoSlice: InfoSlice,
		typ:       NewMIKETYP(Name, []int{}, InfoSlice),
		DQ:        cb.New(ParSlice[0]),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (w *MIKEWS) LoadData(data map[string]float64) {
	w.typ.LoadData(data)
	w.DQ.Enqueue(data[w.InfoSlice[0]])

	w.max = 0
	for i := 0; i < w.DQ.Size(); i++ {
		w.max = math.Max(w.max, w.DQ.Vals[i].(float64))
	}
}

func (w *MIKEWS) Eval() float64 {
	return 2*w.typ.Eval() - w.max
}

func (w *MIKEWS) GetName() string {
	return w.Name
}
