package realinfo

import "github.com/spf13/viper"

type Info struct {
	IM map[string]interface{}
}

func NewInfo(info map[string]interface{}) *Info {
	return &Info{IM: info}
}

// NewInfoFromConfig reads the configuration file and returns a Info struct
// filename: accountinfo.yaml
func NewInfoFromConfig(configpath string, filename string) *Info {
	// read accountinfo configuration from file  viper is not thread safe
	viper.SetConfigName(filename)
	viper.AddConfigPath(configpath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	cm := viper.AllSettings()
	return NewInfo(cm)
}
