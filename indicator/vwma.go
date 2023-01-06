// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Maminghui (Digital Office Product Department #2)
// revisor:
package indicator
import (
	cb "github.com/wonderstone/QuantTools/indicator/tools"
)

type VWMA struct {
	ParSlice  []int
	InfoSlice []string //[closing,volume]
	period int
	closing,volume float64
	DQ       *cb.Queue
}

func NewVWMA(ParSlice []int, infoslice []string) *VWMA {
	return &VWMA{
		ParSlice:  ParSlice,
		InfoSlice: infoslice,
		period: ParSlice[0],
	}
}

// LoadData loads 1 tick info datas into the indicator
func (v *VWMA) LoadData(data map[string]float64) {
	v.period=v.ParSlice[0]
	v.closing=data[v.InfoSlice[0]]
	v.volume=data[v.InfoSlice[1]]
}

func (v *VWMA) Eval() float64 {
	result:=(v.period+int(v.closing*v.volume))/(v.period+int(v.volume))
	return float64(result)
}