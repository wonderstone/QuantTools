// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Zhangweixuan (Digital Office Product Department #2)
// revisor:
package indicator

type MACD struct {
	Name      string
	ParSlice  []int
	InfoSlice []string
	dif       *DIF
	dea       *DEA
}

func NewMACD(Name string, ParSlice []int, infoslice []string) *MACD {
	return &MACD{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		dif:       NewDIF(Name, ParSlice, infoslice),
		dea:       NewDEA(Name, ParSlice, infoslice),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (m *MACD) LoadData(data map[string]float64) {
	m.dif.LoadData(data)
	m.dea.LoadData(data)
}

// Eval evaluates the indicator
func (m *MACD) Eval() float64 {
	return 2 * (m.dif.Eval() - m.dea.Eval())
}

func (m *MACD) GetName() string {
	return m.Name
}
