package datatunnel

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
	StgName        string
	StgFreq        string
	StgTimeTrigger []string
	// StgTargetsElement slice
	STES []StgTargetsElement
}

type IDatatunnel interface {
	RegisterSTG(sti StgTargetsInfo)
	RemoveSTG(stgName string)
	SubProcessData()
}

type InfoFT struct {
	TargetName string
	FreqType   string
}
