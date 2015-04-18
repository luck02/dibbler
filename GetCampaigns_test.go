package dibbler

import (
	"encoding/json"
	"sort"
	"strings"
	"testing"
)

func getOtbQueryObject(otbString string) map[string]interface{} {
	otbData := map[string]interface{}{}
	decoder := json.NewDecoder(strings.NewReader(otbString))
	decoder.Decode(&otbData)
	return otbData
}

func TestCampaignIsApplicablePlacementTarget(t *testing.T) {
	otbData := getOtbQueryObject(otbPlacement)
	expected := campaignApplicable(otbData, CampaignTests[0])
	if !expected {
		t.Error("CampaignTest for placement should be true")
	}

	expected = campaignApplicable(otbData, CampaignTests[1])
	if expected {
		t.Error("CampaignTest for placement should be false")
	}
}

func TestCampaignIsApplicableAdTarget(t *testing.T) {
	otbData := getOtbQueryObject(otbAd)
	expected := campaignApplicable(otbData, CampaignTests[1])
	if !expected {
		t.Error("CampaignTest for AdTarget should be true")
	}

	expected = campaignApplicable(otbData, CampaignTests[0])
	if expected {
		t.Error("CampaignTest for AdTarget should be false")
	}
}

func TestCampaignIsApplicableCountryTarget(t *testing.T) {
	otbData := getOtbQueryObject(otbAd)
	expected := campaignApplicable(otbData, CampaignTests[2])
	if !expected {
		t.Error("CampaignTest for Country should be true")
	}

	expected = campaignApplicable(otbData, CampaignTests[4])
	if expected {
		t.Error("CampaignTest for Country should be false")
	}
}

func TestCampaignIsApplicableOsTarget(t *testing.T) {
	otbData := getOtbQueryObject(otbAd)
	expected := campaignApplicable(otbData, CampaignTests[3])
	if !expected {
		t.Error("CampaignTest for Os should be true")
	}

	expected = campaignApplicable(otbData, CampaignTests[5])
	if expected {
		t.Error("CampaignTest for Os should be false")
	}
}

func TestGetCampaigns(t *testing.T) {
	fakeBidRepository := FakeBidRepository{CampaignCollection: CampaignTests}
	sortedList, err := GetApplicableCampaigns(otbPlacement, fakeBidRepository)

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
	CampaignTests[0].BidCpm = 0.31
	fakeBidRepository := FakeBidRepository{CampaignCollection: CampaignTests}
	sortedList, err := GetApplicableCampaigns(otbPlacement, fakeBidRepository)
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
	for i := range CampaignTests {
		CampaignTests[i].BidCpm = bidCpm
		bidCpm -= float32(5)
	}
	sort.Sort(SortedCampaigns(CampaignTests))
	current := float32(99)
	for _, campaign := range CampaignTests {
		if campaign.BidCpm > current {
			t.Error("campaignList out of order")
		}
		current = campaign.BidCpm
	}
}
