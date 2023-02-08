package framework

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wonderstone/QuantTools/account/virtualaccount"
	"github.com/wonderstone/QuantTools/realinfo"
)

// test NewBackTestConfig
func TestNewBackTestConfig(t *testing.T) {
	bt := NewBackTestConfig("../config/Manual/", "BackTest.yaml", "Default")
	assert.Equal(t, "2017/10/9 9:39", bt.BeginDate)

}

func TestNewRealTimeConfig(t *testing.T) {
	va := virtualaccount.NewVirtualAccountFromConfig("../config/Manual", "")
	info := realinfo.NewInfoFromConfig("../config/Manual", "accountinfo")
	rt := NewRealTimeConfig("../config/Manual", "realtime", info.IM, &va)
	assert.Equal(t, "2017/10/9 9:39", rt.VA.SAcct.InitTime)

}

// test PrepareData
func TestPrepareData(t *testing.T) {
	bt := NewBackTestConfig("../config/Manual/", "BackTest.yaml", "Default")
	bt.PrepareData("VDS")
	assert.Equal(t, "2017/10/9 9:39", bt.BCM.BeginDate)

}
