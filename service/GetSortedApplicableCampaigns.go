package service

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/jmoiron/jsonq"
	"github.com/luck02/dibbler/models"
	"github.com/luck02/dibbler/repo"
)

// GetApplicableCampaigns Returns an ordered list of campaigns applicable to the given OTB
func GetSortedApplicableCampaigns(otbJSON string, bidRepository repo.BidRepository) ([]models.Campaign, error) {

	campaigns := bidRepository.GetCampaigns()
	var applicableCampaigns []models.Campaign
	requestToBidData := map[string]interface{}{}
	decoder := json.NewDecoder(strings.NewReader(otbJSON))
	decoder.Decode(&requestToBidData)

	for _, campaign := range campaigns {
		if campaignApplicable(requestToBidData, campaign) {
			applicableCampaigns = append(applicableCampaigns, campaign)
			fmt.Println("Applicable Campaign:")
			fmt.Println(campaign)
		}
	}

	sort.Sort(models.SortedCampaigns(applicableCampaigns))
	return applicableCampaigns, nil
}

func campaignApplicable(otbJSON map[string]interface{}, campaign models.Campaign) bool {

	requestToBidQuery := jsonq.NewQuery(otbJSON)

	switch target := campaign.Targeting.(type) {
	case models.PlacementTarget:
		if appName, err := requestToBidQuery.String("app", "name"); err != nil {
			fmt.Println(err)
		} else {
			return appName == target.AppName
		}
	case models.AdTarget:
		width, err := requestToBidQuery.Int("imp", "0", "banner", "w")
		if err != nil {
			fmt.Println(err)
		}

		height, err := requestToBidQuery.Int("imp", "0", "banner", "h")
		if err != nil {
			fmt.Println(err)
		}

		if width == target.Width && height == target.Height {
			return true
		}
	case models.CountryTarget:
		if country, err := requestToBidQuery.String("device", "geo", "country"); err != nil {
			fmt.Println(err)
		} else {
			return country == target.Country
		}
	case models.OSTarget:
		if osName, err := requestToBidQuery.String("device", "os"); err != nil {
			fmt.Println(err)
		} else {
			return osName == target.OsType
		}
	}
	return false
}
