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

func TestPlaceBidSuccess(t *testing.T) {
	//redisBidRepo := RedisBidRepository{}

}
