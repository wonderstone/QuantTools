// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Zhangweixuan (Digital Office Product Department #2)
// revisor:

package indicator

type MTM struct {
	Name      string
	ParSlice  []int
	InfoSlice []string
	_close    float64
	ref       *Ref
}

func NewMTM(Name string, ParSlice []int, InfoSlice []string) *MTM {
	return &MTM{
		Name:      Name,
		ParSlice:  ParSlice, //period
		InfoSlice: InfoSlice,
		ref:       NewRef(Name, []int{ParSlice[0] + 1}, InfoSlice),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (m *MTM) LoadData(data map[string]float64) {
	m.ref.LoadData(data)
	m._close = data[m.InfoSlice[0]]
}

func (m *MTM) Eval() float64 {
	return m._close - m.ref.Eval()
}

func (m *MTM) GetName() string {
	return m.Name
}
