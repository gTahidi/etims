package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"etims-test/models"
)

// APIError represents an error returned by the VSCU API
type APIError struct {
	code    string
	message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %s (code: %s)", e.message, e.code)
}

func (e *APIError) Code() string {
	return e.code
}

// VSCUClient handles all API interactions
type VSCUClient struct {
	BaseURL  string
	Tin      string
	BhfId    string
	CmcKey   string
	DvcSrlNo string
	Logger   *logrus.Logger
	Client   *http.Client
}

// NewVSCUClient creates a new VSCU API client
func NewVSCUClient(baseURL, tin, bhfId, cmcKey string, logger *logrus.Logger) *VSCUClient {
	return &VSCUClient{
		BaseURL:  baseURL,
		Tin:      tin,
		BhfId:    bhfId,
		CmcKey:   cmcKey,
		DvcSrlNo: "7ba05e23-850a-44dd-b09a-2eac8405e592",
		Logger:   logger,
		Client:   &http.Client{Timeout: 30 * time.Second},
	}
}

// SendRequest sends a request to the VSCU API
func (c *VSCUClient) SendRequest(method, endpoint string, requestBody interface{}, headers map[string]string) (*models.APIResponse, error) {
	// Validate request if it implements Validator interface
	if validator, ok := requestBody.(models.Validator); ok {
		if err := validator.Validate(); err != nil {
			return nil, fmt.Errorf("request validation failed: %w", err)
		}
	}

	// Convert request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create request
	url := c.BaseURL + endpoint
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("CMC-KEY", c.CmcKey)

	// Set additional headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Log request details
	c.Logger.WithFields(logrus.Fields{
		"method":       method,
		"url":          url,
		"headers":      req.Header,
		"request_body": string(jsonBody),
	}).Info("Sending request")

	// Send request
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log response
	c.Logger.WithFields(logrus.Fields{
		"status_code":   resp.StatusCode,
		"headers":       resp.Header,
		"response_body": string(body),
	}).Info("Received response")

	// Parse response
	var apiResponse models.APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	// Check for API-level errors
	if apiResponse.ResultCd != "000" && apiResponse.ResultCd != "902" && apiResponse.ResultCd != "001" {
		return &apiResponse, &APIError{
			code:    apiResponse.ResultCd,
			message: apiResponse.ResultMsg,
		}
	}

	return &apiResponse, nil
}

// Code List Endpoint
func (c *VSCUClient) GetCodeList(lastReqDt string) (*models.APIResponse, error) {
	req := models.CodeRequest{
		BaseRequest: models.BaseRequest{
			Tin:   c.Tin,
			BhfId: c.BhfId,
		},
		LastReqDt: lastReqDt,
	}
	return c.SendRequest("POST", "/code/selectCodes", req, nil)
}

// Item Classification List Endpoint
func (c *VSCUClient) GetItemClassList(lastReqDt string) (*models.APIResponse, error) {
	req := models.ItemClassRequest{
		BaseRequest: models.BaseRequest{
			Tin:   c.Tin,
			BhfId: c.BhfId,
		},
		LastReqDt: lastReqDt,
	}
	return c.SendRequest("POST", "/itemClass/selectItemsClass", req, nil)
}

// Customer Information Endpoint
func (c *VSCUClient) GetCustomerInfo(custmTin string) (*models.APIResponse, error) {
	req := models.CustomerRequest{
		BaseRequest: models.BaseRequest{
			Tin:   c.Tin,
			BhfId: c.BhfId,
		},
		CustmTin: custmTin,
	}
	return c.SendRequest("POST", "/customer/selectCustomerInfo", req, nil)
}

// Branch Customer Management Endpoint
func (c *VSCUClient) SaveBranchCustomer(customer models.BranchCustomerRequest) (*models.APIResponse, error) {
	customer.BaseRequest = models.BaseRequest{
		Tin:   c.Tin,
		BhfId: c.BhfId,
	}
	return c.SendRequest("POST", "/branches/saveBrancheCustomers", customer, nil)
}

// Item Management Endpoints
func (c *VSCUClient) SaveItem(item models.ItemRequest) (*models.APIResponse, error) {
	item.BaseRequest = models.BaseRequest{
		Tin:   c.Tin,
		BhfId: c.BhfId,
	}
	return c.SendRequest("POST", "/items/saveItems", item, nil)
}

