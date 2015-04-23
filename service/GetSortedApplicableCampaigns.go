package service

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/jsonq"
	"github.com/luck02/dibbler/models"
	"github.com/luck02/dibbler/repo"
)

// GetApplicableCampaigns Returns an ordered list of campaigns applicable to the given requestToBid
func GetSortedApplicableCampaigns(requestToBidJSON string, bidRepository repo.BidRepository) ([]models.Campaign, error) {

	campaigns, err := bidRepository.GetCampaigns()
	if err != nil {
		return nil, err
	}
	var applicableCampaigns []models.Campaign
	requestToBidData := map[string]interface{}{}
	decoder := json.NewDecoder(strings.NewReader(requestToBidJSON))
	decoder.Decode(&requestToBidData)

	for _, campaign := range campaigns {
		if campaignApplicable(requestToBidData, campaign) {
			applicableCampaigns = append(applicableCampaigns, campaign)
		}
	}

	sort.Sort(models.SortedCampaigns(applicableCampaigns))
	return applicableCampaigns, nil
}

func campaignApplicable(requestToBidJSON map[string]interface{}, campaign models.Campaign) bool {

	requestToBidQuery := jsonq.NewQuery(requestToBidJSON)

	switch target := campaign.Targeting.(type) {
	case models.PlacementTarget:
		if appName, err := requestToBidQuery.String("app", "name"); err != nil {
			logrus.Info(err)
		} else {
			return appName == target.AppName
		}
	case models.AdTarget:
		arrayOfImpressions, err := requestToBidQuery.ArrayOfObjects("imp")
		if err != nil {
			logrus.Info(err)
		}
		for _, impression := range arrayOfImpressions {
			impressionQuery := jsonq.NewQuery(impression)
			width, err := impressionQuery.Int("banner", "w")
			if err != nil {
				logrus.Info(err)
			}

			height, err := impressionQuery.Int("banner", "h")
			if err != nil {
				logrus.Info(err)
			}

			if width == target.Width && height == target.Height {
				return true
			}
		}
	case models.CountryTarget:
		if country, err := requestToBidQuery.String("device", "geo", "country"); err != nil {
			logrus.Info(err)
		} else {
			return country == target.Country
		}
	case models.OSTarget:
		if osName, err := requestToBidQuery.String("device", "os"); err != nil {
			logrus.Info(err)
		} else {
			return osName == target.OsType
		}
	}
	return false
}
