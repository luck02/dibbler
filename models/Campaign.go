package models

// Campaign holds informtion about the bidding strategy of the given campaign.
type Campaign struct {
	ID              int32       `json:"id"`
	BidCpm          float32     `json:"bidcpm"`
	DailyBudget     float32     `json:"daily_budget"`
	RemainingBudget float32     `json:"-"`
	Targeting       interface{} `json:"-"`
}

// AdTarget will match bids based on Height / Width of impression
type AdTarget struct {
	Height int
	Width  int
}

// PlacementTarget will match bids based on AppName - Is placement a technical term or could this be better named (ApplicationNameTarget)
type PlacementTarget struct {
	AppName string
}

// CountryTarget will match bids based on Country
type CountryTarget struct {
	Country string
}

// OSTarget will match bids based on operating system
type OSTarget struct {
	OsType string
}

// SortedCampaigns will sort based on BidCPM then based on (TBD)
type SortedCampaigns []Campaign

func (s SortedCampaigns) Len() int {
	return len(s)
}
func (s SortedCampaigns) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortedCampaigns) Less(i, j int) bool {
	return s[i].BidCpm > s[j].BidCpm
}
