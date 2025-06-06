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
	"flag"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	tracing "github.com/vhive-serverless/vSwarm/utils/tracing/go"
)

var (
	zipkin        = flag.String("zipkin", "http://localhost:9411/api/v2/spans", "zipkin url")
	url           = flag.String("url", "0.0.0.0", "Address of the service")
	port          = flag.Int("port", 8083, "Port of the service")
	database_addr = flag.String("db_addr", "0.0.0.0:5432", "Address of the data base server")
)

func main() {

	flag.Parse()

	// Setup tracing ---
	if tracing.IsTracingEnabled() {
		log.Printf("Start tracing on : %s\n", *zipkin)
		shutdown, err := tracing.InitBasicTracer(*zipkin, "Hotel app - recommendation function")
		if err != nil {
			log.Warn(err)
		}
		defer shutdown()
	}

	// Initialize database ---
	splited := strings.Split(*database_addr, ":")
	db := initializeDatabase(splited[0], splited[1], "postgres", "postgres", "hotel_recommendation")
	for db == nil {
		time.Sleep(5 * time.Second)
		db = initializeDatabase(splited[0], splited[1], "postgres", "postgres", "hotel_recommendation")
	}
	defer db.Close()

	// Start the gRPC server ---
	srv := &Server{
		Port:   *port,
		IpAddr: *url,
		DB:     db,
	}
	log.Fatal(srv.Run())
}
