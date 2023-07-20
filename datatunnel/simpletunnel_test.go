package datatunnel

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// test RegisterSTG
func TestRegisterandRemoveSTG(t *testing.T) {
	//test data
	sti := StgTargetsInfo{
		StgName: "test",

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
		InfoET: InfoET{
			StgFreq: "daily",
			StgTimeTriggers: []string{
				"14:55:00",
			},
		},
	}
	sti1 := StgTargetsInfo{
		StgName: "test1",
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
		InfoET: InfoET{
			StgFreq:         "daily",
			StgTimeTriggers: []string{},
		},
	}
	sti2 := StgTargetsInfo{
		StgName: "test2",
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
		InfoET: InfoET{
			StgFreq: "daily",
			StgTimeTriggers: []string{
				"14:55:00",
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
	ip = "10.1.90.91"
	// port
	port := 9004
	port = 9904
	// channel
	ch := make(chan bool)
	// new a data tunnel
	dt := NewDataTunnel()
	// new a channel for targets change as map[string]string
	TargetsChangechan := make(chan []map[string]string)

	// new the channel
	BarDEFreqchan := make(chan *BarDEFreq)
	// test GetTargetsData
	tmpMap := map[string]string{"Msgtype": "Snapshot", "Symbol": "600000.SH"}
	tmpMaps := []map[string]string{tmpMap}
	go addNewSubDataMaps(TargetsChangechan)
	go dt.GetTargetsData(ip, port, ch, tmpMaps, TargetsChangechan, BarDEFreqchan)
	// output the channel data
	for {
		tmpBarDE := <-BarDEFreqchan
		fmt.Println(tmpBarDE)
	}

}

// func to add new subdata maps in 1 minute later
func addNewSubDataMaps(TargetsChangechan chan []map[string]string) {
	// wait 1 minute to execute following
	// new the channel
	// TargetsChangechan := make(chan []map[string]string)
	// new the subdata maps

	time.Sleep(15 * time.Second)
	fmt.Println("add new subdata maps!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	tmpMap := map[string]string{"Msgtype": "Snapshot", "Symbol": "600000.SH,510050.SH"}
	tmpMaps := []map[string]string{tmpMap}
	// add the subdata maps in 1 minute later
	// time.Sleep(1 * time.Minute)
	TargetsChangechan <- tmpMaps
	fmt.Println("new subdata maps")
}

// test float64ToByte
