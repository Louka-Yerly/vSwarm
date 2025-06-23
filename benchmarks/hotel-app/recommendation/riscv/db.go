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

type HotelDB struct {
	HId    string
	HLat   float64
	HLon   float64
	HRate  float64
	HPrice float64
}

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
		CREATE TABLE IF NOT EXISTS hotels (
			id SERIAL PRIMARY KEY,	
			hotel_id TEXT NOT NULL,
			lat REAL NOT NULL,
			lon REAL NOT NULL,
			rate REAL NOT NULL,
			price REAL NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create indexes
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_hotels_id ON hotels(hotel_id)`)
	if err != nil {
		log.Fatal(err)
	}
}

func populateInitialData(db *sql.DB) {

	hotels := []*HotelDB{
		{
			HId:    "1",
			HLat:   37.7867,
			HLon:   -122.4112,
			HRate:  109.00,
			HPrice: 150.00,
		},
		{
			HId:    "2",
			HLat:   37.7854,
			HLon:   -122.4005,
			HRate:  139.00,
			HPrice: 120.00,
		},
		{
			HId:    "3",
			HLat:   37.7834,
			HLon:   -122.4071,
			HRate:  109.00,
			HPrice: 190.00,
		},
		{
			HId:    "4",
			HLat:   37.7936,
			HLon:   -122.3930,
			HRate:  129.00,
			HPrice: 160.00,
		},
		{
			HId:    "5",
			HLat:   37.7831,
			HLon:   -122.4181,
			HRate:  119.00,
			HPrice: 140.00,
		},
		{
			HId:    "6",
			HLat:   37.7863,
			HLon:   -122.4015,
			HRate:  149.00,
			HPrice: 200.00,
		},
	}

	for _, hotel := range hotels {
		insertHotel(db, hotel)
	}

	// add up to 80 hotels
	for i := 7; i <= 80; i++ {
		hotel_id := strconv.Itoa(i)

		lat := 37.7835 + float64(i)/500.0*3
		lon := -122.41 + float64(i)/500.0*4

		rate := 135.00
		rate_inc := 179.00
		if i%3 == 0 {
			if i%5 == 0 {
				rate = 109.00
				rate_inc = 123.17
			} else if i%5 == 1 {
				rate = 120.00
				rate_inc = 140.00
			} else if i%5 == 2 {
				rate = 124.00
				rate_inc = 144.00
			} else if i%5 == 3 {
				rate = 132.00
				rate_inc = 158.00
			} else if i%5 == 4 {
				rate = 232.00
				rate_inc = 258.00
			}
		}

		hotel := &HotelDB{
			HId:    hotel_id,
			HLat:   lat,
			HLon:   lon,
			HRate:  rate,
			HPrice: rate_inc,
		}

		insertHotel(db, hotel)
	}

}

func insertHotel(db *sql.DB, hotel *HotelDB) {
	count := 0
	err := db.QueryRow("SELECT COUNT(*) FROM hotels WHERE id = $1", hotel.HId).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count != 0 {
		log.Fatal("Hotel id alread exists")
	}

	// Insert hotel
	_, err = db.Exec(`
		INSERT INTO hotels (hotel_id, lat, lon, rate, price) 
		VALUES ($1, $2, $3, $4, $5)
	`, hotel.HId, hotel.HLat, hotel.HLon, hotel.HRate, hotel.HPrice)
	if err != nil {
		log.Fatal(err)
	}
}

func getHotels(db *sql.DB) ([]*Hotel, error) {
	hotels := []*Hotel{}
	rows, err := db.Query(`
		SELECT hotel_id, lat, lon, rate, price
		FROM hotels
	`)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		hotel := &Hotel{}
		if err := rows.Scan(
			&hotel.HId, &hotel.HLat, &hotel.HLon, &hotel.HRate, &hotel.HPrice); err != nil {
			return hotels, err
		}
		hotels = append(hotels, hotel)
	}

	if err = rows.Err(); err != nil {
		return hotels, err
	}
	return hotels, nil
}
