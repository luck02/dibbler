package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/luck02/dibbler/models"
	"github.com/luck02/dibbler/repo"
)

// PlaceBids takes a list of sorted campaigns and a bidRepository.
// Returns the campaign that the bid was successfully placed on.
// If return value Campaign is null then an error will Have been encountered
// error if an error was encountered
// Mechanism:
// * try to place bid in order of campaign in list.
// * BidRepository will go to redis in production.
// * So if I bid on first, fail...  But there's still budget left.  That means the campaign is in contention.
//      I should probably do something like build a function that determines retries based on relative value between cpm's and such
//		however, I think for a proof of concept / practice assignment that's probably more than required.
func PlaceBids(sortedCampaigns []models.Campaign, bidRepository repo.BidRepository) (bool, error) {
	retryCount := 3 //This needs to be configurable and possibly a ratio as mentioned above

	for _, campaign := range sortedCampaigns {

		for i := 0; i < retryCount; i++ {
			if resultCampaign, success, err := bidRepository.PlaceBid(campaign); err != nil {
				return false, err
			} else if success {
				return true, nil
			} else if resultCampaign.RemainingBudget <= 0 {
				fmt.Printf("Campaign %v exhausted", resultCampaign)
				break // campaign exhausted
			} else {
				time.Sleep(time.Duration(rand.Intn(100)) * time.Microsecond) // I think this is a smell
			}
		}
	}
	return false, errors.New("Unable to bid on any eligible campaign")
}
