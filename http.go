package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	httpEndpoint = "0.0.0.0:3000"
)

func runHttp() {
	r := mux.NewRouter()

	r.HandleFunc("/api/weather", httpGetWeather).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./www")))

	go func() {
		fmt.Fprintf(os.Stdout, "http: listening on \"%s\"\n", httpEndpoint)
		http.ListenAndServe(httpEndpoint, r)
	}()
}

func httpGetWeather(w http.ResponseWriter, r *http.Request) {
	weather := weatherGet()

	bytes, err := json.Marshal(weather)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(bytes)
	if err != nil {
		panic(err)
	}
}
