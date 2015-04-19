package repo

import (
	"github.com/luck02/dibbler/models"
)

type BidRepository interface {
	PlaceBid(models.Campaign) (models.Campaign, bool)
	GetCampaigns() []models.Campaign
}

type FakeBidRepository struct {
	CampaignCollection []models.Campaign
	BidAttempts        int
}

func (r FakeBidRepository) PlaceBid(campaign models.Campaign) (models.Campaign, bool) {
	return campaign, true
}

func (r FakeBidRepository) GetCampaigns() []models.Campaign {
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
