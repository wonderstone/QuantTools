package dataprocessor

// test for csvprocessor_daily.go
import (
	"testing"
)

func TestCsvProcessor(t *testing.T) {
	yamldir := "../config/Manual/"
	BTfilename := "BackTest.yaml"
	STGfilename := "Strategy.yaml"
	Data_pretreatment(yamldir, BTfilename, STGfilename)

}
