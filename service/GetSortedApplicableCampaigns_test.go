package service

import (
	"encoding/json"
	"sort"
	"strings"
	"testing"

	"github.com/luck02/dibbler/fixtures"
	"github.com/luck02/dibbler/models"
	"github.com/luck02/dibbler/repo"
)

func getRequestToBidQueryObject(requestToBidJson string) map[string]interface{} {
	requestToBidData := map[string]interface{}{}
	decoder := json.NewDecoder(strings.NewReader(requestToBidJson))
	decoder.Decode(&requestToBidData)
	return requestToBidData
}

func TestCampaignIsApplicablePlacementTarget(t *testing.T) {
	requestToBidData := getRequestToBidQueryObject(fixtures.RequestToBidPlacement)
	expected := campaignApplicable(requestToBidData, fixtures.CampaignTests[0])
	if !expected {
		t.Error("CampaignTest for placement should be true")
	}

	expected = campaignApplicable(requestToBidData, fixtures.CampaignTests[1])
	if expected {
		t.Error("CampaignTest for placement should be false")
	}
}

func TestCampaignIsApplicableAdTarget(t *testing.T) {
	requestToBidData := getRequestToBidQueryObject(fixtures.RequestToBidAdSize)
	expected := campaignApplicable(requestToBidData, fixtures.CampaignTests[1])
	if !expected {
		t.Error("CampaignTest for AdTarget should be true")
	}

	expected = campaignApplicable(requestToBidData, fixtures.CampaignTests[0])
	if expected {
		t.Error("CampaignTest for AdTarget should be false")
	}
}

func TestCampaignIsApplicableCountryTarget(t *testing.T) {
	requestToBidData := getRequestToBidQueryObject(fixtures.RequestToBidAdSize)
	expected := campaignApplicable(requestToBidData, fixtures.CampaignTests[2])
	if !expected {
		t.Error("CampaignTest for Country should be true")
	}

	expected = campaignApplicable(requestToBidData, fixtures.CampaignTests[4])
	if expected {
		t.Error("CampaignTest for Country should be false")
	}
}

func TestCampaignIsApplicableOsTarget(t *testing.T) {
	requestToBidData := getRequestToBidQueryObject(fixtures.RequestToBidAdSize)
	expected := campaignApplicable(requestToBidData, fixtures.CampaignTests[3])
	if !expected {
		t.Error("CampaignTest for Os should be true")
	}

	expected = campaignApplicable(requestToBidData, fixtures.CampaignTests[5])
	if expected {
		t.Error("CampaignTest for Os should be false")
	}
}

func TestGetCampaigns(t *testing.T) {
	fakeBidRepository := new(repo.FakeBidRepository)
	fakeBidRepository.CampaignCollection = fixtures.CampaignTests
	sortedList, err := GetSortedApplicableCampaigns(fixtures.RequestToBidPlacement, fakeBidRepository)

	if err != nil {
		t.Error("Error returned", sortedList)
	}

	if sortedList[0].ID != 100101 {
		t.Error("incorrect order [0] returned", sortedList[0])
	}

	if sortedList[1].ID != 100103 {
		t.Error("incorrect order [1] returned", sortedList[1])
	}
}
func TestGetCampaignsReordered(t *testing.T) {
	fixtures.CampaignTests[0].BidCpm = 0.31
	fakeBidRepository := new(repo.FakeBidRepository)
	fakeBidRepository.CampaignCollection = fixtures.CampaignTests
	sortedList, err := GetSortedApplicableCampaigns(fixtures.RequestToBidPlacement, fakeBidRepository)
	if err != nil {
		t.Error("Error returned", sortedList)
	}

	if sortedList[0].ID != 100103 {
		t.Error("incorrect order [0] returned", sortedList[0])
	}

	if sortedList[1].ID != 100101 {
		t.Error("incorrect order [1] returned", sortedList[1])
	}
}

func TestCampaignSorter(t *testing.T) {
	bidCpm := float32(99)
	for i := range fixtures.CampaignTests {
		fixtures.CampaignTests[i].BidCpm = bidCpm
		bidCpm -= float32(5)
	}
	sort.Sort(models.SortedCampaigns(fixtures.CampaignTests))
	current := float32(99)
	for _, campaign := range fixtures.CampaignTests {
		if campaign.BidCpm > current {
			t.Error("campaignList out of order")
		}
		current = campaign.BidCpm
	}
}
