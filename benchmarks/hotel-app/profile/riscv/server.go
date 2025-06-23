package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	log "github.com/sirupsen/logrus"
	pb "github.com/vhive-serverless/vSwarm-proto/proto/hotel_reserv/profile"
	tracing "github.com/vhive-serverless/vSwarm/utils/tracing/go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// Server implements the profile service
type Server struct {
	pb.UnimplementedProfileServer

	Port       int
	IpAddr     string
	DB         *sql.DB
	MemcClient *memcache.Client
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

	// Process each hotel ID in the request
	for _, hotelID := range req.HotelIds {
		// First check memcached
		item, err := s.MemcClient.Get(hotelID)
		if err == nil {
			// Memcached hit
			log.Print("memcached hit")
			hotelProf := new(pb.Hotel)
			if err = json.Unmarshal(item.Value, hotelProf); err != nil {
				log.Warn(err)
			}
			hotels = append(hotels, hotelProf)
		} else if err == memcache.ErrCacheMiss {
			// Memcached miss, query from PostgreSQL
			log.Printf("memcached miss")
			hotelProf, err := getHotelByID(s.DB, hotelID)
			if err != nil {
				log.Printf("Error fetching hotel data from PostgreSQL: %v", err)
				continue
			}

			hotels = append(hotels, hotelProf)

			// Convert hotel data to JSON and store in memcached
			profJSON, err := json.Marshal(hotelProf)
			if err != nil {
				log.Warn(err)
				continue
			}

			// Write to memcached
			err = s.MemcClient.Set(&memcache.Item{Key: hotelID, Value: profJSON})
			if err != nil {
				log.Warn("Memcached error: ", err)
			}
		} else {
			log.Printf("Memcached error = %v\n", err)
			// Instead of panicking, we'll just log the error and continue
			log.Warn("Memcached error: ", err)
		}
	}

	res.Hotels = hotels
	return res, nil
}