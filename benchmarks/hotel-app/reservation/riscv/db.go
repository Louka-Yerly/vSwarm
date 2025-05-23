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
		CREATE TABLE IF NOT EXISTS numbers (
			hotel_id TEXT PRIMARY KEY,	
			number INTEGER NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS reservations (
			order_id SERIAL PRIMARY KEY,
			hotel_id TEXT NOT NULL,
			customername TEXT NOT NULL,
			indate TEXT NOT NULL,
			outdate TEXT NOT NULL,
			number INTEGER NOT NULL,
			FOREIGN KEY (hotel_id) REFERENCES numbers(hotel_id) ON DELETE CASCADE
		)
	`)

	if err != nil {
		log.Fatal(err)
	}
}

func populateInitialData(db *sql.DB) {

	numbers := []*Number{
		{
			HotelId: "1",
			Number:  200,
		},
		{
			HotelId: "2",
			Number:  10,
		},
		{
			HotelId: "3",
			Number:  200,
		},
		{
			HotelId: "4",
			Number:  200,
		},
		{
			HotelId: "5",
			Number:  200,
		},
		{
			HotelId: "6",
			Number:  200,
		},
	}

	for _, number := range numbers {
		insertNumber(db, number)
	}

	for i := 7; i <= 80; i++ {
		hotel_id := strconv.Itoa(i)
		room_num := 200
		if i%3 == 1 {
			room_num = 300
		} else if i%3 == 2 {
			room_num = 250
		}
		insertNumber(db, &Number{hotel_id, room_num})
	}

	res := &Reservation{"4", "Alice", "2015-04-09", "2015-04-10", 1}

	insertReservation(db, res)

}

func insertNumber(db *sql.DB, number *Number) {
	count := 0
	err := db.QueryRow("SELECT COUNT(*) FROM numbers WHERE hotel_id = $1", number.HotelId).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count != 0 {
		log.Fatal("number id alread exists")
	}

	// Insert number
	_, err = db.Exec(`
		INSERT INTO numbers (hotel_id, number) 
		VALUES ($1, $2)
	`, number.HotelId, number.Number)
	if err != nil {
		log.Fatal(err)
	}
}

func insertReservation(db *sql.DB, reservation *Reservation) {
	count := 0
	err := db.QueryRow(`
	SELECT COUNT(*)
	FROM reservations
	WHERE hotel_id = $1 AND customername = $2 AND indate = $3 AND outdate = $4 AND number = $5
	`, reservation.HotelId, reservation.CustomerName,
		reservation.InDate, reservation.OutDate, reservation.Number).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count != 0 {
		fmt.Printf("Reservation alread exists")
		return
	}

	// Insert reservation
	_, err = db.Exec(`
		INSERT INTO reservations (hotel_id, customername, indate, outdate, number) 
		VALUES ($1, $2, $3, $4, $5)
	`, reservation.HotelId, reservation.CustomerName,
		reservation.InDate, reservation.OutDate, reservation.Number)
	if err != nil {
		log.Fatal(err)
	}
}

func getReservations(db *sql.DB, reservation *Reservation) ([]*Reservation, error) {
	reservations := []*Reservation{}

	rows, err := db.Query(`
		SELECT hotel_id, customername, indate, outdate, number
		FROM reservations
		WHERE hotel_id = $1 AND indate = $2 AND outdate = $3
	`, reservation.HotelId, reservation.InDate, reservation.OutDate)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		res := &Reservation{}
		if err := rows.Scan(
			&res.HotelId, &res.CustomerName, &res.InDate, &res.OutDate, &res.Number); err != nil {
			return reservations, err
		}
		reservations = append(reservations, res)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}

func getNumber(db *sql.DB, hotel_id string) (Number, error) {
	number := Number{}

	err := db.QueryRow(`
		SELECT hotel_id, number
		FROM numbers
		WHERE hotel_id = $1
	`, hotel_id).Scan(
		&number.HotelId, &number.Number,
	)

	return number, err
}
