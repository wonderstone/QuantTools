package datatunnel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test RegisterSTG
func TestRegisterandRemoveSTG(t *testing.T) {
	//test data
	sti := StgTargetsInfo{
		StgName: "test",
		StgFreq: "daily",
		StgTimeTrigger: []string{
			"14:55:00",
		},
		STES: []StgTargetsElement{
			{
				FreqType: "SnapShot",
				StgTargets: StgTargets{
					STargets: map[string]map[string]void{
						"600000.sh": {"high": void{}, "low": void{}},
						"600001.sh": {"high": void{}, "low": void{}},
					},
					FTargets: map[string]map[string]void{
						"cu2303": {"high": void{}, "low": void{}},
					},
				},
			},
		},
	}
	sti1 := StgTargetsInfo{
		StgName: "test1",
		StgFreq: "daily",
		STES: []StgTargetsElement{
			{
				FreqType: "SnapShot",
				StgTargets: StgTargets{
					STargets: map[string]map[string]void{
						"600000.sh": {"high": void{}, "low": void{}, "open": void{}},
						"600002.sh": {"high": void{}, "low": void{}},
					},
					FTargets: map[string]map[string]void{
						"cu2303": {"high": void{}, "low": void{}},
					},
				},
			},
			{
				FreqType: "1d",
				StgTargets: StgTargets{
					STargets: map[string]map[string]void{
						"600000.sh": {"high": void{}, "low": void{}},
						"600001.sh": {"high": void{}, "low": void{}},
						"600005.sh": {"high": void{}, "low": void{}},
					},
					FTargets: map[string]map[string]void{
						"cu2303": {"high": void{}, "low": void{}},
					},
				},
			},
		},
	}
	sti2 := StgTargetsInfo{
		StgName: "test2",
		StgFreq: "daily",
		StgTimeTrigger: []string{
			"14:55:00",
		},
		STES: []StgTargetsElement{
			{
				FreqType: "1d",
				StgTargets: StgTargets{
					STargets: map[string]map[string]void{
						"600000.sh": {"high": void{}, "low": void{}},
						"600001.sh": {"high": void{}, "low": void{}},
						"600002.sh": {"high": void{}, "low": void{}},
					},
					FTargets: map[string]map[string]void{
						"cu2303": {"high": void{}, "low": void{}},
					},
				},
			},
		},
	}

	dt := NewDataTunnel()
	//test
	dt.RegisterSTG(sti)
	//assert
	// assert.Equal(t, dt.InfoM["test"].StgName, "test")
	dt.RegisterSTG(sti1)
	// assert.Equal(t, dt.InfoM["test1"].StgName, "test1")
	dt.RegisterSTG(sti2)
	// assert.Equal(t, dt.InfoM["test2"].StgName, "test2")

	// test RemoveSTG
	dt.RemoveSTG(sti1)
	assert.Equal(t, 7, len(dt.InfoM), "dt.InfoM should have 7 elements")

}

// test GetTargetsData

func TestGetTargetsData(t *testing.T) {
	// ip
	ip := "123.138.216.197"
	// port
	port := 9004
	// channel
	ch := make(chan bool)
	// new a data tunnel
	dt := NewDataTunnel()

	// test GetTargetsData
	dt.GetTargetsData(ip, port, ch)

}

// test float64ToByte
