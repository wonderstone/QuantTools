package datatunnel

import (
	"sync"

	"github.com/wonderstone/QuantTools/dataprocessor"
	"github.com/wonderstone/QuantTools/indicator"
)

type void struct{}

// StgTargets struct key: TargetInstrument, value: indicator set
type StgTargets struct {
	STargets map[string]map[string]void
	FTargets map[string]map[string]void
}

// StgTargetsElement struct
type StgTargetsElement struct {
	FreqType string
	StgTargets
}

type StgTargetsInfo struct {
	StgName string

	// StgTimeTriggers []string
	// StgTargetsElement slice
	STES []StgTargetsElement
	InfoET
}

type IDatatunnel interface {
	RegisterSTG(sti StgTargetsInfo)
	RemoveSTG(stgName string)
	SubProcessData()
}

// a struct for VDS sub Info FeaTure with target name and freq type
type InfoFT struct {
	TargetName string
	FreqType   string
}

// a struct with stg freq , data freq and specific Exec Time labels
type InfoET struct {
	StgFreq         string
	StgTimeTriggers []string
}

// a Struct define for dataprocessor.BarDE as promoted field with freq
type BarDEFreq struct {
	dataprocessor.BarDE
	Freq string
}

// datatunnel info struct
type InfoTunnel struct {
	// use sync.Map to store the []indicator.IndiInfo for each stg
	IISM sync.Map
	//
	StgTargetsInfo
}

// InfoTunnel method update the IISM with []indicator.IndiInfo
func (it *InfoTunnel) UpdateIISM(stgName string, IIS []indicator.IndiInfo) {
	it.IISM.Store(stgName, IIS)
}
