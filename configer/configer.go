package configer

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"
)

// configer struct to deal with yaml config file
type configer struct {
	// the config file path
	path string
	// the config file content
	content []byte
	// the config file content after unmarshaling
	unmarshaledContent map[string]interface{}
}

// New returns a new configer
func New(path string) *configer {
	return &configer{
		path: path,
	}
}

// Load loads the config file
func (c *configer) Load() error {
	content, err := ioutil.ReadFile(c.path)
	if err != nil {
		return err
	}
	c.content = content
	return nil
}

// Unmarshal unmarshals the config file content
func (c *configer) Unmarshal() error {
	err := yaml.Unmarshal(c.content, &c.unmarshaledContent)
	if err != nil {
		return err
	}
	return nil
}

// Unmarshal method for slice
func (c *configer) UnmarshalSlice_old() ([]interface{}, error) {
	var res []interface{}
	err := yaml.Unmarshal(c.content, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Unmarshal method for slice
func (c *configer) UnmarshalSlice(sec string) ([]interface{}, error) {
	err := yaml.Unmarshal(c.content, &c.unmarshaledContent)
	if err != nil {
		return nil, err
	}
	tmp := c.GetContent()[sec].([]interface{})

	// fmt.Println(tmp)
	return tmp, nil

}

// GetContent returns the content of the config file
func (c *configer) GetContent() map[string]interface{} {
	return c.unmarshaledContent
}

// Get returns the value of the key
func (c *configer) Get(key string) interface{} {
	return c.unmarshaledContent[key]
}

// GetStringMap returns the value of the key as a map
func (c *configer) GetStringMap(key string) map[string]interface{} {
	// the key could be multiple levels, e.g. "a.b.c"
	// tmp := strings.Split(strings.ToLower(key), ".")
	tmp := strings.Split(key, ".")
	tmpmap := c.unmarshaledContent
	for _, v := range tmp {
		tmpmap = tmpmap[v].(map[string]interface{})
	}

	return tmpmap
}

// GetStringSlice returns the value of the key as a slice of string
func (c *configer) GetStringSlice(key string) []string {
	var res []string
	// the key could be multiple levels, e.g. "a.b.c"
	// tmp := strings.Split(strings.ToLower(key), ".")
	tmp := strings.Split(key, ".")
	tmpmap := c.unmarshaledContent

	// iter the tmp, when the last item is reached, return the value
	for i, v := range tmp {
		if i == len(tmp)-1 {
			// check if the value is a slice,if not , make it into a slice
			switch tmpmap[v].(type) {
			case []interface{}:
				for _, v := range tmpmap[v].([]interface{}) {
					res = append(res, v.(string))
				}
			case string:
				res = []string{tmpmap[v].(string)}
			default:
				res = tmpmap[v].([]string)
			}

		} else {
			tmpmap = tmpmap[v].(map[string]interface{})
		}

	}
	return res
}

func (c *configer) GetIntSlice(key string) []int {
	var res []int
	// the key could be multiple levels, e.g. "a.b.c"
	// tmp := strings.Split(strings.ToLower(key), ".")
	tmp := strings.Split(key, ".")
	tmpmap := c.unmarshaledContent

	// iter the tmp, when the last item is reached, return the value
	for i, v := range tmp {
		if i == len(tmp)-1 {
			// check if the value is a slice,if not , make it into a slice
			switch tmpmap[v].(type) {
			case []interface{}:
				for _, v := range tmpmap[v].([]interface{}) {
					res = append(res, v.(int))
				}
			case int:
				res = []int{tmpmap[v].(int)}
			default:
				res = tmpmap[v].([]int)
			}

		} else {
			tmpmap = tmpmap[v].(map[string]interface{})
		}

	}
	return res
}

// GetFloat64 returns the value of the key as a float64
func (c *configer) GetFloat64(key string) float64 {
	// the key could be multiple levels, e.g. "a.b.c"
	// tmp := strings.Split(strings.ToLower(key), ".")
	tmp := strings.Split(key, ".")
	tmpmap := c.unmarshaledContent
	var res float64
	// iter the tmp, when the last item is reached, return the value
	for i, v := range tmp {
		if i == len(tmp)-1 {
			// check the type of the value and convert it to float64
			switch tmpmap[v].(type) {
			case int:
				res = float64(tmpmap[v].(int))
			case int64:
				res = float64(tmpmap[v].(int64))
			case float64:
				res = tmpmap[v].(float64)
			}
		} else {
			tmpmap = tmpmap[v].(map[string]interface{})
		}
	}
	return res
}

func (c *configer) GetInt(key string) int {
	// the key could be multiple levels, e.g. "a.b.c"
	// tmp := strings.Split(strings.ToLower(key), ".")
	tmp := strings.Split(key, ".")
	tmpmap := c.unmarshaledContent
	var res int
	// iter the tmp, when the last item is reached, return the value
	for i, v := range tmp {
		if i == len(tmp)-1 {
			// check the type of the value and convert it to float64
			switch tmpmap[v].(type) {
			case int:
				res = tmpmap[v].(int)
			case int64:
				res = int(tmpmap[v].(int64))
			case float64:
				res = int(tmpmap[v].(float64))
			}
		} else {
			tmpmap = tmpmap[v].(map[string]interface{})
		}
	}
	return res
}

// GetString returns the value of the key as a string
func (c *configer) GetString(key string) string {
	// the key could be multiple levels, e.g. "a.b.c"
	// tmp := strings.Split(strings.ToLower(key), ".")
	tmp := strings.Split(key, ".")
	tmpmap := c.unmarshaledContent
	var res string
	// iter the tmp, when the last item is reached, return the value
	for i, v := range tmp {
		if i == len(tmp)-1 {
			res = tmpmap[v].(string)
		} else {
			tmpmap = tmpmap[v].(map[string]interface{})
		}
	}
	return res
}
