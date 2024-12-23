package client

import (
	"etims-test/internal/app/models"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     *logrus.Logger
	Tin        string
	BhfId      string
	CmcKey     string
}

func NewClient(baseURL, tin, bhfId, cmcKey string, logger *logrus.Logger) VSCUClient {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		logger: logger,
		Tin:    tin,
		BhfId:  bhfId,
		CmcKey: cmcKey,
	}
}

func (c *Client) InitializeDevice() (*models.InitResponse, error) {
	return nil, nil
}

func (c *Client) GetCodeList(lastReqDt time.Time) (*models.CodeListResponse, error) {
	return nil, nil
}

func (c *Client) GetItemClassList(lastReqDt time.Time) (*models.ItemClassResponse, error) {
	return nil, nil
}

func (c *Client) GetBranchList(lastReqDt time.Time) (*models.BranchListResponse, error) {
	return nil, nil
}

func (c *Client) GetNoticeList(lastReqDt time.Time) (*models.NoticeListResponse, error) {
	return nil, nil
}

func (c *Client) GetImportedItems(lastReqDt time.Time) (*models.ImportedItemsResponse, error) {
	return nil, nil
}

func (c *Client) SaveBranchCustomer(req models.BranchCustomerRequest) (*models.Response, error) {
	return nil, nil
}

func (c *Client) SaveItem(req models.ItemRequest) (*models.Response, error) {
	return nil, nil
}

func (c *Client) SaveStock(req models.StockRequest) (*models.Response, error) {
	return nil, nil
}

func (c *Client) SavePurchase(req models.PurchaseRequest) (*models.Response, error) {
	return nil, nil
}

func (c *Client) CreateStockMovementFromPurchase(data []byte) (*models.StockRequest, error) {
	return nil, nil
}

func (c *Client) SaveSales(req models.SalesRequest) (*models.Response, error) {
	return nil, nil
}

func (c *Client) CreateStockMovementFromSale(data []byte) (*models.StockRequest, error) {
	return nil, nil
}

func (c *Client) sendRequest(method, endpoint string, body io.Reader, headers map[string]string) (*http.Response, error) {
	url := c.baseURL + endpoint
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.CmcKey)

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	c.logger.WithFields(logrus.Fields{
		"method":  method,
		"url":     url,
		"headers": req.Header,
	}).Debug("Sending request")

	return c.httpClient.Do(req)
}

func (c *Client) GetCurrentTimestamp() time.Time {
	return time.Now()
}

func (c *Client) IsHeadOffice() bool {
	return c.BhfId == "00"
}

func extractLastReqDt(lastReqDt time.Time) string {
	return lastReqDt.Format("20060102150405")
}

func (c *Client) createSignature(params ...string) string {
	concatenatedParams := strings.Join(params, "|")
	return concatenatedParams
}
