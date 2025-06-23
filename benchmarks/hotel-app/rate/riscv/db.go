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
	pb "github.com/vhive-serverless/vSwarm-proto/proto/hotel_reserv/rate"
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
	// Create rate_plans table
	_, err := db.Exec(`
		DROP TABLE IF EXISTS rate_plans CASCADE;
		CREATE TABLE rate_plans (
			hotel_id TEXT NOT NULL,
			code TEXT NOT NULL,
			in_date DATE NOT NULL,
			out_date DATE NOT NULL,
			room_type_code TEXT NOT NULL,
			room_description TEXT NOT NULL,
			bookable_rate DECIMAL(10, 2) NOT NULL,
			total_rate DECIMAL(10, 2) NOT NULL,
			total_rate_inclusive DECIMAL(10, 2) NOT NULL,
			PRIMARY KEY (hotel_id, code, in_date, out_date)
		)
	`)
	if err != nil {
		log.Fatalf("Error creating rate_plans table: %s", err.Error())
	}

	// Create indexes
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_rate_plans_hotel_id ON rate_plans(hotel_id)`)
	if err != nil {
		log.Fatalf("Error creating index on rate_plans: %s", err.Error())
	}
}

func populateInitialData(db *sql.DB) {

	// Clear existing data
	_, err := db.Exec("DELETE FROM rate_plans")
	if err != nil {
		log.Print("Error clearing rate_plans table: ", err)
	}

	rate_plans := []*pb.RatePlan{
		{
			HotelId: "1",
			Code:    "RACK",
			InDate:  "2015-04-09",
			OutDate: "2015-04-10",
			RoomType: &pb.RoomType{
				BookableRate:       109.00,
				Code:               "KNG",
				RoomDescription:    "King sized bed",
				TotalRate:          109.00,
				TotalRateInclusive: 123.17,
			},
		},
		{
			HotelId: "2",
			Code:    "RACK",
			InDate:  "2015-04-09",
			OutDate: "2015-04-10",
			RoomType: &pb.RoomType{
				BookableRate:       139.00,
				Code:               "QN",
				RoomDescription:    "Queen sized bed",
				TotalRate:          139.00,
				TotalRateInclusive: 153.09,
			},
		},
		{
			HotelId: "3",
			Code:    "RACK",
			InDate:  "2015-04-09",
			OutDate: "2015-04-10",
			RoomType: &pb.RoomType{
				BookableRate:       109.00,
				Code:               "KNG",
				RoomDescription:    "King sized bed",
				TotalRate:          109.00,
				TotalRateInclusive: 123.17,
			},
		},
	}

	// Insert initial hotels
	for _, rate_plan := range rate_plans {
		insertRatePlan(db, rate_plan)
	}

	// add up to 80 hotels
	item := &pb.RatePlan{
		RoomType: &pb.RoomType{},
	}

	for i := 7; i <= 80; i++ {
		if i%3 == 0 {
			hotel_id := strconv.Itoa(i)
			end_date := "2015-04-"
			rate := 109.00
			rate_inc := 123.17
			if i%2 == 0 {
				end_date = end_date + "17"
			} else {
				end_date = end_date + "24"
			}

			if i%5 == 1 {
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

			item.HotelId = hotel_id
			item.Code = "RACK"
			item.InDate = "2015-04-09"
			item.OutDate = end_date
			item.RoomType.BookableRate = rate
			item.RoomType.Code = "KNG"
			item.RoomType.RoomDescription = "King sized bed"
			item.RoomType.TotalRate = rate
			item.RoomType.TotalRateInclusive = rate_inc

			insertRatePlan(db, item)
		}
	}
}

func insertRatePlan(db *sql.DB, rate_plan *pb.RatePlan) {
	_, err := db.Exec(`
		INSERT INTO rate_plans (
			hotel_id, code, in_date, out_date, 
			room_type_code, room_description, 
			bookable_rate, total_rate, total_rate_inclusive
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (hotel_id, code, in_date, out_date) DO NOTHING
	`, rate_plan.HotelId, rate_plan.Code, rate_plan.InDate, rate_plan.OutDate, rate_plan.RoomType.Code, rate_plan.RoomType.RoomDescription, rate_plan.RoomType.BookableRate, rate_plan.RoomType.TotalRate, rate_plan.RoomType.TotalRateInclusive)

	if err != nil {
		log.Fatalf("Error inserting rate plan data: %s", err.Error())
	}
}

// GetHotelByID fetches a hotel by its ID
func getHotelByID(db *sql.DB, id string) ([]*pb.RatePlan, error) {
	rate_plans := []*pb.RatePlan{}
	rows, err := db.Query(`
		SELECT hotel_id, code, in_date, out_date, room_type_code, room_description, 
			bookable_rate, total_rate, total_rate_inclusive
		FROM rate_plans
		WHERE hotel_id = $1
	`, id)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		rate_plan := &pb.RatePlan{
			RoomType: &pb.RoomType{},
		}
		if err := rows.Scan(
			&rate_plan.HotelId, &rate_plan.Code, &rate_plan.InDate, &rate_plan.OutDate, &rate_plan.RoomType.Code, &rate_plan.RoomType.RoomDescription, &rate_plan.RoomType.BookableRate, &rate_plan.RoomType.TotalRate, &rate_plan.RoomType.TotalRateInclusive); err != nil {
			return rate_plans, err
		}
		rate_plans = append(rate_plans, rate_plan)
	}

	if err = rows.Err(); err != nil {
		return rate_plans, err
	}
	return rate_plans, nil
}
