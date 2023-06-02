package dataprocessor

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeParse(t *testing.T) {
	tm, err := time.Parse("2006.01.02T15:04:05.000", "2023.01.18T09:35:00.000")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tm)
	tm, err = time.Parse("20060102150405000", "20230118093500000")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tm)
}
