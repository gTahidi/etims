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
	"etims-test/validator"
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
	BaseURL    string
	Tin        string
	BhfId      string
	CmcKey     string
	DvcSrlNo   string
	Logger     *logrus.Logger
	Client     *http.Client
	Validator  *validator.SchemaValidator
}

// NewVSCUClient creates a new VSCU API client
func NewVSCUClient(baseURL, tin, bhfId, cmcKey string, logger *logrus.Logger) (*VSCUClient, error) {
	// Initialize schema validator
	schemaValidator, err := validator.NewSchemaValidator("/home/daniel/repos/etims/etims-test/shemas", logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize schema validator: %w", err)
	}

	return &VSCUClient{
		BaseURL:    baseURL,
		Tin:        tin,
		BhfId:      bhfId,
		CmcKey:     cmcKey,
		DvcSrlNo:   "7ba05e23-850a-44dd-b09a-2eac8405e592",
		Logger:     logger,
		Validator:  schemaValidator,
		Client: &http.Client{
			Timeout: time.Second * 480,
			Transport: &http.Transport{
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   100,
				IdleConnTimeout:       90 * time.Second,
				Proxy:                 http.ProxyFromEnvironment,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
	}, nil
}

// validateRequest validates the request against its schema before sending
func (c *VSCUClient) validateRequest(request interface{}) error {
	if request == nil {
		return nil
	}

	// Get schema path for the request type
	schemaPath := validator.GetSchemaPath(request)
	if schemaPath == "" {
		c.Logger.Warnf("No schema found for request type %T", request)
		return nil
	}

	// Validate request against schema
	if err := c.Validator.ValidateRequest(request, schemaPath); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	return nil
}

// SendRequest sends a request to the VSCU API
func (c *VSCUClient) SendRequest(method, endpoint string, requestBody interface{}, headers map[string]string) (*models.APIResponse, error) {
	// Validate request body against schema
	if err := c.validateRequest(requestBody); err != nil {
		return nil, err
	}

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
	resp, err := c.SendRequest("POST", "/code/selectCodes", models.CodeListRequest{
		Tin:       c.Tin,
		BhfId:     c.BhfId,
		LastReqDt: lastReqDt,
	}, nil)
	if err != nil {
		return nil, err
	}

	// Handle 001 case
	if resp.ResultCd == "001" {
		c.Logger.Info("No code list found")
		return resp, nil
	}

	// Handle other non-success cases
	if resp.ResultCd != "000" {
		return nil, fmt.Errorf("API error: %s (code: %s)", resp.ResultMsg, resp.ResultCd)
	}

	return resp, nil
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

// Purchase Management Endpoints
func (c *VSCUClient) SavePurchase(purchase models.PurchaseRequest) (*models.APIResponse, error) {
	return c.SendRequest("POST", "/trnsPurchase/savePurchases", purchase, nil)
}

func (c *VSCUClient) GetPurchases(lastReqDt string) (*models.APIResponse, error) {
	req := models.CodeListRequest{
		Tin:       c.Tin,
		BhfId:     c.BhfId,
		LastReqDt: lastReqDt,
	}
	return c.SendRequest("POST", "/trnsPurchase/selectTrnsPurchaseSales", req, nil)
}

// Helper function to get current timestamp in required format
func (c *VSCUClient) GetCurrentTimestamp() string {
	return time.Now().Format("20060102150405")
}

// Helper function to get current date in required format
func (c *VSCUClient) GetCurrentDate() string {
	return time.Now().Format("20060102")
}
