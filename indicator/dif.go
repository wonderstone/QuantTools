// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Zhangweixuan (Digital Office Product Department #2)
// revisor:
package indicator

// DIF returns differential value
type DIF struct {
	Name      string
	ParSlice  []int
	InfoSlice []string

	EMA12, EMA26 *EMA
}

func NewDIF(Name string, ParSlice []int, infoslice []string) *DIF {
	return &DIF{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		EMA12:     NewEMA(Name, []int{12}, infoslice),
		EMA26:     NewEMA(Name, []int{26}, infoslice),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (d *DIF) LoadData(data map[string]float64) {
	d.EMA12.LoadData(data)
	d.EMA26.LoadData(data)
}

// Eval evaluates the indicator
func (d *DIF) Eval() float64 {
	return d.EMA12.Eval() - d.EMA26.Eval()
}

func (d *DIF) GetName() string {
	return d.Name
}
