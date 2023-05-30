package dataprocessor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test acquire_token
func TestAcquireToken(t *testing.T) {
	//test data
	url := "http://123.138.216.197:9002/xbzq/vds/v1/user/login"
	var user_login Login
	user_login.Uname = "admin"
	user_login.Upwd = "123456"
	token := acquire_token(&user_login, url)
	fmt.Println(token)
	// test the token is not empty
	assert.NotEqual(t, token, "")
}

func TestDataDownload(t *testing.T) {
	hr := historyData_request{
		Symbol:     "688003.SH",
		StartDt:    "20230122180000000",
		EndDt:      "20230525180000000",
		Count:      20000,
		Field:      "*",
		CandleType: "1min",
	}
	vdsdata := historydata_download(&hr, "admin", "123456")
	csv_download(vdsdata, "./", "Template.csv")
}
