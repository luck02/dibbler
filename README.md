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

ToDo:
* Come up with a better 'rebid' implementation.  There should be a smarter way to determine when and how quickly to move to the next campaign.  Waiting a random amount of time (even a random amount of time) is not a technically acceptable solution generally speaking.
* Campaign.Target is an interface{}, however it shouldn't be.  Interface{} means 'any type' literally... In reality it's a union of the various target types.  However golang does not have a mechanism to reflect this efficiently.  I have a couple of ideas in mind for this (http://www.jerf.org/iri/post/2917). This would entail moving the JSON query logic into the TargetTypes themselves.  Originally I thought this would be a smell, but if it's a tradeoff between logic in models vs. logic in a service file.  If it lets me strongly type my models it's probably worthwhile.  In addition the really important part is that the selection logic lives entirely in one place, not spread apart over multiple implementation concerns ie cohesive.
* Finish the concrete redis datalayer.  I have this more or less working but I'm pretty sure it will be cleaner as a json representation for the base model object + a json representation of the target object and 2 key value pairs for: TargetType and RemainingAmount
