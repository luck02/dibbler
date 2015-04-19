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
func GetApplicableCampaigns(otbJSON string, bidRepository repo.BidRepository) ([]models.Campaign, error) {

	campaigns := bidRepository.GetCampaigns()
	var applicableCampaigns []models.Campaign
	otbData := map[string]interface{}{}
	decoder := json.NewDecoder(strings.NewReader(otbJSON))
	decoder.Decode(&otbData)

	for _, campaign := range campaigns {
		if campaignApplicable(otbData, campaign) {
			applicableCampaigns = append(applicableCampaigns, campaign)
			fmt.Println("Applicable Campaign:")
			fmt.Println(campaign)
		}
	}

	sort.Sort(models.SortedCampaigns(applicableCampaigns))
	return applicableCampaigns, nil
}

func campaignApplicable(otbJSON map[string]interface{}, campaign models.Campaign) bool {

	otbQuery := jsonq.NewQuery(otbJSON)

	switch target := campaign.Targeting.(type) {
	case models.PlacementTarget:
		if appName, err := otbQuery.String("app", "name"); err != nil {
			fmt.Println(err)
		} else {
			return appName == target.AppName
		}
	case models.AdTarget:
		width, err := otbQuery.Int("imp", "0", "banner", "w")
		if err != nil {
			fmt.Println(err)
		}

		height, err := otbQuery.Int("imp", "0", "banner", "h")
		if err != nil {
			fmt.Println(err)
		}

		if width == target.Width && height == target.Height {
			return true
		}
	case models.CountryTarget:
		if country, err := otbQuery.String("device", "geo", "country"); err != nil {
			fmt.Println(err)
		} else {
			return country == target.Country
		}
	case models.OSTarget:
		if osName, err := otbQuery.String("device", "os"); err != nil {
			fmt.Println(err)
		} else {
			return osName == target.OsType
		}
	}

	return false

}
