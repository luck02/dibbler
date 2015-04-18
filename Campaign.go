package dibbler

type Campaign struct {
	Id              int32
	BidCpm          float32
	DailyBudget     float32
	RemainingBudget float32
	Targeting       interface{}
}

// We could theoretically put the json query here as well... But I think we'd be dirtying up
//   The otherwise clear model
type AdTarget struct {
	Height int
	Width  int
}

type PlacementTarget struct {
	AppName string
}

type CountryTarget struct {
	Country string
}

type OSTarget struct {
	OsType string
}

type SortedCampaigns []Campaign

func (s SortedCampaigns) Len() int {
	return len(s)
}
func (s SortedCampaigns) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortedCampaigns) Less(i, j int) bool {
	return s[i].BidCpm < s[j].BidCpm
}
