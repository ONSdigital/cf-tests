package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
)

// DAO is the database layer access interface
type DAO interface {
	Open(host, user, password, dbName string) (*sql.DB, error)
	CreateTable(db *sql.DB, tableName, name string) error
	QueryTable(db *sql.DB, tableName string) (string, error)
}

// Credentialiser is an abstraction for reading credentials from VCAP services
type Credentialiser interface {
	GetCreds(serviceName string) (host, user, password, dbName string, err error)
}

func main() {
	dao := &PostgresDAO{}
	port := os.Getenv("PORT")
	serviceName := os.Getenv("DB_SERVICENAME")
	creds := &CFCredentialiser{}
	handler := WebHandler(dao, creds, serviceName, "test_table", "Fred")
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// WebHandler provides a test endpoint
func WebHandler(dao DAO, creds Credentialiser, serviceName, tableName, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := Query(dao, creds, serviceName, tableName, name)
		if err != nil || name != result {
			w.WriteHeader(http.StatusFailedDependency)
			fmt.Fprintf(w, "Failed to read database: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "RDS service is OK")
	}
}

type CFCredentialiser struct{}

func (CFCredentialiser) GetCreds(serviceName string) (host, user, password, dbName string, err error) {
	app, err := cfenv.Current()
	if err != nil {
		return
	}

	pg, err := app.Services.WithName(serviceName)
	if err != nil {
		return
	}

	host, _ = pg.CredentialString("host")
	user, _ = pg.CredentialString("username")
	password, _ = pg.CredentialString("password")
	dbName, _ = pg.CredentialString("db_name")

	return
}

// Query connects to a postgres services and runs a basic query on the test table
func Query(dao DAO, creds Credentialiser, serviceName, tableName, name string) (string, error) {
	host, user, password, dbName, err := creds.GetCreds(serviceName)
	if err != nil {
		return "", err
	}

	db, err := dao.Open(host, user, password, dbName)
	if err != nil {
		return "", err
	}

	if err := dao.CreateTable(db, tableName, name); err != nil {
		return "", err
	}

	return dao.QueryTable(db, tableName)
}

// PostgresDAO is a specific dao for postgres
type PostgresDAO struct{}

// Open creates a connection to a postgres instance
func (PostgresDAO) Open(host, user, password, dbName string) (*sql.DB, error) {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		host, user, password, dbName)
	return sql.Open("postgres", dbinfo)
}

// CreateTable creates a simple test table in the attached database
func (PostgresDAO) CreateTable(db *sql.DB, tableName, name string) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}

	}()

	if _, err = tx.Exec("DROP TABLE IF EXISTS " + tableName); err != nil {
		return
	}

	if _, err = tx.Exec("CREATE TABLE " + tableName + "(name VARCHAR(30) primary key)"); err != nil {
		return
	}

	_, err = tx.Exec("INSERT INTO "+tableName+"(name) VALUES($1)", name)
	return
}

// QueryTable runs a simple query on the test table and returns the first row
func (PostgresDAO) QueryTable(db *sql.DB, tableName string) (name string, err error) {
	rows, err := db.Query("SELECT name FROM " + tableName + " LIMIT 1")
	if err != nil {
		return
	}

	if !rows.Next() {
		return "", errors.New("No rows found")
	}

	err = rows.Scan(&name)
	return
}
