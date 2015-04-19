package repo

import (
	"github.com/luck02/dibbler/models"
)

type BidRepository interface {
	PlaceBid(models.Campaign) (models.Campaign, bool)
	GetCampaigns() []models.Campaign
}

type FakeBidRepository struct {
	CampaignCollection       []models.Campaign
	BidAttempts              int
	TestingExhaustedCampaign bool
}

func (r *FakeBidRepository) PlaceBid(campaign models.Campaign) (models.Campaign, bool) {
	r.BidAttempts += 1
	if r.TestingExhaustedCampaign {
		campaign.RemainingBudget = 0
	}
	return campaign, false
}

func (r *FakeBidRepository) GetCampaigns() []models.Campaign {
	return r.CampaignCollection
}

type RedisBidRepository struct {
}

func (r RedisBidRepository) PlaceBid(campaign models.Campaign) (models.Campaign, bool) {

	return models.Campaign{}, true
}

func (r RedisBidRepository) GetCampaigns() []models.Campaign {
	return nil
}
