// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Zhangweixuan (Digital Office Product Department #2)
// revisor:

package indicator

type OBV struct {
	Name            string
	ParSlice        []int
	InfoSlice       []string
	va, bov, _close float64
}

func NewOBV(Name string, ParSlice []int, InfoSlice []string) *OBV {
	return &OBV{
		Name:      Name,
		ParSlice:  ParSlice, //period
		InfoSlice: InfoSlice,
		bov:       0,
	}
}

// LoadData loads 1 tick info datas into the indicator
func (o *OBV) LoadData(data map[string]float64) {
	if o._close < data[o.InfoSlice[0]] {
		o.va = data[o.InfoSlice[1]]
	} else if o._close > data[o.InfoSlice[0]] {
		o.va = -data[o.InfoSlice[1]]
	} else {
		o.va = 0
	}

	if o._close != 0 {
		o.bov += o.va
	}
	o._close = data[o.InfoSlice[0]]

}

func (o *OBV) Eval() float64 {
	return o.bov
}

func (o *OBV) GetName() string {
	return o.Name
}
