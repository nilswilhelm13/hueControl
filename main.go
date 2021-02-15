package main

import (
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
)

var client *http.Client
var db *bolt.DB
var scheduler *cron.Cron

func main() {
	db = openDb()

	// initialize http client
	client = &http.Client{}

	scheduler = cron.New()
	UpdateSchedule(db, scheduler)
	scheduler.Start()

	router := mux.NewRouter()
	router.HandleFunc("/", scheduleHandler)
	log.Fatal(http.ListenAndServe(":9000", router))

}
