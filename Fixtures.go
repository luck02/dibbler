package dibbler

var otbPlacement = `
{
    "app": {
        "name": "Words With Friends 2 iPad"
    },
    "imp": [
        {
            "banner": {
                "h": 250,
                "w": 300
            },
            "id": "1"
        }
    ],
    "device": {
        "os": "iOS",
        "geo": {
            "city": "Irwin",
            "region": "PA",
            "zip": "15642",
            "country": "USA"
        }
    }
}
`
var otbAd = `
{
    "app": {
        "name": "Test App For ad placement"
    },
    "imp": [
        {
            "banner": {
                "h": 728,
                "w": 1024
            },
            "id": "1"
        }
    ],
    "device": {
        "os": "Android",
        "geo": {
            "city": "Irwin",
            "region": "PA",
            "zip": "15642",
            "country": "USA"
        }
    }
}
`

var CampaignTests = []Campaign{
	{Id: 100101, BidCpm: 0.32, DailyBudget: 35.50, RemainingBudget: 35.50, Targeting: PlacementTarget{AppName: "Words With Friends 2 iPad"}},
	{Id: 100102, BidCpm: 0.04, DailyBudget: 5.25, RemainingBudget: 5.25, Targeting: AdTarget{Height: 728, Width: 1024}},
	{Id: 100103, BidCpm: 0.32, DailyBudget: 15.00, RemainingBudget: 15.00, Targeting: CountryTarget{Country: "USA"}},
	{Id: 100104, BidCpm: 0.15, DailyBudget: 22.00, RemainingBudget: 22.00, Targeting: OSTarget{OsType: "Android"}},
	{Id: 100105, BidCpm: 0.02, DailyBudget: 2.25, RemainingBudget: 2.25, Targeting: CountryTarget{Country: "MEX"}},
	{Id: 100106, BidCpm: 0.16, DailyBudget: 125.00, RemainingBudget: 125.00, Targeting: OSTarget{OsType: "iOS"}},
}
