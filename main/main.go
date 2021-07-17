package main

import (
	"encoding/json"
	"github.com/gen2brain/beeep"
	"log"
	"net/http"
	"time"
)

type ServerStatus struct {
	StatusMaintenance bool `json:"status_maintenance"`
	StatusLogin       bool `json:"status_login"`
	StatusBattles     int  `json:"status_battles"`
	StatusGraphql     bool `json:"status_graphql"`
	StatusCloudflare  bool `json:"status_cloudflare"`
}

func main() {
	// Check the status every 30 seconds
	for t := range time.Tick(30 * time.Second) {
		go checkStatus(t)
	}
}

func checkStatus(t time.Time) {
	log.Printf("Checking status [%v]\n", t)

	// Send the request to the axie.zone endpoint
	res, err := http.Get("https://axie.zone:3000/server_status")
	if err != nil {
		log.Println("Endpoint not responding")
		return
	}

	// Decode the response
	var s ServerStatus
	err = json.NewDecoder(res.Body).Decode(&s)
	if err != nil {
		log.Println("Error while parsing response body")
		return
	}

	notifyStatus(s)
}

// Used to keep track of how many times that the server has been working
var battleStatus = 0

func notifyStatus(s ServerStatus) {
	// Update the battle status count based on the response
	if s.StatusBattles == 0 {
		battleStatus = 0
		log.Println("Battle service status: down / not responding")
		return
	}
	battleStatus++

	// If battle status has been ok for 2 consecutive status check, send a notification message
	log.Println("Battle service status: OK")
	if battleStatus == 2 {
		_ = beeep.Notify("Battle service running", "You may try playing the game now", "assets/logo.png")
	}
}
