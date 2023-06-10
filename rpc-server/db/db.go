package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

// DB represents the database connection.
type DB struct {
	conn *sql.DB
}

const maxConnectionAttempts = 10

// NewDB creates a new DB instance with the provided MySQL connection details.
func ConnectToDB() (*DB, error) {

	fmt.Println("Connecting to db")
	conn, err := sql.Open("mysql", "root:root@tcp(db:3306)/instant_messaging_app")
	if err != nil {
		log.Fatal(err)
	}

	attempts := 0
	for conn.Ping() != nil {
		fmt.Println("Attempting connection to db")
		time.Sleep(5 * time.Second)
		attempts++
		if attempts >= maxConnectionAttempts {
			return nil, errors.New("Failed to connect to the database")
		}
	}
	fmt.Println("Connected to db")

	return &DB{conn: conn}, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.conn.Close()
}

// Create executes an insert query and returns the ID of the newly inserted row.
func (db *DB) Create(query string, args ...interface{}) (int64, error) {
	result, err := db.conn.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Read executes a select query and returns a row result.
func (db *DB) Read(query string, args ...interface{}) (*sql.Row, error) {
	return db.conn.QueryRow(query, args...), nil
}

// ReadAll executes a select query and returns a result set.
func (db *DB) ReadAll(query string, args ...interface{}) (*sql.Rows, error) {
	return db.conn.Query(query, args...)
}

// Update executes an update query and returns the number of affected rows.
func (db *DB) Update(query string, args ...interface{}) (int64, error) {
	result, err := db.conn.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete executes a delete query and returns the number of affected rows.
func (db *DB) Delete(query string, args ...interface{}) (int64, error) {
	result, err := db.conn.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
