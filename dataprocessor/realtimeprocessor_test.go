package dataprocessor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wonderstone/QuantTools/indicator"
)

// test realtimeprocessor
func TestFakeGetHistoryData(t *testing.T) {
	dir := "../tmpdata/stockdata/test/"
	fmt.Println("dir:", dir)

	d, s := FakeGetHistoryData(dir, "VDS")
	// fmt.Println("d:", d)
	// fmt.Println("s:", s)
	for _, v := range s[:30] {
		fmt.Println(v, d[v])

	}
	assert.Equal(t, len(s), 164, "should be 164")

}

// get a map with instID as key and slice of indicator.IndiInfo as value
func GetIndiInfoMap(IDs []string, iis []indicator.IndiInfo) (iim map[string][]indicator.IIndicator) {
	// init the map
	iim = make(map[string][]indicator.IIndicator)
	// iter the IDs
	for _, ID := range IDs {
		for _, ii := range iis {
			iim[ID] = append(iim[ID], indicator.IndiFactory(ii))
		}
	}
	return iim
}

// test GetIndiInfoMap
func TestGetIndiInfoMap(t *testing.T) {
	IDs := []string{"sh600409", "sz000058"}
	iis := []indicator.IndiInfo{
		{Name: "MA10", IndiType: "MA", ParSlice: []int{3}, InfoSlice: []string{"close"}},
		{Name: "Var10", IndiType: "Var", ParSlice: []int{3}, InfoSlice: []string{"close"}},
	}

	iim := GetIndiInfoMap(IDs, iis)
	fmt.Println("iim:", iim["sh600409"][0].GetName())
	assert.Equal(t, len(iim["sh600409"]), 2, "should be 2")

}

// test AddIndicatorsToSData
func TestAddIndicatorsToSData(t *testing.T) {
	// get the history data
	dir := "../tmpdata/stockdata/test/"
	fmt.Println("dir:", dir)
	dataMap, timeStampSlice := FakeGetHistoryData(dir, "VDS")
	// get the indiInfoMap
	IDs := []string{"sh510050"}
	iis := []indicator.IndiInfo{
		{Name: "MA3", IndiType: "MA", ParSlice: []int{3}, InfoSlice: []string{"Close"}},
		{Name: "Var3", IndiType: "Var", ParSlice: []int{3}, InfoSlice: []string{"Close"}},
	}
	// new a temp map for 100 map[string]*BarC data
	tmpMap := make(map[string]*BarC)
	// iter former 100 timeStampSlice data and add to the tmpMap
	for _, timeStamp := range timeStampSlice[:100] {
		tmpMap[timeStamp] = dataMap[timeStamp]
	}
	oldLen := len(tmpMap["2023.01.18T13:32:00.000"].Stockdata["sh510050"].IndiDataMap)
	iim := GetIndiInfoMap(IDs, iis)
	// iter the timeStampSlice
	for _, timeStamp := range timeStampSlice[:100] {
		fmt.Println("timeStamp:", timeStamp)

		// iter the IDs
		for _, ID := range IDs {
			// iter the indicators
			AddIndicatorsToSData(tmpMap[timeStamp], ID, iim[ID])
		}
		fmt.Println("tmpMap[timeStamp]:", tmpMap[timeStamp])
	}
	//
	NewLen := len(tmpMap["2023.01.18T13:32:00.000"].Stockdata["sh510050"].IndiDataMap)

	assert.Equal(t, NewLen, oldLen+2, "difference should be 2")

}
