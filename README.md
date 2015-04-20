# dibbler

Dibbler is intended as a proof of concept.

The goal is to accept OpenRTB requests and determine 
* if a campaign should target the RTB.
* This involves the campaign being configured to look for certain attributes on the RequestToBid 
* The bid placed should always be the most profitable bid.
* caveat: If the most profitable bid is locked by other bid processes the bid will be retried a few times and then the next most profitable bid will be attempted until all applicable campaigns are tried.
* Placement: App Name
* Country:
* Etc, this is a proof of concept and it could get quite a bit more complex:

Architecture:
* Simple golang web server.
* Orchestration that launches in a go routine and returns a 200 if a bid has been placed or a 204 if no bid is to be placed.
* dibbler.go is the main web server.
* Calls GetCampaigns to get list of applicable campaigns
* Calls PlaceBids to confirm money is available in account and place a bid.
* Will use Redis key value store to keep client credentials and use watch / exec to ensure no concurrent modification.
