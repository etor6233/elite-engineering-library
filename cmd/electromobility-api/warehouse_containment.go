package main
import (
 "elite.local/enterprise/internal/inventorycontrol"
 "elite.local/enterprise/internal/platform/httpapi"
)
func selectedInventoryModule(service *inventorycontrol.Service, bulk *inventorycontrol.BulkService, warehouse *inventorycontrol.WarehouseService, transfer *inventorycontrol.BulkTransferService) httpapi.InventoryControlModule {
 return httpapi.InventoryControlModule{Service:service,Bulk:bulk,Warehouse:warehouse,Transfer:transfer,RequireDurableReservations:true}
}
