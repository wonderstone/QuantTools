package datatunnel

import "github.com/wonderstone/QuantTools/dataprocessor"

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

// a Struct define for dataprocessor.BarDE as promoted field with freq
type BarDEFreq struct {
	dataprocessor.BarDE
	Freq string
}

// a struct with stg freq , data freq and specific Exec Time labels
type InfoET struct {
	StgFreq         string
	StgTimeTriggers []string
}
