package dibbler

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/jmoiron/jsonq"
	"github.com/luck02/dibbler/dibbler/repo"
)

// GetApplicableCampaigns Returns an ordered list of campaigns applicable to the given OTB
func GetApplicableCampaigns(otbJSON string, bidRepository repo.BidRepository) ([]Campaign, error) {

	campaigns := bidRepository.GetCampaigns()
	var applicableCampaigns []Campaign
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

	sort.Sort(SortedCampaigns(applicableCampaigns))
	return applicableCampaigns, nil
}

func campaignApplicable(otbJSON map[string]interface{}, campaign Campaign) bool {

	otbQuery := jsonq.NewQuery(otbJSON)

	switch target := campaign.Targeting.(type) {
	case PlacementTarget:
		if appName, err := otbQuery.String("app", "name"); err != nil {
			fmt.Println(err)
		} else {
			return appName == target.AppName
		}
	case AdTarget:
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
	case CountryTarget:
		if country, err := otbQuery.String("device", "geo", "country"); err != nil {
			fmt.Println(err)
		} else {
			return country == target.Country
		}
	case OSTarget:
		if osName, err := otbQuery.String("device", "os"); err != nil {
			fmt.Println(err)
		} else {
			return osName == target.OsType
		}
	}

	return false

}
