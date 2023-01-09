// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Maminghui (Digital Office Product Department #2)
package indicator

type BOP struct {
	Name                        string
	closing, opening, high, low float64

	// info fields for indicator calculation
	InfoSlice []string //[closing,opening,high,low]
}

func NewBOP(Name string, infoslice []string) *BOP {
	return &BOP{
		Name:      Name,
		InfoSlice: infoslice,
	}
}
func (bop *BOP) LoadData(data map[string]float64) {
	bop.closing = data[bop.InfoSlice[0]]
	bop.opening = data[bop.InfoSlice[1]]
	bop.high = data[bop.InfoSlice[2]]
	bop.low = data[bop.InfoSlice[3]]
}
func (bop *BOP) Eval() float64 {
	//result := (bop.closing-bop.opening)/(bop.high- bop.low)
	return (bop.closing - bop.opening) / (bop.high - bop.low)
}
func (bop *BOP) GetName() string {
	return bop.Name
}
