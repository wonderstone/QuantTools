// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Zhangweixuan (Digital Office Product Department #2)
// revisor:

package indicator

type OSC struct {
	Name      string
	ParSlice  []int
	InfoSlice []string
	_close    float64
	ma        *MA
}

func NewOSC(Name string, ParSlice []int, InfoSlice []string) *OSC {
	return &OSC{
		Name:      Name,
		ParSlice:  ParSlice, //period
		InfoSlice: InfoSlice,
		ma:        NewMA(Name, ParSlice, InfoSlice),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (o *OSC) LoadData(data map[string]float64) {
	o.ma.LoadData(data)
	o._close = data[o.InfoSlice[0]]

}

func (o *OSC) Eval() float64 {
	return (o._close - o.ma.Eval()) * 100
}

func (o *OSC) GetName() string {
	return o.Name
}