func (c *VSCUClient) GetItems(lastReqDt string) (*models.APIResponse, error) {
	req := models.CodeRequest{
		BaseRequest: models.BaseRequest{
			Tin:   c.Tin,
			BhfId: c.BhfId,
		},
		LastReqDt: lastReqDt,
	}
	return c.SendRequest("POST", "/items/selectItems", req, nil)
}

// Sales Management Endpoints
func (c *VSCUClient) SaveSales(sales models.SalesRequest) (*models.APIResponse, error) {
	sales.BaseRequest = models.BaseRequest{
		Tin:   c.Tin,
		BhfId: c.BhfId,
	}
	return c.SendRequest("POST", "/trnsSales/saveSales", sales, nil)
}

func (c *VSCUClient) GetSales(lastReqDt string) (*models.APIResponse, error) {
	req := models.CodeRequest{
		BaseRequest: models.BaseRequest{
			Tin:   c.Tin,
			BhfId: c.BhfId,
		},
		LastReqDt: lastReqDt,
	}
	return c.SendRequest("POST", "/trnsSales/selectSales", req, nil)
}

// Stock Management Endpoints
func (c *VSCUClient) SaveStock(stock models.StockRequest) (*models.APIResponse, error) {
	// Only set Tin and BhfId if they are empty
	if stock.Tin == "" {
		stock.Tin = c.Tin
	}
	if stock.BhfId == "" {
		stock.BhfId = c.BhfId
	}
	return c.SendRequest("POST", "/stock/saveStockItems", stock, nil)
}

func (c *VSCUClient) GetStockMovements(lastReqDt string) (*models.APIResponse, error) {
	req := models.CodeRequest{
		BaseRequest: models.BaseRequest{
			Tin:   c.Tin,
			BhfId: c.BhfId,
		},
		LastReqDt: lastReqDt,
	}
	return c.SendRequest("POST", "/stock/selectStockMovement", req, nil)
}

// Device Initialization Endpoint
func (c *VSCUClient) InitializeDevice() (*models.APIResponse, error) {
	return c.SendRequest("POST", "/initializer/selectInitInfo", models.InitRequest{
		Tin:      c.Tin,
		BhfId:    c.BhfId,
		DvcSrlNo: c.DvcSrlNo,
	}, nil)
}

// Get Branch List Endpoint
func (c *VSCUClient) GetBranchList(lastReqDt string) (*models.APIResponse, error) {
	return c.SendRequest("POST", "/branches/selectBranches", models.BranchListRequest{
		Tin:       c.Tin,
		BhfId:     c.BhfId,
		LastReqDt: lastReqDt,
	}, nil)
}

// Get Notice List Endpoint
func (c *VSCUClient) GetNoticeList(lastReqDt string) (*models.APIResponse, error) {
	return c.SendRequest("POST", "/notices/selectNotices", models.NoticeRequest{
		Tin:       c.Tin,
		BhfId:     c.BhfId,
		LastReqDt: lastReqDt,
	}, nil)
}

// Save Item Composition Endpoint
func (c *VSCUClient) SaveItemComposition(composition models.ItemCompositionRequest) (*models.APIResponse, error) {
	return c.SendRequest("POST", "/items/saveItemComposition", composition, nil)
}

// Get Imported Items Endpoint
func (c *VSCUClient) GetImportedItems(lastReqDt string) (*models.APIResponse, error) {
	return c.SendRequest("POST", "/imports/selectImportItems", models.ImportedItemRequest{
		Tin:       c.Tin,
		BhfId:     c.BhfId,
		LastReqDt: lastReqDt,
	}, nil)
}

// Save Imported Items Endpoint
func (c *VSCUClient) SaveImportedItems(items models.ImportedItemSaveRequest) (*models.APIResponse, error) {
	return c.SendRequest("POST", "/imports/saveImportItems", items, nil)
}

// Save Purchase Transaction Confirmation
func (c *VSCUClient) SavePurchaseConfirmation(confirmation models.PurchaseConfirmationRequest) (*models.APIResponse, error) {
	return c.SendRequest("POST", "/trnsPurchase/savePurchaseConfirmation", confirmation, nil)
}

// Helper function to get current timestamp in required format
func (c *VSCUClient) GetCurrentTimestamp() string {
	return time.Now().Format("20060102150405")
}

// Helper function to get current date in required format
func (c *VSCUClient) GetCurrentDate() string {
	return time.Now().Format("20060102")
}
