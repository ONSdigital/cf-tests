package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/go-redis/redis"
)

type DAO interface {
	Connect(url, password string) error
	SetValue(label, value string) error
	GetValue(label string) (string, error)
	UnsetValue(label string) error
	Close()
}

func main() {
	dao := &RedisDAO{}
	creds := CFCredentialiser{}
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, WebHandler(dao, creds)))
}

type Tester struct {
	dao   DAO
	creds CFCredentialiser
}

func NewTester(dao DAO, creds CFCredentialiser) *Tester {
	return &Tester{dao: dao, creds: creds}
}

func WebHandler(dao DAO, creds CFCredentialiser) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		serviceName := os.Getenv("ELASTICACHE_SERVICE_NAME")
		if err := NewTester(dao, creds).PerformTest(serviceName); err != nil {
			w.WriteHeader(http.StatusFailedDependency)
			fmt.Fprintf(w, "Failed to access ElastiCache: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Elasticache service is OK")
	}
}

func (t *Tester) PerformTest(serviceName string) error {
	uri, password, err := t.creds.GetCreds(serviceName)
	if err != nil {
		return err
	}
	if err := t.dao.Connect(uri, password); err != nil {
		return err
	}
	if err := t.dao.SetValue("foo", "bar"); err != nil {
		return err
	}
	defer t.dao.UnsetValue("foo")
	if value, err := t.dao.GetValue("foo"); err != nil {
		return err
	} else if value != "bar" {
		return fmt.Errorf("Value set but not retrieved")
	}
	return nil
}

type CFCredentialiser struct {
}

func (CFCredentialiser) GetCreds(serviceName string) (uri, password string, err error) {
	app, err := cfenv.Current()
	if err != nil {
		return
	}

	pg, err := app.Services.WithName(serviceName)
	if err != nil {
		return
	}

	host, _ := pg.CredentialString("host")
	port, _ := pg.Credentials["port"].(float64)
	uri = fmt.Sprintf("%s:%0.0f", host, port)
	password, _ = pg.CredentialString("password")

	return
}

type RedisDAO struct {
	client *redis.Client
}

func (r *RedisDAO) Connect(uri, password string) error {
	r.client = redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: password,
		DB:       0,
	})
	_, err := r.client.Ping().Result()
	return err
}

func (r *RedisDAO) SetValue(label, value string) error {
	return r.client.Set(label, value, 0).Err()
}

func (r *RedisDAO) GetValue(label string) (string, error) {
	return r.client.Get(label).Result()
}

func (r *RedisDAO) UnsetValue(label string) error {
	_, err := r.client.Del(label).Result()
	return err
}

func (r *RedisDAO) Close() {
	if r.client != nil {
		r.client.Close()
	}
}
