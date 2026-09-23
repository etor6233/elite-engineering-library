package main
import "testing"
func TestSelectedInventoryModuleContainsReceiptlessReservations(t *testing.T) {
 if !selectedInventoryModule(nil,nil,nil,nil).RequireDurableReservations {t.Fatal("legacy reservation creation remains enabled in host")}
}
