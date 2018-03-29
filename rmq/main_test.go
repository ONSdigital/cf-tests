package main

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestElastiCacheConnectAndSet(t *testing.T) {
	SetEnv()
	err := PerformTest(FakeFactory, "test-rmq")
	require.NoError(t, err)
}

func TestWeb(t *testing.T) {
	SetEnv()
	w := httptest.NewRecorder()
	handler := WebHandler(FakeFactory, "test-rmq")
	req := httptest.NewRequest("GET", "http://x/", nil)
	handler(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "RMQ service is OK", string(body))
}

func TestGetURI(t *testing.T) {
	SetEnv()
	ssl, uri, err := GetURI("test-rmq")
	require.NoError(t, err)
	assert.False(t, ssl)
	assert.Equal(t, "amqp://foobar", uri)
}

func SetEnv() {
	vcap_services := `{
			"rabbitmq": [
			 {
			  "credentials": {
			   "ssl": false,
			   "uri": "amqp://foobar"
			  },
			  "label": "rabbitmq",
			  "name": "test-rmq",
			  "tags": [
			   "rabbitmq"
			  ]
			 }
			]
	}`
	os.Setenv("VCAP_SERVICES", vcap_services)
	os.Setenv("VCAP_APPLICATION", "{}")
}
