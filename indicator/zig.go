// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Zhangweixuan (Digital Office Product Department #2)
// revisor:

package indicator

import "math"

type ZIG struct {
	Name                                     string
	ParSlice                                 []int
	InfoSlice                                []string
	inflection, extremum, rate, dist, ieDist float64
	index                                    int
	upTrend                                  bool
}

func NewZIG(Name string, ParSlice []int, InfoSlice []string) *ZIG {
	return &ZIG{
		Name:       Name,
		ParSlice:   ParSlice,
		InfoSlice:  InfoSlice,
		inflection: 0,
		rate:       float64(ParSlice[1]),
		index:      ParSlice[0],
		upTrend:    true,
		dist:       -1,
	}
}

// LoadData loads 1 tick info datas into the indicator
func (z *ZIG) LoadData(data map[string]float64) {
	z.dist++
	tmpVal := data[z.InfoSlice[z.index]]
	if z.inflection == 0 {
		z.inflection = tmpVal
		z.extremum = tmpVal
		z.ieDist = math.Inf(0)
	}

	sameTrend := false
	if (z.upTrend && tmpVal >= z.extremum) || (!z.upTrend && tmpVal <= z.extremum) {
		sameTrend = true
	}

	deltaInflect := math.Abs(tmpVal-z.inflection) / z.inflection * 100
	deltaExt5ra := math.Abs(tmpVal-z.extremum) / z.extremum * 100
	if sameTrend && deltaInflect >= z.rate {
		z.extremum = tmpVal
		z.ieDist = z.dist
	} else if !sameTrend && deltaExt5ra >= z.rate {
		z.inflection = z.extremum
		z.extremum = tmpVal
		z.upTrend = !z.upTrend
		z.dist -= z.ieDist
		z.ieDist = z.dist
	}
}

// Eval return the point on the line
func (z *ZIG) Eval() float64 {
	k := (z.extremum - z.inflection) / z.ieDist
	return k*z.dist + z.inflection
}

func (z *ZIG) GetName() string {
	return z.Name
}
