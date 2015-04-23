package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"code.google.com/p/go-uuid/uuid"

	"github.com/Sirupsen/logrus"
	"github.com/luck02/dibbler/repo"
	"github.com/luck02/dibbler/service"
)

func main() {
	fmt.Println("start")
	http.HandleFunc("/", requestToBidHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func requestToBidHandler(w http.ResponseWriter, r *http.Request) {
	correlationId := uuid.NewUUID()
	bidRepository := repo.NewRedisBidRepository("localhost:6379")

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := string(buf.Bytes())
	logrus.WithFields(logrus.Fields{"Event": "RequestToBidReceived", "correlationId": correlationId, "content": body}).Info("Received Request")

	eligibleCampaigns, err := service.GetSortedApplicableCampaigns(body, bidRepository)
	if err != nil {
		logrus.Error(err)
	}

	var success bool
	if len(eligibleCampaigns) <= 0 {
		logrus.WithFields(logrus.Fields{"Event": "GetCampaigns", "correlationId": correlationId, "content": "No Eligible Campaigns"}).Info("Campaigns")
	} else {
		logrus.WithFields(logrus.Fields{"Event": "GotCampaigns", "correlationId": correlationId, "content": eligibleCampaigns}).Info("SortedCampaigns")
		success, err = service.PlaceBids(eligibleCampaigns, bidRepository)
		if err != nil {
			fmt.Println(err)
		}

	}

	if success {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(204)
	}
}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	file, err := os.Create("./logs.txt")
	if err != nil {
		panic(err)
	}
	logrus.SetOutput(file)
}
