package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq"
	pb "github.com/vhive-serverless/vSwarm-proto/proto/hotel_reserv/profile"
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
		CREATE TABLE IF NOT EXISTS hotels (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			phone_number TEXT NOT NULL,
			description TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create addresses table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS addresses (
			hotel_id TEXT PRIMARY KEY REFERENCES hotels(id) ON DELETE CASCADE,
			street_number TEXT NOT NULL,
			street_name TEXT NOT NULL,
			city TEXT NOT NULL,
			state TEXT NOT NULL,
			country TEXT NOT NULL,
			postal_code TEXT NOT NULL,
			lat REAL NOT NULL,
			lon REAL NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create indexes
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_hotels_id ON hotels(id)`)
	if err != nil {
		log.Fatal(err)
	}
}

func populateInitialData(db *sql.DB) {
	// Clear existing data
	_, err := db.Exec("DELETE FROM addresses")
	if err != nil {
		log.Print("Error clearing addresses table: ", err)
	}

	_, err = db.Exec("DELETE FROM hotels")
	if err != nil {
		log.Print("Error clearing hotels table: ", err)
	}
	hotels := []*pb.Hotel{
		{
			Id:          "1",
			Name:        "Clift Hotel",
			PhoneNumber: "(415) 775-4700",
			Description: "A 6-minute walk from Union Square and 4 minutes from a Muni Metro station, this luxury hotel designed by Philippe Starck features an artsy furniture collection in the lobby, including work by Salvador Dali.",
			Address: &pb.Address{
				StreetNumber: "495",
				StreetName:   "Geary St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94102",
				Lat:          37.7867,
				Lon:          -122.4112,
			},
		},
		{
			Id:          "2",
			Name:        "W San Francisco",
			PhoneNumber: "(415) 777-5300",
			Description: "Less than a block from the Yerba Buena Center for the Arts, this trendy hotel is a 12-minute walk from Union Square.",
			Address: &pb.Address{
				StreetNumber: "181",
				StreetName:   "3rd St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94103",
				Lat:          37.7854,
				Lon:          -122.4005,
			},
		},
		{
			Id:          "3",
			Name:        "Hotel Zetta",
			PhoneNumber: "(415) 543-8555",
			Description: "A 3-minute walk from the Powell Street cable-car turnaround and BART rail station, this hip hotel 9 minutes from Union Square combines high-tech lodging with artsy touches.",
			Address: &pb.Address{
				StreetNumber: "55",
				StreetName:   "5th St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94103",
				Lat:          37.7834,
				Lon:          -122.4071,
			},
		},
		{
			Id:          "4",
			Name:        "Hotel Vitale",
			PhoneNumber: "(415) 278-3700",
			Description: "This waterfront hotel with Bay Bridge views is 3 blocks from the Financial District and a 4-minute walk from the Ferry Building.",
			Address: &pb.Address{
				StreetNumber: "8",
				StreetName:   "Mission St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94105",
				Lat:          37.7936,
				Lon:          -122.3930,
			},
		},
		{
			Id:          "5",
			Name:        "Phoenix Hotel",
			PhoneNumber: "(415) 776-1380",
			Description: "Located in the Tenderloin neighborhood, a 10-minute walk from a BART rail station, this retro motor lodge has hosted many rock musicians and other celebrities since the 1950s. It's a 4-minute walk from the historic Great American Music Hall nightclub.",
			Address: &pb.Address{
				StreetNumber: "601",
				StreetName:   "Eddy St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94109",
				Lat:          37.7831,
				Lon:          -122.4181,
			},
		},
		{
			Id:          "6",
			Name:        "St. Regis San Francisco",
			PhoneNumber: "(415) 284-4000",
			Description: "St. Regis Museum Tower is a 42-story, 484 ft skyscraper in the South of Market district of San Francisco, California, adjacent to Yerba Buena Gardens, Moscone Center, PacBell Building and the San Francisco Museum of Modern Art.",
			Address: &pb.Address{
				StreetNumber: "125",
				StreetName:   "3rd St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94109",
				Lat:          37.7863,
				Lon:          -122.4015,
			},
		},
	}

	// Insert initial hotels
	for _, hotel := range hotels {
		insertHotel(db, hotel)
	}

	// Add additional hotels (7-80)
	for i := 7; i <= 80; i++ {
		hotelID := strconv.Itoa(i)

		hotel := pb.Hotel{
			Id:          hotelID,
			Name:        "St. Regis San Francisco",
			PhoneNumber: "(415) 284-40" + hotelID,
			Description: "St. Regis Museum Tower is a 42-story, 484 ft skyscraper in the South of Market district of San Francisco, California, adjacent to Yerba Buena Gardens, Moscone Center, PacBell Building and the San Francisco Museum of Modern Art.",
			Address: &pb.Address{
				StreetNumber: "125",
				StreetName:   "3rd St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94109",
				Lat:          37.7835 + float32(i)/500.0*3,
				Lon:          -122.41 + float32(i)/500.0*4,
			},
		}

		insertHotel(db, &hotel)
	}
}

func insertHotel(db *sql.DB, hotel *pb.Hotel) {
	count := 0
	err := db.QueryRow("SELECT COUNT(*) FROM hotels WHERE id = $1", hotel.Id).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count != 0 {
		log.Fatal("Hotel id alread exists")
	}

	// Insert hotel
	_, err = db.Exec(`
		INSERT INTO hotels (id, name, phone_number, description) 
		VALUES ($1, $2, $3, $4)
	`, hotel.Id, hotel.Name, hotel.PhoneNumber, hotel.Description)
	if err != nil {
		log.Fatal(err)
	}

	// Insert address
	_, err = db.Exec(`
		INSERT INTO addresses (hotel_id, street_number, street_name, city, state, country, postal_code, lat, lon) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, hotel.Id, hotel.Address.StreetNumber, hotel.Address.StreetName, hotel.Address.City,
		hotel.Address.State, hotel.Address.Country, hotel.Address.PostalCode, hotel.Address.Lat, hotel.Address.Lon)
	if err != nil {
		log.Fatal(err)
	}
}

// GetHotelByID fetches a hotel by its ID
func getHotelByID(db *sql.DB, id string) (*pb.Hotel, error) {
	hotel := &pb.Hotel{
		Address: &pb.Address{},
	}

	err := db.QueryRow(`
		SELECT h.id, h.name, h.phone_number, h.description,
		       a.street_number, a.street_name, a.city, a.state, a.country, a.postal_code, a.lat, a.lon
		FROM hotels h
		JOIN addresses a ON h.id = a.hotel_id
		WHERE h.id = $1
	`, id).Scan(
		&hotel.Id, &hotel.Name, &hotel.PhoneNumber, &hotel.Description,
		&hotel.Address.StreetNumber, &hotel.Address.StreetName, &hotel.Address.City,
		&hotel.Address.State, &hotel.Address.Country, &hotel.Address.PostalCode,
		&hotel.Address.Lat, &hotel.Address.Lon,
	)

	if err != nil {
		return nil, err
	}

	return hotel, nil
}