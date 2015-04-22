package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/luck02/dibbler/repo"
	"github.com/luck02/dibbler/service"
)

func main() {
	fmt.Println("start")
	http.HandleFunc("/", requestToBidHandler)
	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err)
}

func requestToBidHandler(w http.ResponseWriter, r *http.Request) {
	// Create GUID
	bidRepository := repo.NewRedisBidRepository("localhost:6379")

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := string(buf.Bytes())

	// log (guid, timestamp, body)
	eligibleCampaigns, err := service.GetSortedApplicableCampaigns(body, bidRepository)
	// log (guid, timestamp, eligibleCampaigns)
	if err != nil {
		fmt.Println(err)
	}
	// if no eligible campatins no need to call place bids
	success, err := service.PlaceBids(eligibleCampaigns, bidRepository)
	//log (guid, timestamp, bid placed)
	if err != nil {
		fmt.Println(err)
	}

	if success {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(204)
	}
}
