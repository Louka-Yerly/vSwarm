package main

import (
	"strings"
	"flag"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	log "github.com/sirupsen/logrus"
	tracing "github.com/vhive-serverless/vSwarm/utils/tracing/go"
)

var (
	zipkin        = flag.String("zipkin", "http://localhost:9411/api/v2/spans", "zipkin url")
	url           = flag.String("url", "0.0.0.0", "Address of the service")
	port          = flag.Int("port", 8083, "Port of the service")
	database_addr = flag.String("db_addr", "0.0.0.0:5432", "Address of the data base server")
	memc_addr     = flag.String("memcached_addr", "0.0.0.0:11211", "Address of the memcached server")
)

func main() {
	flag.Parse()

	// Setup tracing ---
	if tracing.IsTracingEnabled() {
		log.Printf("Start tracing on : %s\n", *zipkin)
		shutdown, err := tracing.InitBasicTracer(*zipkin, "Hotel app - profile function")
		if err != nil {
			log.Warn(err)
		}
		defer shutdown()
	}

	// Initialize database ---
	splited := strings.Split(*database_addr, ":")
	db := initializeDatabase(splited[0], splited[1], "postgres", "postgres", "hotel_profile")
	for db == nil {
		time.Sleep(5*time.Second)
		db = initializeDatabase(splited[0], splited[1], "postgres", "postgres", "hotel_profile")
	}
	defer db.Close()

	// Initialize Memcached ---
	memc_client := memcache.New(*memc_addr)
	memc_client.Timeout = time.Second * 2
	memc_client.MaxIdleConns = 512

	// Start the gRPC server ---
	srv := &Server{
		Port:       *port,
		IpAddr:     *url,
		DB:         db,
		MemcClient: memc_client,
	}
	log.Fatal(srv.Run())
}
