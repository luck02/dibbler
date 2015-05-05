package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"code.google.com/p/go-uuid/uuid"

	"github.com/Sirupsen/logrus"
	"github.com/luck02/dibbler/repo"
	"github.com/luck02/dibbler/service"
)

var (
	BidRepository = repo.NewRedisBidRepository("localhost:6379", 5)
)

func main() {
	var port = flag.Int("port", 8080, "Port to run webserver on")
	var logFile = flag.String("logfile", "log.txt", "Logfile path and name")
	flag.Parse()

	initLogger(*logFile)
	logrus.Info(fmt.Sprintf("Server starting on port:%d", *port))
	http.HandleFunc("/", requestToBidHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", *port), nil))
}

func requestToBidHandler(w http.ResponseWriter, r *http.Request) {
	correlationId := uuid.NewUUID()
	//bidRepository := repo.NewRedisBidRepository("localhost:6379")

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := string(buf.Bytes())
	logrus.WithFields(logrus.Fields{"Event": "RequestToBidReceived", "correlationId": correlationId, "content": body}).Info("Received Request")

	eligibleCampaigns, err := service.GetSortedApplicableCampaigns(body, BidRepository)
	if err != nil {
		logrus.Error(err)
	}

	var success bool
	if len(eligibleCampaigns) <= 0 {
		logrus.WithFields(logrus.Fields{"Event": "GetCampaigns", "correlationId": correlationId, "content": "No Eligible Campaigns"}).Info("Campaigns")
	} else {
		logrus.WithFields(logrus.Fields{"Event": "GotCampaigns", "correlationId": correlationId, "content": eligibleCampaigns}).Info("SortedCampaigns")
		success, err = service.PlaceBids(eligibleCampaigns, BidRepository)
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

func initLogger(logFileAndPath string) {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	file, err := os.Create(logFileAndPath)
	if err != nil {
		panic(err)
	}
	logrus.SetOutput(file)
}
