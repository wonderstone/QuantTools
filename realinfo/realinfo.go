package realinfo

import (
	"github.com/wonderstone/QuantTools/configer"
)

type Info struct {
	IM map[string]interface{}
}

func NewInfo(info map[string]interface{}) *Info {
	return &Info{IM: info}
}

// NewInfoFromConfig reads the configuration file and returns a Info struct
// filename: accountinfo.yaml
func NewInfoFromConfig(configpath string, filename string) *Info {
	c := configer.New(configpath + filename)
	err := c.Load()
	if err != nil {
		panic(err)
	}
	err = c.Unmarshal()
	if err != nil {
		panic(err)
	}
	cm := c.GetContent()

	return NewInfo(cm)
}
