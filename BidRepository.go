package dibbler

type BidRepository interface {
	PlaceBid(Campaign) bool
	GetCampaigns() []Campaign
}

type FakeBidRepository struct {
	CampaignCollection []Campaign
}

func (r FakeBidRepository) PlaceBid(campaign Campaign) bool {
	return true
}

func (r FakeBidRepository) GetCampaigns() []Campaign {
	return r.CampaignCollection
}
