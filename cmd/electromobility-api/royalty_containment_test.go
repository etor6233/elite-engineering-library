package main

import "testing"

func TestSelectedRoyaltyModuleContainsLegacyReceiptLoss(t *testing.T) {
 if !selectedRoyaltyModule(nil).RequireDurableCommands { t.Fatal("selected host exposes receipt-less creation") }
}
