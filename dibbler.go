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
	http.HandleFunc("/dibbler", requestToBidHandler)
	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err)
}

func requestToBidHandler(w http.ResponseWriter, r *http.Request) {
	bidRepository := repo.NewRedisBidRepository("localhost:6379")

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := string(buf.Bytes())
	eligibleCampaigns, err := service.GetSortedApplicableCampaigns(body, bidRepository)

	if err != nil {
		fmt.Println(err)
	}

	success, err := service.PlaceBids(eligibleCampaigns, bidRepository)

	if err != nil {
		fmt.Println(err)
	}

	if success {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(204)
	}
}
