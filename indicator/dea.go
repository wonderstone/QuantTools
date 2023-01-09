// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Zhangweixuan (Digital Office Product Department #2)
// revisor:
package indicator

type DEA struct {
	Name       string
	ParSlice   []int
	InfoSlice  []string
	ptoday, lv float64
	dif        *DIF
}

func NewDEA(Name string, ParSlice []int, infoslice []string) *DEA {
	return &DEA{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		dif:       NewDIF(Name, ParSlice, infoslice),
		lv:        0,
		ptoday:    0,
	}
}

// LoadData loads 1 tick info datas into the indicator
func (d *DEA) LoadData(data map[string]float64) {
	d.lv = d.ptoday
	d.dif.LoadData(data)
	d.ptoday = d.Eval()
}

// Eval evaluates the indicator
func (d *DEA) Eval() float64 {
	return (d.dif.Eval()*2 + d.lv*8) / 10
}

func (d *DEA) GetName() string {
	return d.Name
}
