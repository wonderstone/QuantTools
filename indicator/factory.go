// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Wonderstone (Digital Office Product Department #2)
// revisor:

package indicator

import (
	"github.com/wonderstone/QuantTools/configer"
)

type IndiInfo struct {
	Name      string
	IndiType  string
	ParSlice  []int
	InfoSlice []string
}

type IIndicator interface {
	GetName() string
	LoadData(data map[string]float64)
	Eval() float64
}

// get indiinfo slice from yaml
func GetIndiInfoSlice(dir string) []IndiInfo {
	var indiInfoSlice []IndiInfo
	// read yaml
	c := configer.New(dir + "IndicatorInfo.yaml")
	err := c.Load()
	if err != nil {
		panic(err)
	}
	// get the unmarshalSlice result with type []interface{} and err
	// iis := make([]interface{}, 0)

	iis, err := c.UnmarshalSlice()
	if err != nil {
		panic(err)
	}
	// iter the iis slice to get indiInfoSlice
	for _, ii := range iis {
		// convert ii to struct IndiInfo
		iiMap := ii.(map[string]interface{})
		tmp := iiMap["ParSlice"].([]interface{})
		// convert []interface{} tmp to a new []int
		var tmp2 []int
		for _, v := range tmp {
			tmp2 = append(tmp2, v.(int))
		}

		tmpstring := iiMap["InfoSlice"].([]interface{})
		// convert []interface{} tmpstring to a new []string
		var tmpstring2 []string
		for _, v := range tmpstring {
			tmpstring2 = append(tmpstring2, v.(string))
		}

		indiInfo := IndiInfo{
			Name:     iiMap["Name"].(string),
			IndiType: iiMap["IndiType"].(string),
			// convert tmp to []int
			ParSlice:  tmp2,
			InfoSlice: tmpstring2,
		}
		indiInfoSlice = append(indiInfoSlice, indiInfo)
	}

	// unmarshal yaml to indiInfoSlice
	return indiInfoSlice
}

// factory pattern
func IndiFactory(ii IndiInfo) IIndicator {
	switch ii.IndiType {
	case "MA":
		return NewMA(ii.Name, ii.ParSlice, ii.InfoSlice)
	case "Var":
		return NewVar(ii.Name, ii.ParSlice, ii.InfoSlice)
	case "BETA":
		return NewBeta(ii.Name, ii.ParSlice, ii.InfoSlice)
	case "Cov":
		return NewConv(ii.Name, ii.ParSlice, ii.InfoSlice)
	case "AD":
		return NewAvgDev(ii.Name, ii.ParSlice, ii.InfoSlice)
	case "Ref":
		return NewRef(ii.Name, ii.ParSlice, ii.InfoSlice)
	default:
		return nil
	}
}
