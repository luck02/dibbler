package dibbler

import (
	"encoding/json"
	"strings"
	"testing"
)

var otbPlacement = `
{
    "app": {
        "name": "Words With Friends 2 iPad"
    },
    "imp": [
        {
            "banner": {
                "h": 250,
                "w": 300
            },
            "id": "1"
        }
    ],
    "device": {
        "os": "iOS",
        "geo": {
            "city": "Irwin",
            "region": "PA",
            "zip": "15642",
            "country": "USA"
        }
    }
}
`

var CampaignTests = []Campaign{
	{Id: 100101, BidCpm: 0.32, DailyBudget: 35.50, RemainingBudget: 35.50, Targeting: PlacementTarget{AppName: "Words With Friends 2 iPad"}},
	{Id: 100102, BidCpm: 0.04, DailyBudget: 5.25, RemainingBudget: 5.25, Targeting: AdTarget{Height: 728, Width: 1024}},
	{Id: 100103, BidCpm: 0.32, DailyBudget: 15.00, RemainingBudget: 15.00, Targeting: CountryTarget{Country: "USA"}},
	{Id: 100104, BidCpm: 0.15, DailyBudget: 22.00, RemainingBudget: 22.00, Targeting: OSTarget{OsType: "Android"}},
	{Id: 100105, BidCpm: 0.02, DailyBudget: 2.25, RemainingBudget: 2.25, Targeting: CountryTarget{Country: "MEX"}},
}

func TestCampaignIsApplicablePlacementTarget(t *testing.T) {
	otbData := map[string]interface{}{}
	decoder := json.NewDecoder(strings.NewReader(otbPlacement))
	decoder.Decode(&otbData)

	expected := campaignApplicable(otbData, CampaignTests[0])
	if !expected {
		t.Error("CampaignTest for placement should be true")
	}

	expected = campaignApplicable(otbData, CampaignTests[1])
	if expected {
		t.Error("CampaignTest for placement should be false")
	}
}

func TestGetCampaigns(t *testing.T) {
	sortedList, err := GetApplicableCampaigns(otbPlacement, CampaignTests)

	if err != nil {
		t.Error("Error returned", sortedList)
	}

	if sortedList[0].Id != 100101 {
		t.Error("incorrect order [0] returned", sortedList)
	}

	if sortedList[1].Id != 100103 {
		t.Error("incorrect order [1] returned", sortedList)
	}
}
func TestGetCampaignsReordered(t *testing.T) {
	CampaignTests[0].BidCpm = 0.31
	sortedList, err := GetApplicableCampaigns(otbPlacement, CampaignTests)

	if err != nil {
		t.Error("Error returned", sortedList)
	}

	if sortedList[0].Id != 100103 {
		t.Error("incorrect order [0] returned", sortedList)
	}

	if sortedList[1].Id != 100101 {
		t.Error("incorrect order [1] returned", sortedList)
	}

}
