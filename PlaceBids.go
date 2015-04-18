package dibbler

// PlaceBids takes a list of sorted campaigns and a bidRepository.
// Returns the campaign that the bid was successfully placed on.
// If return value Campaign is null then an error will Have been encountered
// error if an error was encountered
// Mechanism:
// * try to place bid in order of campaign in list.
// * BidRepository will go to redis in production.
func PlaceBids(sortedCampaigns []Campaign, bidRepository BidRepository) (Campaign, error) {

	return Campaign{}, nil
}
