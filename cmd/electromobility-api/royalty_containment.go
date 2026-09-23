package main

import (
 "elite.local/enterprise/internal/platform/httpapi"
 "elite.local/enterprise/internal/royalty"
)

// AUTHORED: the selected host must not expose receipt-less legacy creation.
func selectedRoyaltyModule(service *royalty.Service) httpapi.RoyaltyModule {
 return httpapi.RoyaltyModule{Service: service, RequireDurableCommands: true}
}
