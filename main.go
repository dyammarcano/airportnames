package main

import (
	"embed"
	_ "embed"
	"encoding/csv"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jszwec/csvutil"
	"net/http"
	"strings"
)

//go:embed airports.csv
var airportsFile embed.FS

var airports []Airport

type Airport struct {
	ID               int    `json:"id" csv:"id"`
	Ident            string `json:"ident" csv:"ident"`
	Type             string `json:"type" csv:"type"`
	Name             string `json:"name" csv:"name"`
	LatitudeDeg      string `json:"latitude_deg" csv:"latitude_deg"`
	LongitudeDeg     string `json:"longitude_deg" csv:"longitude_deg"`
	ElevationFt      string `json:"elevation_ft" csv:"elevation_ft"`
	Continent        string `json:"continent" csv:"continent"`
	IsoCountry       string `json:"iso_country" csv:"iso_country"`
	IsoRegion        string `json:"iso_region" csv:"iso_region"`
	Municipality     string `json:"municipality" csv:"municipality"`
	ScheduledService string `json:"scheduled_service" csv:"scheduled_service"`
	GpsCode          string `json:"gps_code" csv:"gps_code"`
	IataCode         string `json:"iata_code" csv:"iata_code"`
	LocalCode        string `json:"local_code" csv:"local_code"`
	HomeLink         string `json:"home_link" csv:"home_link"`
	WikipediaLink    string `json:"wikipedia_link" csv:"wikipedia_link"`
	Keywords         string `json:"keywords" csv:"keywords"`
}

func init() {
	loadAirports()
}

func loadAirports() {
	data, err := airportsFile.Open("airports.csv")
	if err != nil {
		panic(err)
	}

	csvReader := csv.NewReader(data)

	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		panic(err)
	}

	for {
		var airport Airport
		if err = dec.Decode(&airport); err != nil {
			break
		}
		airports = append(airports, airport)
	}
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Get("/airportName", func(w http.ResponseWriter, r *http.Request) {
		iataCode := r.URL.Query().Get("iataCode")
		if iataCode == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No IATA code provided."))
			return
		}

		iataCode = strings.ToUpper(iataCode)

		for _, airport := range airports {
			if airport.IataCode == iataCode {
				w.WriteHeader(http.StatusOK)
				
				return
			}
		}
	})

	http.ListenAndServe(":8080", r)
}
