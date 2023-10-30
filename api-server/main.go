package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	openapispec := flag.String("openapi", "", "openapispec")
	datafile := flag.String("data", "", "file to put data in")
	latency := flag.Duration("latency", 0, "latency to add")
	errorrate := flag.Int("errorrate", 0, "latency to add")
	flag.Parse()

	var openapi []byte

	if openapispec != nil && *openapispec != "" {
		var err error
		openapi, err = os.ReadFile(*openapispec)
		if err != nil {
			log.Fatal(err)
		}
	}

	var data []byte
	if datafile != nil && *datafile != "" {
		var err error
		data, err = os.ReadFile(*datafile)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Fatal(http.ListenAndServe(":3000", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/openapi.yaml" {
			_, _ = rw.Write(openapi)
			return
		}

		if latency != nil && *latency > 0 {
			time.Sleep(*latency)
		}

		status := req.URL.Query().Get("status")
		if status == "" {
			if errorrate != nil && *errorrate > 0 {
				if rand.Int()%100 < *errorrate {
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			switch req.Method {
			case http.MethodGet, http.MethodPut:
				rw.WriteHeader(http.StatusOK)
				_, _ = rw.Write(data)
			case http.MethodPost:
				rw.WriteHeader(http.StatusCreated)
				_, _ = rw.Write([]byte(`{"id":4}`))
			case http.MethodDelete:
				rw.WriteHeader(http.StatusNoContent)
			}

			return
		}

		atoi, err := strconv.Atoi(status)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(atoi)
	})))
}
