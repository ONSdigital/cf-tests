package main

import (
	"database/sql"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestDAOCreateTable(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DROP TABLE IF EXISTS test_data").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE test_data").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO test_data").WithArgs("Fred").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	dao := &PostgresDAO{}
	err = dao.CreateTable(db, "test_data", "Fred")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDAOQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Fred")
	mock.ExpectQuery("SELECT name FROM test_data LIMIT 1").WillReturnRows(rows)

	dao := &PostgresDAO{}
	name, err := dao.QueryTable(db, "test_data")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, "Fred", name)
}

type FakeDAO struct {
	Host        string
	User        string
	Password    string
	DBName      string
	TableName   string
	Name        string
	OpenError   error
	CreateError error
	QueryError  error
}

func (f *FakeDAO) Open(host, user, password, dbName string) (*sql.DB, error) {
	f.Host = host
	f.User = user
	f.Password = password
	f.DBName = dbName
	return nil, f.OpenError
}

func (f *FakeDAO) CreateTable(_ *sql.DB, tableName, name string) error {
	f.TableName = tableName
	f.Name = name
	return f.CreateError
}

func (f *FakeDAO) QueryTable(_ *sql.DB, tableName string) (string, error) {
	f.TableName = tableName
	return f.Name, f.QueryError
}

func setupFake() (*FakeDAO, Credentialiser) {
	dao := &FakeDAO{}
	vcap_services := `
	{
		"rds": [
		  {
			"credentials": {
			  "db_name": "test_db",
			  "host": "test_host",
			  "password": "test_password",
			  "uri": "you don't want to use this",
			  "username": "test_user"
			},
			"label": "rds",
			"name": "test-psql"
		  }
		]
	}
	`
	os.Setenv("VCAP_SERVICES", vcap_services)
	os.Setenv("VCAP_APPLICATION", "{}")
	os.Setenv("DB_SERVICENAME", "test-psql")
	return dao, &CFCredentialiser{}
}

func teardownFake() {
	os.Unsetenv("VCAP_SERVICES")
	os.Unsetenv("VCAP_APPLICATION")
	os.Unsetenv("DB_SERVICENAME")
}

func TestQuery(t *testing.T) {
	dao, creds := setupFake()
	defer teardownFake()
	name, err := Query(dao, creds, "test-psql", "test_data", "Fred")
	require.NoError(t, err)
	assert.Equal(t, "Fred", name)
	assert.Equal(t, "test_host", dao.Host)
	assert.Equal(t, "test_user", dao.User)
	assert.Equal(t, "test_password", dao.Password)
	assert.Equal(t, "test_db", dao.DBName)
	assert.Equal(t, "test_data", dao.TableName)
	assert.Equal(t, "Fred", dao.Name)
}

func TestWeb(t *testing.T) {
	dao, creds := setupFake()
	defer teardownFake()
	w := httptest.NewRecorder()
	handler := WebHandler(dao, creds, "test-psql", "test_data", "Fred")
	req := httptest.NewRequest("GET", "http://x/", nil)
	handler(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "RDS service is OK", string(body))
}
