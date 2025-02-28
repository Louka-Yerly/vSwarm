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
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	pb "github.com/vhive-serverless/vSwarm-proto/proto/hotel_reserv/profile"
	tracing "github.com/vhive-serverless/vSwarm/utils/tracing/go"
)

// Server implements the profile service
type Server struct {
	pb.UnimplementedProfileServer

	Port             int
	IpAddr           string
	CassandraSession *gocql.Session
	MemcClient       *memcache.Client
}

// Run starts the server
func (s *Server) Run() error {
	if s.Port == 0 {
		return fmt.Errorf("server port must be set")
	}

	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Timeout: 120 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			PermitWithoutStream: true,
		}),
	}

	if tracing.IsTracingEnabled() {
		opts = append(opts, tracing.GetServerInterceptor())
	}

	srv := grpc.NewServer(opts...)
	pb.RegisterProfileServer(srv, s)

	// Register reflection service on gRPC server.
	reflection.Register(srv)

	// listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Start Profile server. Addr: %s:%d\n", s.IpAddr, s.Port)
	return srv.Serve(lis)
}

// GetProfiles returns hotel profiles for requested IDs
func (s *Server) GetProfiles(ctx context.Context, req *pb.Request) (*pb.Result, error) {
	res := new(pb.Result)
	hotels := make([]*pb.Hotel, 0)

	// Process each requested hotel ID
	for _, hotelID := range req.HotelIds {
		// First check memcached
		item, err := s.MemcClient.Get(hotelID)
		if err == nil {
			// Memcached hit
			hotelProf := new(pb.Hotel)
			if err = json.Unmarshal(item.Value, hotelProf); err != nil {
				log.Warn(err)
			}
			hotels = append(hotels, hotelProf)
		} else if err == memcache.ErrCacheMiss {
			// Memcached miss, query Cassandra
			hotelProf := fetchHotelFromCassandra(s.CassandraSession, hotelID)
			
			if hotelProf != nil {
				hotels = append(hotels, hotelProf)

				// Marshal the hotel profile to JSON for memcached
				profJSON, err := json.Marshal(hotelProf)
				if err != nil {
					log.Warn(err)
				} else {
					// Write to memcached
					err = s.MemcClient.Set(&memcache.Item{Key: hotelID, Value: profJSON})
					if err != nil {
						log.Warn("Memcached error: ", err)
					}
				}
			} else {
				log.Printf("Hotel with ID %s not found in database", hotelID)
			}
		} else {
			log.Printf("Memcached error = %s\n", err)
			// Instead of panicking, just log the error and continue
			log.Warn("Memcached error when fetching hotel profile: ", err)
		}
	}

	res.Hotels = hotels
	return res, nil
}

// Helper function to fetch a hotel from Cassandra by ID
func fetchHotelFromCassandra(session *gocql.Session, hotelID string) *pb.Hotel {
	var name, phoneNumber, description string
	var streetNumber, streetName, city, state, country, postalCode string
	var lat, lon float32

	// Query Cassandra for the hotel data
	err := session.Query(`
		SELECT name, phone_number, description, 
		       street_number, street_name, city, state, 
		       country, postal_code, lat, lon 
		FROM profile_keyspace.hotels WHERE id = ?`, 
		hotelID).Scan(
		&name, &phoneNumber, &description,
		&streetNumber, &streetName, &city, &state,
		&country, &postalCode, &lat, &lon)

	if err != nil {
		if err != gocql.ErrNotFound {
			log.Warn("Cassandra query error: ", err)
		}
		return nil
	}

	// Construct the hotel object with the retrieved data
	hotel := &pb.Hotel{
		Id:          hotelID,
		Name:        name,
		PhoneNumber: phoneNumber,
		Description: description,
		Address: &pb.Address{
			StreetNumber: streetNumber,
			StreetName:   streetName,
			City:         city,
			State:        state,
			Country:      country,
			PostalCode:   postalCode,
			Lat:          lat,
			Lon:          lon,
		},
	}

	return hotel
}