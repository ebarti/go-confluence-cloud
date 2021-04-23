package confluentcloud

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type apiTestValue struct {
	Input []string
	Error error
}

func TestNewAPI(t *testing.T) {
	testValues := []apiTestValue{
		{[]string{"", "username", "token"}, fmt.Errorf("url empty")},
		{[]string{"https://test.test", "", ""}, nil},
		{[]string{"https://test.test", "username", "token"}, nil},
		{[]string{"test", "username", "token"}, fmt.Errorf("parse \"test\": invalid URI for request")},
	}

	for _, test := range testValues {
		api, err := newAPI(test.Input[0], test.Input[1], test.Input[2])
		if err != nil {
			assert.Equal(t, test.Error.Error(), err.Error())
		} else {
			assert.Equal(t, test.Input[0], api.endPoint.String())
			assert.Equal(t, test.Input[1], api.username)
			assert.Equal(t, test.Input[2], api.token)
		}
	}
}

func TestSetDebug(t *testing.T) {
	assert.False(t, DebugFlag)
	SetDebug(true)
	assert.True(t, DebugFlag)
	SetDebug(false)
	assert.False(t, DebugFlag)
}

func Test_api_VerifyTLS(t *testing.T) {
	a, _ := newAPI("https://test.test", "test", "test")
	a.VerifyTLS(true)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	if !reflect.DeepEqual(a.client.Transport, tr) {
		t.Fail()
	}
	a.VerifyTLS(false)
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if !reflect.DeepEqual(a.client.Transport, tr) {
		t.Fail()
	}
}
