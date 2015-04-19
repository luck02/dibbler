package dibbler

type BidRepository interface {
	PlaceBid(Campaign) (Campaign, bool)
	GetCampaigns() []Campaign
}

type FakeBidRepository struct {
	CampaignCollection []Campaign
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
