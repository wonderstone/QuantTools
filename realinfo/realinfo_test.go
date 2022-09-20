package realinfo

import (
	"fmt"
	"testing"
)

func TestNewInfoFromConfig(t *testing.T) {
	// NewInfoFromConfig
	i := NewInfoFromConfig("../config/Manual", "accountinfo")
	for k, v := range i.IM {
		fmt.Println(k, v)
		for k1, v1 := range v.(map[string]interface{}) {
			fmt.Println(k1, v1)
		}
	}
}
