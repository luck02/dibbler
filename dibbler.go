package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/luck02/dibbler/dibbler"
)

func main() {
	fmt.Println("start")
	http.HandleFunc("/dibbler", otbHandler)
	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err)
}

func otbHandler(w http.ResponseWriter, r *http.Request) {
	redisBidRepository := dibbler.FakeBidRepository{CampaignCollection: dibbler.CampaignTests}
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := string(buf.Bytes())
	eligibleCampaigns, err := dibbler.GetApplicableCampaigns(body, redisBidRepository)

	if err != nil {
		fmt.Println(err)
	}

	success, err := dibbler.PlaceBids(eligibleCampaigns, redisBidRepository)

	if err != nil {
		fmt.Println(err)
	}

	if success {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(204)
	}
}
