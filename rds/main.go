package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

// DAO is the database layer access interface
type DAO interface {
	Open(user, password, dbName string) (*sql.DB, error)
	CreateTable(db *sql.DB, tableName, name string) error
	QueryTable(db *sql.DB, tableName string) (string, error)
}

func main() {
	dao := &PostgresDAO{}
	name, err := Query(dao, "test_table", "Bob")
	if err != nil || name != "Bob" {
		fmt.Println("Failed to read database:", err)
		os.Exit(1)
	}

	fmt.Println("Database read successfully")
}

func Query(dao DAO, tableName, name string) (string, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	db, err := dao.Open(user, password, dbName)
	if err != nil {
		return "", err
	}

	if err := dao.CreateTable(db, tableName, name); err != nil {
		return "", err
	}

	return dao.QueryTable(db, tableName)
}

type PostgresDAO struct{}

func (PostgresDAO) Open(user, password, dbName string) (*sql.DB, error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		user, password, dbName)
	return sql.Open("postgres", dbinfo)
}

func (PostgresDAO) CreateTable(db *sql.DB, dbName, name string) (err error) {
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

	if _, err = tx.Exec("CREATE TABLE " + dbName + "(name VARCHAR(30) primary key)"); err != nil {
		return
	}

	_, err = tx.Exec("INSERT INTO "+dbName+"(name) VALUES(?)", name)
	return
}

func (PostgresDAO) QueryTable(db *sql.DB, dbName string) (name string, err error) {
	rows, err := db.Query("SELECT name FROM " + dbName + " LIMIT 1")
	if err != nil {
		return
	}

	if !rows.Next() {
		return "", errors.New("No rows found")
	}

	err = rows.Scan(&name)
	return
}
