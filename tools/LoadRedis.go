package main

import (
	"fmt"

	"github.com/luck02/dibbler/fixtures"
	"github.com/luck02/dibbler/repo"
)

func main() {

	repo := repo.NewRedisBidRepository("localhost:6379")

	for _, campaign := range fixtures.CampaignTests {
		err := repo.SaveCampaign(campaign)
		fmt.Println(err)
	}
}
