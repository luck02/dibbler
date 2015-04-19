package repo

import "github.com/luck02/dibbler/dibbler"

type BidRepository interface {
	PlaceBid(dibbler.Campaign) (Campaign, bool)
	GetCampaigns() []Campaign
}

type FakeBidRepository struct {
	CampaignCollection []Campaign
	BidAttempts        int
}

func (r FakeBidRepository) PlaceBid(campaign Campaign) (Campaign, bool) {
	return campaign, true
}

func (r FakeBidRepository) GetCampaigns() []Campaign {
	return r.CampaignCollection
}

type RedisBidRepository struct {
}

func (r RedisBidRepository) PlaceBid(campaign Campaign) (Campaign, bool) {

	return Campaign{}, true
}

func (r RedisBidRepository) GetCampaigns() []Campaign {
	return nil
}
