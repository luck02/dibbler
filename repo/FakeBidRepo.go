package repo

import (
	"time"

	"github.com/luck02/dibbler/models"
)

type FakeBidRepository struct {
	CampaignCollection       []models.Campaign
	BidAttempts              int
	TestingExhaustedCampaign bool
	BidShouldBeSuccesfull    bool
}

func (r *FakeBidRepository) PlaceBid(campaign models.Campaign) (models.Campaign, bool, error) {
	r.BidAttempts += 1
	if r.TestingExhaustedCampaign {
		campaign.RemainingBudget = 0
	}

	return campaign, r.BidShouldBeSuccesfull, nil
}

func (r *FakeBidRepository) GetCampaigns() ([]models.Campaign, error) {
	return r.CampaignCollection, nil
}
func (r *FakeBidRepository) GetCampaignsCached(currentTime time.Time) ([]models.Campaign, error) {
	return r.CampaignCollection, nil
}
