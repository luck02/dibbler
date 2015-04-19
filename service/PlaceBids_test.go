package service

import (
	"testing"

	"github.com/luck02/dibbler/fixtures"
	"github.com/luck02/dibbler/repo"
)

func TestPlaceBidsTries3TimesOnEachCampaign(t *testing.T) {
	bidRepo := new(repo.FakeBidRepository)
	success, err := PlaceBids(fixtures.CampaignTests[0:1], bidRepo)
	if err == nil {
		t.Error("error should not be empty")
	}
	if success {
		t.Error("Should not be successfull")
	}
	if bidRepo.BidAttempts != 3 {
		t.Errorf("place bids should try 3 times to place a bid when it fails %d", bidRepo.BidAttempts)
	}
}

func TestPlaceBidsTriesOnceIfCampaignIsExhausted(t *testing.T) {
	bidRepo := new(repo.FakeBidRepository)
	bidRepo.TestingExhaustedCampaign = true
	success, err := PlaceBids(fixtures.CampaignTests[0:1], bidRepo)

	if err == nil {
		t.Error("error should not be empty")
	}
	if success {
		t.Error("Should not be successfull")
	}
	if bidRepo.BidAttempts != 1 {
		t.Errorf("place bid should try once when the campaign is exhausted")
	}
}
