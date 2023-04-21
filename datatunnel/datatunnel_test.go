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
		STES: []StgTargetsElement{
			{
				Frequency: "1m",
				StgTargets: StgTargets{
					STargets: map[string]void{
						"600000.sh": {},
						"600001.sh": {},
					},
					FTargets: map[string]void{
						"cu2303": {},
					},
				},
			},
		},
	}
	sti1 := StgTargetsInfo{
		StgName: "test1",
		STES: []StgTargetsElement{
			{
				Frequency: "1m",
				StgTargets: StgTargets{
					STargets: map[string]void{
						"600000.sh": {},
						"600002.sh": {},
					},
					FTargets: map[string]void{
						"cu2303": {},
					},
				},
			},
			{
				Frequency: "1d",
				StgTargets: StgTargets{
					STargets: map[string]void{
						"600000.sh": {},
						"600001.sh": {},
						"600005.sh": {},
					},
					FTargets: map[string]void{
						"cu2303": {},
					},
				},
			},
		},
	}
	sti2 := StgTargetsInfo{
		StgName: "test2",
		STES: []StgTargetsElement{
			{
				Frequency: "1d",
				StgTargets: StgTargets{
					STargets: map[string]void{
						"600000.sh": {},
						"600001.sh": {},
						"600002.sh": {},
					},
					FTargets: map[string]void{
						"cu2303": {},
					},
				},
			},
		},
	}

	dt := NewDataTunnel()
	//test
	dt.RegisterSTG(sti)
	//assert
	assert.Equal(t, dt.StgM["test"].StgName, "test")
	dt.RegisterSTG(sti1)
	assert.Equal(t, dt.StgM["test1"].StgName, "test1")
	dt.RegisterSTG(sti2)
	assert.Equal(t, dt.StgM["test2"].StgName, "test2")

	// test RemoveSTG
	dt.RemoveSTG(sti1)
	assert.Equal(t, len(dt.StgM), 2)

}
