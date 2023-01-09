// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Maminghui (Digital Office Product Department #2)
// revisor:

package indicator

//AD calculation without period (ParSlice[0])
type AD struct {
	Name                               string
	ParSlice                           []int
	high, low, closing, volume, result float64

	// info fields for indicator calculation
	InfoSlice []string //[high,low,closing,volume]
}

func NewAD(Name string, ParSlice []int, infoslice []string) *AD {
	return &AD{
		Name:      Name,
		ParSlice:  ParSlice, //period
		InfoSlice: infoslice,
	}
}

// LoadData loads 1 tick info datas into the indicator
func (ad *AD) LoadData(data map[string]float64) {
	ad.high = data[ad.InfoSlice[0]]
	ad.low = data[ad.InfoSlice[1]]
	ad.closing = data[ad.InfoSlice[2]]
	ad.volume = data[ad.InfoSlice[3]]
}

func (ad *AD) Eval() float64 {
	ad.result += float64(ad.volume) * (((ad.closing - ad.low) - (ad.high - ad.closing)) / (ad.high - ad.low))

	return ad.result
}
func (ad *AD) GetName() string {
	return ad.Name
}
