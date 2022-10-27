package configer

import (
	"testing"
)

func TestConfiger(t *testing.T) {
	c := New("../config/Manual/accountinfo.yaml")
	err := c.Load()
	if err != nil {
		t.Error(err)
	}
	err = c.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	// test the GetContent method
	content := c.GetContent()
	if len(content) == 0 {
		t.Error("content is empty")
	}

}
