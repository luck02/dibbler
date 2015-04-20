# dibbler

Dibbler is intended as a proof of concept.

The goal is to accept OpenRTB requests and determine 
* if a campaign should target the RequestToBid.
* This involves the campaign being configured to look for certain attributes on the RequestToBid (Placement, Ad Size, etc)
* The bid placed should always be the most profitable bid.
* caveat: If the most profitable bid is locked by other bid processes the bid will be retried a few times and then the next most profitable bid will be attempted until all applicable campaigns are tried and retried.  There is a hard upper limit to how long this can take due to the nature of the bidding process.
* Etc, this is a proof of concept and it could get quite a bit more complex:

Architecture:
* Simple golang web server.  I will probably drop in Tiger Tonic later in order to capitalize on someone elses go routine handling.
* Orchestration that launches in a go routine (not yet) and returns a 200 if a bid has been placed or a 204 if no bid is to be placed.
* dibbler.go is the main web server.
* Calls GetCampaigns to get list of applicable campaigns
* Calls PlaceBids to confirm money is available in account and place a bid.
* Will use Redis key value store to keep client credentials and use watch / exec to ensure no concurrent modification of the amount of cash remaining in a given campaign.
