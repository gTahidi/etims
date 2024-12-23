package client

import (
	"etims-test/internal/app/models"
	"time"
)

type VSCUClient interface {
	InitializeDevice() (*models.InitResponse, error)
	GetCodeList(lastReqDt time.Time) (*models.CodeListResponse, error)
	GetItemClassList(lastReqDt time.Time) (*models.ItemClassResponse, error)
	GetBranchList(lastReqDt time.Time) (*models.BranchListResponse, error)
	GetNoticeList(lastReqDt time.Time) (*models.NoticeListResponse, error)
	GetImportedItems(lastReqDt time.Time) (*models.ImportedItemsResponse, error)
	SaveBranchCustomer(req models.BranchCustomerRequest) (*models.Response, error)
	SaveItem(req models.ItemRequest) (*models.Response, error)
	SaveStock(req models.StockRequest) (*models.Response, error)
	SavePurchase(req models.PurchaseRequest) (*models.Response, error)
	CreateStockMovementFromPurchase(data []byte) (*models.StockRequest, error)
	SaveSales(req models.SalesRequest) (*models.Response, error)
	CreateStockMovementFromSale(data []byte) (*models.StockRequest, error)
	GetCurrentTimestamp() time.Time
}
