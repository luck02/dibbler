package repo

import (
	"reflect"
	"testing"

	"github.com/luck02/dibbler/fixtures"
)

var (
	redisExists           bool
	redisConnectionString string = "localhost:6379"
	bidRepository         *RedisBidRepository
)

func init() {
	bidRepository = NewRedisBidRepository(redisConnectionString)
}

func TestICanSaveAndLoadACampaign(t *testing.T) {
	err := bidRepository.SaveCampaign(fixtures.CampaignTests[0])

	if err != nil {
		t.Errorf("Error should be nil, was: %v", err)
		return
	}

	savedCampaign, err := bidRepository.getCampaign(fixtures.CampaignTests[0].ID)
	if err != nil {
		t.Errorf("Error should be nil, was: %v", err)

	}
	if !reflect.DeepEqual(fixtures.CampaignTests[0], savedCampaign) {
		t.Errorf("campaigns should be equal\n %+v \n %+v", fixtures.CampaignTests[0], savedCampaign)
	}
}

func TestICanSaveFixturesAndLoadThem(t *testing.T) {
	for _, value := range fixtures.CampaignTests {
		err := bidRepository.SaveCampaign(value)
		if err != nil {
			t.Error(err)
		}
	}

	campaigns, err := bidRepository.GetCampaigns()
	if err != nil {
		t.Error(err)
	}

	if len(campaigns) != len(fixtures.CampaignTests) {
		t.Errorf("Expected %+v\n Got %+v\n", fixtures.CampaignTests, campaigns)
	}
}

func TestPlaceBidSuccess(t *testing.T) {
	err := bidRepository.SaveCampaign(fixtures.CampaignTests[0])
	if err != nil {
		t.Errorf("Failed to save campaign")
	}

	resultCampaign, success, err := bidRepository.PlaceBid(fixtures.CampaignTests[0])
	if !success || err != nil {
		t.Errorf("Failed to place bid %v, %v", success, err)
	}
	expected := fixtures.CampaignTests[0].RemainingBudget - fixtures.CampaignTests[0].BidCpm/1000

	if resultCampaign.RemainingBudget != expected {
		t.Errorf("Got: %d : Expected %d", resultCampaign.RemainingBudget, expected)
	}
}

func TestPlaceBidExhaustion(t *testing.T) {
	campaign := fixtures.CampaignTests[0]

	campaign.RemainingBudget = 0.25
	totalRuns := int(campaign.RemainingBudget / (campaign.BidCpm / 1000))

	err := bidRepository.SaveCampaign(campaign)
	if err != nil {
		t.Error(err)
	}
	i := 0
	success := true
	for success {
		campaign, success, err = bidRepository.PlaceBid(campaign)
		if !success {
			break
		}
		if err != nil {
			t.Error(err)
		}
		i++
	}

	if campaign.RemainingBudget != 0 {
		t.Errorf("RemainingBudget Expected: 0\n actual: %d", campaign.RemainingBudget)
	}

	if totalRuns != i {
		t.Errorf("Expected total debits: %d   actual total Debits %d", totalRuns, i)
	}

}

func BenchmarkPlaceBid(b *testing.B) {
	err := bidRepository.SaveCampaign(fixtures.CampaignTests[0])
	if err != nil {
		b.Error(err)
	}

	_, _, err = bidRepository.PlaceBid(fixtures.CampaignTests[0])
	if err != nil {
		b.Error(err)
	}

}
