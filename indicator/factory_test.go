// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Wonderstone (Digital Office Product Department #2)
// revisor:

package indicator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test the indicator factory
func TestFactory(t *testing.T) {
	// Test the indicator factory
	indis := []IndiInfo{
		{"MA10", "MA", []int{3}, []string{"Close"}},
		{"Var10", "Var", []int{3}, []string{"Close"}},
	}

	for _, indi := range indis {
		indicator := IndiFactory(indi)
		fmt.Println(indicator)
		assert.NotNil(t, indicator, "Indicator should not be nil")
	}

}
func TestGetIndiInfoSlice(t *testing.T) {
	indis := GetIndiInfoSlice("../config/Manual/")
	fmt.Println(indis)
	assert.NotNil(t, indis, "Indicator should not be nil")
}

// test GetCondIndiInfoSlice

func TestGetCondIndiInfoSlice(t *testing.T) {
	tmpinfoMap := make(map[string]interface{})
	tmpinfoMap["Simple"] = "1450000000"
	tmpinfoMap["Default"] = ""
	condIndis := GetCondIndiInfoSlice("../config/Manual/CondIndicatorInfo.yaml", tmpinfoMap)
	fmt.Println(condIndis)
	assert.NotNil(t, condIndis, "Indicator should not be nil")
}

type void struct{}

type OptInfo struct {
	SubInfoS  []string
	IndiInfoS []string
}

// a method for OptInfo to optimize the IndiInfoS and SubInfoS
func (optInfo *OptInfo) OptimizeIndiInfoSlice(indiInfo string, IndiInfoS []IndiInfo, SubInfoS []string) {
	IndiInfoM := make(map[string][]string)
	// convert IndiInfoS to map
	for _, v := range IndiInfoS {
		IndiInfoM[v.Name] = v.InfoSlice
	}
	SubInfoM := make(map[string]void)
	// convert SubInfoS to map
	for _, v := range SubInfoS {
		SubInfoM[v] = void{}
	}
	newSubInfoM := make(map[string]void)
	for _, v := range optInfo.SubInfoS {
		newSubInfoM[v] = void{}
	}
	newIndiInfoM := make(map[string]void)
	for _, v := range optInfo.IndiInfoS {
		newIndiInfoM[v] = void{}
	}
	// check if indiInfo is in SubInfoS
	if _, ok := SubInfoM[indiInfo]; ok {
		// 属于订阅信息
		// check if indiInfo is in newSubInfoM
		if _, ok := newSubInfoM[indiInfo]; ok {
			// 已经添加
			// do nothing
		} else {
			optInfo.SubInfoS = append(optInfo.SubInfoS, indiInfo)
			newSubInfoM[indiInfo] = void{}
		}

	} else {
		// 属于指标信息
		// check if indiInfo is in newIndiInfoM
		if infos, ok := IndiInfoM[indiInfo]; ok {
			// check if indiInfo is in newIndiInfoM
			if _, ok := newIndiInfoM[indiInfo]; ok {
				// 已经添加
				// do nothing
			} else {
				optInfo.IndiInfoS = append([]string{indiInfo}, optInfo.IndiInfoS...)
				newIndiInfoM[indiInfo] = void{}
				for _, v := range infos {
					optInfo.OptimizeIndiInfoSlice(v, IndiInfoS, SubInfoS)
				}
			}
		} else {
			panic("IndiInfoM should have the key: " + indiInfo)
		}

	}
}

// test the OptimizeIndiInfoSlice
func TestOptimizeIndiInfoSlice(t *testing.T) {
	optInfo := OptInfo{
		SubInfoS:  []string{},
		IndiInfoS: []string{},
	}
	IndiInfoS := []IndiInfo{
		{"MA10", "MA", []int{3}, []string{"Close"}},
		{"Var10", "Var", []int{3}, []string{"MA10"}},
		{"Var12", "Var", []int{3}, []string{"Var10"}},
	}
	SubInfoS := []string{"Close", "Open", "High", "Low"}
	optInfo.OptimizeIndiInfoSlice("Var12", IndiInfoS, SubInfoS)
	fmt.Println(optInfo)
	assert.NotNil(t, optInfo, "optInfo should not be nil")
}
