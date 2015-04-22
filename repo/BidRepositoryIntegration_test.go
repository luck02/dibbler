package repo

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/luck02/dibbler/fixtures"
	"github.com/luck02/dibbler/models"
)

func dontRunTestRedisBasicFunctions(t *testing.T) {
	pool := newPool("localhost:6379")

	conn := pool.Get()
	defer conn.Close()

	n, err := conn.Do("SET", "derp", "DERPS")

	fmt.Println(n)
	fmt.Println(err)

	b, err := json.Marshal(fixtures.CampaignTests[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

	sampleCampaign := models.Campaign{}
	err = json.Unmarshal(b, &sampleCampaign)
	fmt.Println(err)
	fmt.Printf("%+v\n", sampleCampaign)
}

func TestICanSaveAndLoadACampaign(t *testing.T) {
	bidRepository := NewRedisBidRepository("localhost:6379")
	err := bidRepository.saveCampaign(fixtures.CampaignTests[0])

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
	bidRepository := NewRedisBidRepository("localhost:6379")

	for _, value := range fixtures.CampaignTests {
		err := bidRepository.saveCampaign(value)
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
	bidRepository := NewRedisBidRepository("localhost:6379")

	err := bidRepository.saveCampaign(fixtures.CampaignTests[0])
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
	bidRepository := NewRedisBidRepository("localhost:6379")

	err := bidRepository.saveCampaign(campaign)
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

var (
	bidRepository = NewRedisBidRepository("localhost:6379")
)

func BenchmarkPlaceBid(b *testing.B) {
	//	campaign := fixtures.CampaignTests[0]
	err := bidRepository.saveCampaign(fixtures.CampaignTests[0])
	if err != nil {
		b.Error(err)
	}

	_, _, err = bidRepository.PlaceBid(fixtures.CampaignTests[0])
	if err != nil {
		b.Error(err)
	}

}
