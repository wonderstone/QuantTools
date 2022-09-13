package framework

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test NewBackTestConfig
func TestNewBackTestConfig(t *testing.T) {
	bt := NewBackTestConfig("Default", "../config/Manual")
	assert.Equal(t, "2017/10/9 9:39", bt.BeginDate)

}

func TestNewRealTimeConfig(t *testing.T) {
	rt := NewRealTimeConfig("../config/Manual")
	assert.Equal(t, "2017/10/9 9:39", rt.VA.SAcct.InitTime)

}

// test PrepareData
func TestPrepareData(t *testing.T) {
	bt := NewBackTestConfig("Default", "../config/Manual")
	bt.PrepareData()
	assert.Equal(t, "2017/10/9 9:39", bt.BCM.BeginDate)

}
