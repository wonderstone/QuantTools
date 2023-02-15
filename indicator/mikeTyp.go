// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Zhangweixuan (Digital Office Product Department #2)
// revisor:

package indicator

type MIKETYP struct {
	Name              string
	ParSlice          []int
	InfoSlice         []string
	high, low, _close float64
}

func NewMIKETYP(Name string, ParSlice []int, InfoSlice []string) *MIKETYP {
	return &MIKETYP{
		Name:      Name,
		ParSlice:  ParSlice, //period
		InfoSlice: InfoSlice,
	}
}

// LoadData loads 1 tick info datas into the indicator
func (t *MIKETYP) LoadData(data map[string]float64) {
	t.high = data[t.InfoSlice[0]]
	t.low = data[t.InfoSlice[1]]
	t._close = data[t.InfoSlice[2]]
}

func (t *MIKETYP) Eval() float64 {
	return (t.high + t.low + 2*t._close) / 4
}

func (t *MIKETYP) GetName() string {
	return t.Name
}
