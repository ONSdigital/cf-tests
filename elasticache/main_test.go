package main

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FakeDAO struct {
	URL          string
	Password     string
	ConnectError error
	SetError     error
	GetError     error

	store map[string]string
}

func NewFakeDAO() *FakeDAO {
	return &FakeDAO{store: make(map[string]string)}
}

func (f *FakeDAO) Connect(url, password string) error {
	f.URL = url
	f.Password = password
	return f.ConnectError
}

func (f *FakeDAO) SetValue(label, value string) error {
	f.store[label] = value
	return f.SetError
}

func (f *FakeDAO) GetValue(label string) (string, error) {
	return f.store[label], f.GetError
}

func (f *FakeDAO) UnsetValue(label string) error {
	delete(f.store, label)
	return nil
}

func (f *FakeDAO) Close() {}

func setupFake() (*FakeDAO, CFCredentialiser) {
	dao := NewFakeDAO()
	vcap_services := `
    {
        "elasticache": [
          {
            "credentials": {
			  "host": "redis_host",
			  "port": 6379,
              "password": "redis_password"
            },
            "label": "elasticache",
            "name": "test-elasticache"
          }
        ]
    }
    `
	os.Setenv("VCAP_SERVICES", vcap_services)
	os.Setenv("VCAP_APPLICATION", "{}")
	os.Setenv("ELASTICACHE_SERVICE_NAME", "test-elasticache")
	return dao, CFCredentialiser{}
}

func teardownFake() {
	os.Unsetenv("VCAP_SERVICES")
	os.Unsetenv("VCAP_APPLICATION")
	os.Unsetenv("ELASTICACHE_SERVICENAME")
}

func TestElastiCacheConnectAndSet(t *testing.T) {
	dao, creds := setupFake()
	tester := NewTester(dao, creds)
	err := tester.PerformTest("test-elasticache")
	require.NoError(t, err)
	assert.Equal(t, "redis_host:6379", dao.URL)
	assert.Equal(t, "redis_password", dao.Password)
}

func TestWeb(t *testing.T) {
	dao, creds := setupFake()
	defer teardownFake()
	w := httptest.NewRecorder()
	handler := WebHandler(dao, creds)
	req := httptest.NewRequest("GET", "http://x/", nil)
	handler(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "Elasticache service is OK", string(body))
}
