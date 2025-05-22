// MIT License

// Copyright (c) 2022 EASE lab

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"

	"strconv"

	_ "github.com/lib/pq"
)

func initializeDatabase(host string, port string, user string, password string, dbname string) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	fmt.Printf("profile db connection string = %s\n", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Error while connecting... (%s)", err.Error())
		return nil
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		fmt.Printf("Error while pinging... (%s)", err.Error())
		return nil
	}

	// Create tables if they don't exist
	createTables(db)

	// Populate initial data
	populateInitialData(db)

	return db
}

func createTables(db *sql.DB) {
	// Create hotels table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create indexes
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_users ON users(username)`)
	if err != nil {
		log.Fatal(err)
	}
}

func populateInitialData(db *sql.DB) {
	users := make([]User, 500)
	users[0].Username = "hello"
	users[0].Password = "hello"

	// Create users
	for i := 1; i < len(users); i++ {
		suffix := strconv.Itoa(i)
		users[i].Username = "user_" + suffix
		users[i].Password = "pass_" + suffix
	}

	// Encrypt password
	for i := range users {
		sum := sha256.Sum256([]byte(users[i].Password))
		pass := fmt.Sprintf("%x", sum)
		users[i].Password = pass
	}

	for _, user := range users {
		insertUser(db, &user)
	}
}

func insertUser(db *sql.DB, user *User) {
	_, err := db.Exec(`
		INSERT INTO users (username, password) VALUES ($1, $2)
	`, user.Username, user.Password)
	if err != nil {
		log.Fatal(err)
	}
}

func getUsers(db *sql.DB) ([]*User, error) {
	users := []*User{}
	rows, err := db.Query(`
		SELECT username, password
		FROM users
	`)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user := &User{}
		if err := rows.Scan(
			&user.Username, &user.Password); err != nil {
			return users, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return users, err
	}
	return users, nil
}

func getUser(db *sql.DB, username string) (*User, error) {
	user := &User{}

	err := db.QueryRow(`
		SELECT username, password
		FROM users
		WHERE username = $1
	`, username).Scan(
		&user.Username, &user.Password,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
