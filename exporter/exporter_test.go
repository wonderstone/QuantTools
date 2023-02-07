package exporter

import (
	"testing"

	"github.com/wonderstone/QuantTools/account/virtualaccount"
)

// test ExportRealtimeYaml function
func TestExportRealtimeYaml(t *testing.T) {
	va := virtualaccount.NewVirtualAccount("2017/10/9 9:39", 100000, 100000)
	va.SAcct.ResetMVSlice()

	ExportRealtimeYaml("../config/Manual/", "BackTest", "Default", va)
}

func TestReplaceVA(t *testing.T) {
	va := virtualaccount.NewVirtualAccount("2017/10/9 9:39", 200000, 100000)
	ReplaceVA("../config/Manual/", "realtime", va)
}

func TestExportSKE(t *testing.T) {
	ExportSKE("../config/Manual/", "GEP", "GEP", []string{"S1"})
}
