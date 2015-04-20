package repo

import (
	"github.com/luck02/dibbler/models"
)

type BidRepository interface {
	PlaceBid(models.Campaign) (models.Campaign, bool)
	GetCampaigns() []models.Campaign
}

type RedisBidRepository struct {
	URL      string
	Port     int
	Password string
}

func (r *RedisBidRepository) PlaceBid(campaign models.Campaign) (models.Campaign, bool) {

	return models.Campaign{}, true
}

func (r *RedisBidRepository) GetCampaigns() []models.Campaign {
	return nil
}
