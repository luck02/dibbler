package repo

import "github.com/luck02/dibbler/models"

type FakeBidRepository struct {
	CampaignCollection       []models.Campaign
	BidAttempts              int
	TestingExhaustedCampaign bool
	BidShouldBeSuccesfull    bool
}

func (r *FakeBidRepository) PlaceBid(campaign models.Campaign) (models.Campaign, bool) {
	r.BidAttempts += 1
	if r.TestingExhaustedCampaign {
		campaign.RemainingBudget = 0
	}
	return campaign, r.BidShouldBeSuccesfull
}

func (r *FakeBidRepository) GetCampaigns() []models.Campaign {
	return r.CampaignCollection
}
