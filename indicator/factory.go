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

// 基础指标信息结构体
type IndiInfo struct {
	Name      string   // 指标名称
	IndiType  string   // 指标类型
	ParSlice  []int    // 指标参数切片
	InfoSlice []string // 指标信息切片
}

// 指标接口
type IIndicator interface {
	GetName() string
	LoadData(data map[string]float64)
	Eval() float64
}

// 条件加载器接口
type ConditionalLoader interface {
	Judge(ts string) bool
}

// 条件指标信息结构体
type CondIndiInfo struct {
	II IIndicator
	CL ConditionalLoader
}

// 设定简易条件加载器结构体，结构体内保存所需信息
type SimpleConditionalLoader struct {
	CLPeriodCritic string // or other useful info
}

// 实现条件加载器接口
func (scl *SimpleConditionalLoader) Judge(ts string) bool {
	// all kinds of meaningful or meaningless judgements
	// with the damned CLPeriodCritic!!!!
	return scl.CLPeriodCritic == ts
}

// 设定默认条件加载器结构体
type DefaultConditionalLoader struct {
}

// 实现默认条件加载器接口
func (dcl *DefaultConditionalLoader) Judge(ts string) bool {
	return true
}
func CIIFactory(freqlabel string, ii IndiInfo, info interface{}) CondIndiInfo {
	switch freqlabel {
	case "simple":
		return CondIndiInfo{
			II: IndiFactory(ii),
			CL: &SimpleConditionalLoader{
				CLPeriodCritic: info.(string),
			},
		}
	default:
		return CondIndiInfo{
			II: IndiFactory(ii),
			CL: &DefaultConditionalLoader{},
		}
	}
}

// get condindiinfo slice from yaml
func GetCondIndiInfoSlice(path string, infomap map[string]interface{}) []CondIndiInfo {
	var condIndiInfoSlice []CondIndiInfo
	// read yaml
	c := configer.New(path)
	err := c.Load()
	if err != nil {
		panic(err)
	}
	// get the unmarshalSlice result with type []interface{} and err
	// iis := make([]interface{}, 0)

	iis, err := c.UnmarshalSlice("ci")
	if err != nil {
		panic(err)
	}
	// iter the iis slice to get indiInfoSlice
	for _, ii := range iis {
		// convert ii to struct IndiInfo
		iiMap := ii.(map[string]interface{})
		PStmp := iiMap["parslice"].([]interface{})
		// convert []interface{} tmp to a new []int
		var parslice []int
		for _, v := range PStmp {
			parslice = append(parslice, v.(int))
		}

		IStmp := iiMap["infoslice"].([]interface{})
		// convert []interface{} tmpstring to a new []string
		var infoslice []string
		for _, v := range IStmp {
			infoslice = append(infoslice, v.(string))
		}

		indiInfo := IndiInfo{
			Name:     iiMap["name"].(string),
			IndiType: iiMap["inditype"].(string),
			// convert tmp to []int
			ParSlice:  parslice,
			InfoSlice: infoslice,
		}
		tmpfreq := iiMap["freq"].(string)
		condIndiInfoSlice = append(condIndiInfoSlice, CIIFactory(tmpfreq, indiInfo, infomap[tmpfreq]))
	}

	// unmarshal yaml to indiInfoSlice
	return condIndiInfoSlice
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

	iis, err := c.UnmarshalSlice("ii")
	if err != nil {
		panic(err)
	}
	// iter the iis slice to get indiInfoSlice
	for _, ii := range iis {
		// convert ii to struct IndiInfo
		iiMap := ii.(map[string]interface{})
		PStmp := iiMap["parslice"].([]interface{})
		// convert []interface{} tmp to a new []int
		var parslice []int
		for _, v := range PStmp {
			parslice = append(parslice, v.(int))
		}

		IStmp := iiMap["infoslice"].([]interface{})
		// convert []interface{} tmpstring to a new []string
		var infoslice []string
		for _, v := range IStmp {
			infoslice = append(infoslice, v.(string))
		}

		indiInfo := IndiInfo{
			Name:     iiMap["name"].(string),
			IndiType: iiMap["inditype"].(string),
			// convert tmp to []int
			ParSlice:  parslice,
			InfoSlice: infoslice,
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
