package main

import (
	"database/sql"
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
	User        string
	Password    string
	DBName      string
	TableName   string
	Name        string
	OpenError   error
	CreateError error
	QueryError  error
}

func (f *FakeDAO) Open(user, password, dbName string) (*sql.DB, error) {
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

func TestQuery(t *testing.T) {
	dao := &FakeDAO{}
	os.Setenv("DB_USER", "foo")
	os.Setenv("DB_PASSWORD", "bar")
	os.Setenv("DB_NAME", "barfoo")
	defer func() {
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	}()
	name, err := Query(dao, "test_data", "Fred")
	require.NoError(t, err)
	assert.Equal(t, "Fred", name)
	assert.Equal(t, "foo", dao.User)
	assert.Equal(t, "bar", dao.Password)
	assert.Equal(t, "barfoo", dao.DBName)
	assert.Equal(t, "test_data", dao.TableName)
	assert.Equal(t, "Fred", dao.Name)
}
