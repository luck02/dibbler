package dibbler

import (
	"testing"

	fakeDibbler "./" //mock

	"code.google.com/p/gomock/gomock"
)

func testPlaceBidsTries3TimesOnEachCampaign(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeDibbler.MOCK().SetController(ctrl)
	fakeBidRepository := FakeBidRepository{}
	success, error := PlaceBids(CampaignTests, fakeBidRepository)
	bidrepo := fakeDibbler.FakeBidRepository{}
	/*if fakeBidRepository.CountBids != 3 {
		t.Error("place bids should try 3 times to place a bid when it fails")
	}*/
}
