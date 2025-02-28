package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gocql/gocql"
	pb "github.com/vhive-serverless/vSwarm-proto/proto/hotel_reserv/profile"
)

func initializeDatabase(url string) *gocql.Session {
	fmt.Printf("profile db ip addr = %s\n", url)

	// Create a cluster configuration
	cluster := gocql.NewCluster(url)
	cluster.Keyspace = "profile_keyspace"
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = time.Second * 10

	// Connect to the cluster
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra cluster: %v", err)
		return nil
	}

	// Create keyspace if not exists
	err = session.Query(`CREATE KEYSPACE IF NOT EXISTS profile_keyspace 
		WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1}`).Exec()
	if err != nil {
		log.Printf("Error creating keyspace: %v", err)
	}

	// Create table if not exists
	err = session.Query(`CREATE TABLE IF NOT EXISTS profile_keyspace.hotels (
		id text PRIMARY KEY,
		name text,
		phone_number text,
		description text,
		street_number text,
		street_name text,
		city text,
		state text,
		country text,
		postal_code text,
		lat float,
		lon float
	)`).Exec()
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	// Clear the table to have a fresh start, similar to MongoDB's DropCollection
	err = session.Query(`TRUNCATE profile_keyspace.hotels`).Exec()
	if err != nil {
		log.Printf("Error truncating table: %v", err)
	}

	// Insert initial hotel data
	insertHotel(session, &pb.Hotel{
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
	})

	insertHotel(session, &pb.Hotel{
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
	})

	insertHotel(session, &pb.Hotel{
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
	})

	insertHotel(session, &pb.Hotel{
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
	})

	insertHotel(session, &pb.Hotel{
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
	})

	insertHotel(session, &pb.Hotel{
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
	})

	// Add hotels 7-80
	for i := 7; i <= 80; i++ {
		hotelID := strconv.Itoa(i)
		phoneNum := "(415) 284-40" + hotelID
		lat := 37.7835 + float32(i)/500.0*3
		lon := -122.41 + float32(i)/500.0*4

		insertHotel(session, &pb.Hotel{
			Id:          hotelID,
			Name:        "St. Regis San Francisco",
			PhoneNumber: phoneNum,
			Description: "St. Regis Museum Tower is a 42-story, 484 ft skyscraper in the South of Market district of San Francisco, California, adjacent to Yerba Buena Gardens, Moscone Center, PacBell Building and the San Francisco Museum of Modern Art.",
			Address: &pb.Address{
				StreetNumber: "125",
				StreetName:   "3rd St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94109",
				Lat:          lat,
				Lon:          lon,
			},
		})
	}

	return session
}

func insertHotel(session *gocql.Session, hotel *pb.Hotel) {
	// Insert hotel data into Cassandra
	err := session.Query(`
		INSERT INTO profile_keyspace.hotels (
			id, name, phone_number, description, 
			street_number, street_name, city, state, 
			country, postal_code, lat, lon
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		hotel.Id, hotel.Name, hotel.PhoneNumber, hotel.Description,
		hotel.Address.StreetNumber, hotel.Address.StreetName, hotel.Address.City, hotel.Address.State,
		hotel.Address.Country, hotel.Address.PostalCode, hotel.Address.Lat, hotel.Address.Lon,
	).Exec()

	if err != nil {
		log.Printf("Error inserting hotel with ID %s: %v", hotel.Id, err)
	}
}
