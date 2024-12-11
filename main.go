package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	baseURL = "http://20.94.40.140:8088"
	tin     = "P052248771P"
	bhfId   = "00"
	cmcKey  = "D3B478EDFBE54536B8DC9DA691A51440E6278C18104D4D6D904F"
)

type CodeRequest struct {
	Tin       string `json:"tin"`
	BhfId     string `json:"bhfId"`
	LastReqDt string `json:"lastReqDt"`
}

type ItemClassRequest struct {
	Tin       string `json:"tin"`
	BhfId     string `json:"bhfId"`
	LastReqDt string `json:"lastReqDt"`
}

type CustomerRequest struct {
	Tin      string `json:"tin"`
	BhfId    string `json:"bhfId"`
	CustmTin string `json:"custmTin"`
}

type BranchCustomerRequest struct {
	Tin     string `json:"tin"`
	BhfId   string `json:"bhfId"`
	CustNo  string `json:"custNo"`
	CustTin string `json:"custTin"`
	CustNm  string `json:"custNm"`
	UseYn   string `json:"useYn"`
	RegrNm  string `json:"regrNm"`
	RegrId  string `json:"regrId"`
	ModrNm  string `json:"modrNm"`
	ModrId  string `json:"modrId"`
}

type BranchInsuranceRequest struct {
	Tin     string  `json:"tin"`
	BhfId   string  `json:"bhfId"`
	IsrccCd string  `json:"isrccCd"`
	IsrccNm string  `json:"isrccNm"`
	IsrcRt  float64 `json:"isrcRt"`
	UseYn   string  `json:"useYn"`
	RegrNm  string  `json:"regrNm"`
	RegrId  string  `json:"regrId"`
	ModrNm  string  `json:"modrNm"`
	ModrId  string  `json:"modrId"`
}

type ItemRequest struct {
	Tin        string  `json:"tin"`
	BhfId      string  `json:"bhfId"`
	ItemCd     string  `json:"itemCd"`
	ItemClsCd  string  `json:"itemClsCd"`
	ItemTyCd   string  `json:"itemTyCd"`
	ItemNm     string  `json:"itemNm"`
	ItemStdNm  string  `json:"itemStdNm"`
	OrgnNatCd  string  `json:"orgnNatCd"`
	PkgUnitCd  string  `json:"pkgUnitCd"`
	QtyUnitCd  string  `json:"qtyUnitCd"`
	TaxTyCd    string  `json:"taxTyCd"`
	BtchNo     string  `json:"btchNo"`
	DftPrc     float64 `json:"dftPrc"`
	IsrcAplcbYn string  `json:"isrcAplcbYn"`
	UseYn      string  `json:"useYn"`
	RegrNm     string  `json:"regrNm"`
	RegrId     string  `json:"regrId"`
	ModrNm     string  `json:"modrNm"`
	ModrId     string  `json:"modrId"`
}

type StockMasterRequest struct {
	Tin    string `json:"tin"`
	BhfId  string `json:"bhfId"`
	ItemCd string `json:"itemCd"`
	RsdQty int    `json:"rsdQty"`
	RegrId string `json:"regrId"`
	RegrNm string `json:"regrNm"`
	ModrId string `json:"modrId"`
	ModrNm string `json:"modrNm"`
}

type StockRequest struct {
	Tin       string `json:"tin"`
	BhfId     string `json:"bhfId"`
	ItemCd    string `json:"itemCd"`
	RsdQty    int    `json:"rsdQty"`
	LastReqDt string `json:"lastReqDt"`
	RegrId    string `json:"regrId"`
	RegrNm    string `json:"regrNm"`
	ModrId    string `json:"modrId"`
	ModrNm    string `json:"modrNm"`
}

type InitRequest struct {
	Tin      string `json:"tin"`
	BhfId    string `json:"bhfId"`
	DvcSrlNo string `json:"dvcSrlNo"`
}

type NoticeRequest struct {
	Tin       string `json:"tin"`
	BhfId     string `json:"bhfId"`
	LastReqDt string `json:"lastReqDt"`
}

type BranchRequest struct {
	Tin       string `json:"tin"`
	BhfId     string `json:"bhfId"`
	LastReqDt string `json:"lastReqDt"`
}

type ImportItemRequest struct {
	Tin       string `json:"tin"`
	BhfId     string `json:"bhfId"`
	LastReqDt string `json:"lastReqDt"`
}

type PurchaseRequest struct {
	Tin       string `json:"tin"`
	BhfId     string `json:"bhfId"`
	LastReqDt string `json:"lastReqDt"`
}

type SalesRequest struct {
	Tin         string      `json:"tin"`
	BhfId       string      `json:"bhfId"`
	InvcNo      string      `json:"invcNo"`
	SalesTyCd   string      `json:"salesTyCd"`
	RcptTyCd    string      `json:"rcptTyCd"`
	PmtTyCd     string      `json:"pmtTyCd"`
	SalesSttsCd string      `json:"salesSttsCd"`
	CfmDt       string      `json:"cfmDt"`
	SalesDt     string      `json:"salesDt"`
	TotItemCnt  int         `json:"totItemCnt"`
	TaxblAmtA   float64     `json:"taxblAmtA"`
	TaxblAmtB   float64     `json:"taxblAmtB"`
	TaxRtA      float64     `json:"taxRtA"`
	TaxRtB      float64     `json:"taxRtB"`
	TaxAmtA     float64     `json:"taxAmtA"`
	TaxAmtB     float64     `json:"taxAmtB"`
	TotTaxblAmt float64     `json:"totTaxblAmt"`
	TotTaxAmt   float64     `json:"totTaxAmt"`
	TotAmt      float64     `json:"totAmt"`
	RegrId      string      `json:"regrId"`
	RegrNm      string      `json:"regrNm"`
	ModrId      string      `json:"modrId"`
	ModrNm      string      `json:"modrNm"`
	Receipt     Receipt     `json:"receipt"`
	ItemList    []SalesItem `json:"itemList"`
}

type Receipt struct {
	CustTin      string `json:"custTin"`
	CustMblNo    string `json:"custMblNo"`
	RptNo        int    `json:"rptNo"`
	TrdeNm       string `json:"trdeNm"`
	Adrs         string `json:"adrs"`
	TopMsg       string `json:"topMsg"`
	BtmMsg       string `json:"btmMsg"`
	PrchrAcptcYn string `json:"prchrAcptcYn"`
}

type SalesItem struct {
	ItemSeq   int     `json:"itemSeq"`
	ItemCd    string  `json:"itemCd"`
	ItemClsCd string  `json:"itemClsCd"`
	ItemNm    string  `json:"itemNm"`
	PkgUnitCd string  `json:"pkgUnitCd"`
	Pkg       int     `json:"pkg"`
	QtyUnitCd string  `json:"qtyUnitCd"`
	Qty       int     `json:"qty"`
	Prc       float64 `json:"prc"`
	SplyAmt   float64 `json:"splyAmt"`
	DcRt      float64 `json:"dcRt"`
	DcAmt     float64 `json:"dcAmt"`
	TaxTyCd   string  `json:"taxTyCd"`
	TaxblAmt  float64 `json:"taxblAmt"`
	TaxAmt    float64 `json:"taxAmt"`
	TotAmt    float64 `json:"totAmt"`
}

type CustomerInfoRequest struct {
	Tin        string `json:"tin"`
	CustmTin   string `json:"custmTin"`
	CustmBhfId string `json:"custmBhfId"`
	LastReqDt  string `json:"lastReqDt"`
}

type BranchCustomerInfoRequest struct {
	Tin     string `json:"tin"`
	BhfId   string `json:"bhfId"`
	CustNo  string `json:"custNo"`
	CustTin string `json:"custTin"`
	CustNm  string `json:"custNm"`
	TelNo   string `json:"telNo"`
	Email   string `json:"email"`
	Fax     string `json:"fax"`
	RegrId  string `json:"regrId"`
	RegrNm  string `json:"regrNm"`
	ModrId  string `json:"modrId"`
	ModrNm  string `json:"modrNm"`
}

type SalesTransactionRequest struct {
	Tin         string     `json:"tin"`
	BhfId       string     `json:"bhfId"`
	SalesTyCd   string     `json:"salesTyCd"`
	RcptTyCd    string     `json:"rcptTyCd"`
	CustTin     string     `json:"custTin"`
	CustNm      string     `json:"custNm"`
	CustBhfId   string     `json:"custBhfId"`
	SalesSttsCd string     `json:"salesSttsCd"`
	CfmDt       string     `json:"cfmDt"`
	SaleItems   []SaleItem `json:"saleItems"`
	TotItemCnt  int        `json:"totItemCnt"`
	TaxblAmtA   float64    `json:"taxblAmtA"`
	TaxblAmtB   float64    `json:"taxblAmtB"`
	TaxblAmtC   float64    `json:"taxblAmtC"`
	TaxblAmtD   float64    `json:"taxblAmtD"`
	TaxblAmtE   float64    `json:"taxblAmtE"`
	TaxRtA      float64    `json:"taxRtA"`
	TaxRtB      float64    `json:"taxRtB"`
	TaxRtC      float64    `json:"taxRtC"`
	TaxRtD      float64    `json:"taxRtD"`
	TaxRtE      float64    `json:"taxRtE"`
	TaxAmtA     float64    `json:"taxAmtA"`
	TaxAmtB     float64    `json:"taxAmtB"`
	TaxAmtC     float64    `json:"taxAmtC"`
	TaxAmtD     float64    `json:"taxAmtD"`
	TaxAmtE     float64    `json:"taxAmtE"`
	TotTaxblAmt float64    `json:"totTaxblAmt"`
	TotTaxAmt   float64    `json:"totTaxAmt"`
	TotAmt      float64    `json:"totAmt"`
	PmtTyCd     string     `json:"pmtTyCd"`
	RegrId      string     `json:"regrId"`
	RegrNm      string     `json:"regrNm"`
	ModrId      string     `json:"modrId"`
	ModrNm      string     `json:"modrNm"`
}

type StockMovementRequest struct {
	Tin       string `json:"tin"`
	BhfId     string `json:"bhfId"`
	LastReqDt string `json:"lastReqDt"`
}

type BranchUserRequest struct {
	Tin    string `json:"tin"`
	BhfId  string `json:"bhfId"`
	UserId string `json:"userId"`
	UserNm string `json:"userNm"`
	Pwd    string `json:"pwd"`
	RoleCd string `json:"roleCd"`
	UseYn  string `json:"useYn"`
	RegrId string `json:"regrId"`
	RegrNm string `json:"regrNm"`
	ModrId string `json:"modrId"`
	ModrNm string `json:"modrNm"`
}

type ItemCompositionRequest struct {
	Tin           string  `json:"tin"`
	BhfId         string  `json:"bhfId"`
	ItemCd        string  `json:"itemCd"`
	CpstItemCd    string  `json:"cpstItemCd"`
	CpstQty       float64 `json:"cpstQty"`
	CpstUnitCd    string  `json:"cpstUnitCd"`
	RegrId        string  `json:"regrId"`
	RegrNm        string  `json:"regrNm"`
	ModrId        string  `json:"modrId"`
	ModrNm        string  `json:"modrNm"`
}

type StockInOutRequest struct {
	Tin        string      `json:"tin"`
	BhfId      string      `json:"bhfId"`
	StockItems []StockItem `json:"stockItems"`
	LastReqDt  string      `json:"lastReqDt"`
	RegrId     string      `json:"regrId"`
	RegrNm     string      `json:"regrNm"`
	ModrId     string      `json:"modrId"`
	ModrNm     string      `json:"modrNm"`
}

type StockItem struct {
	ItemCd     string  `json:"itemCd"`
	ItemClsCd  string  `json:"itemClsCd"`
	ItemNm     string  `json:"itemNm"`
	PkgUnitCd  string  `json:"pkgUnitCd"`
	QtyUnitCd  string  `json:"qtyUnitCd"`
	TaxTyCd    string  `json:"taxTyCd"`
	Bcd        string  `json:"bcd"`
	RegBhfId   string  `json:"regBhfId"`
	Pkg        int     `json:"pkg"`
	Qty        int     `json:"qty"`
	DcRt       float64 `json:"dcRt"`
	SupplrTin  string  `json:"supplrTin"`
	PchsTyCd   string  `json:"pchsTyCd"`
	OrgnNatCd  string  `json:"orgnNatCd"`
	ItemExprDt string  `json:"itemExprDt"`
	ItemSttsCd string  `json:"itemSttsCd"`
	RegrId     string  `json:"regrId"`
	RegrNm     string  `json:"regrNm"`
	ModrId     string  `json:"modrId"`
	ModrNm     string  `json:"modrNm"`
}

type SaleItem struct {
	ItemSeq    int     `json:"itemSeq"`
	ItemCd     string  `json:"itemCd"`
	ItemClsCd  string  `json:"itemClsCd"`
	ItemNm     string  `json:"itemNm"`
	PkgUnitCd  string  `json:"pkgUnitCd"`
	QtyUnitCd  string  `json:"qtyUnitCd"`
	Pkg        int     `json:"pkg"`
	Qty        int     `json:"qty"`
	PrcAmt     float64 `json:"prcAmt"`
	DcRt       float64 `json:"dcRt"`
	DcAmt      float64 `json:"dcAmt"`
	TaxTyCd    string  `json:"taxTyCd"`
	TaxAmt     float64 `json:"taxAmt"`
	TotAmt     float64 `json:"totAmt"`
	ItemExprDt string  `json:"itemExprDt"`
}

type GetItemRequest struct {
	Tin       string `json:"tin"`
	BhfId     string `json:"bhfId"`
	LastReqDt string `json:"lastReqDt"`
}

func main() {
	// Configure logrus
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
		DisableColors:    false,
		DisableTimestamp: false,
	})
	log.SetLevel(logrus.InfoLevel)

	// Add file and line number to log output
	log.SetReportCaller(true)

	startTime := time.Now()
	log.WithFields(logrus.Fields{
		"tin":    tin,
		"bhfId":  bhfId,
		"cmcKey": cmcKey,
	}).Info("Starting data synchronization process")

	// Create a session ID for this run
	sessionID := fmt.Sprintf("session_%d", startTime.UnixNano())
	logger := log.WithField("session_id", sessionID)

	// Track the number of successful and failed requests
	stats := struct {
		successful int
		failed     int
	}{0, 0}

	// Helper function to handle common request/response pattern
	makeRequest := func(endpoint string, requestBody []byte, description string) error {
		requestLog := logger.WithFields(logrus.Fields{
			"endpoint":   endpoint,
			"request_id": fmt.Sprintf("%s_%d", description, time.Now().UnixNano()),
		})

		requestLog.WithField("request_body", string(requestBody)).Info(fmt.Sprintf("Sending %s request", description))

		response, err := sendRequest(baseURL+endpoint, nil, requestBody)
		if err != nil {
			stats.failed++
			requestLog.WithError(err).Error(fmt.Sprintf("Failed to send %s request", description))
			return fmt.Errorf("%s request failed: %w", description, err)
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			stats.failed++
			requestLog.WithError(err).Error(fmt.Sprintf("Failed to read %s response body", description))
			return fmt.Errorf("failed to read %s response body: %w", description, err)
		}

		requestLog.WithFields(logrus.Fields{
			"status_code": response.StatusCode,
			"headers":     response.Header,
			"body":        string(body),
		}).Info(fmt.Sprintf("Received %s response", description))

		if response.StatusCode != http.StatusOK {
			stats.failed++
			return fmt.Errorf("%s request failed with status code: %d", description, response.StatusCode)
		}

		stats.successful++
		return nil
	}

	// First, initialize the device
	initRequest := InitRequest{
		Tin:      tin,
		BhfId:    bhfId,
		DvcSrlNo: "7ba05e23-850a-44dd-b09a-2eac8405e592",
	}

	initRequestBody, err := json.Marshal(initRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal init request")
	}

	logger.WithFields(logrus.Fields{
		"tin":         tin,
		"bhfId":       bhfId,
		"dvcSrlNo":    initRequest.DvcSrlNo,
		"requestBody": string(initRequestBody),
	}).Info("Sending initialization request")

	// Send initialization request
	response, err := sendRequest(baseURL+"/initializer/selectInitInfo", map[string]string{
		"Content-Type": "application/json; charset=utf-8",
		"Accept":       "application/json",
		"CMC-KEY":     cmcKey,
		"User-Agent":  "etims-client/1.0",
	}, initRequestBody)
	if err != nil {
		logger.WithError(err).Fatal("Failed to send initialization request")
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.WithError(err).Fatal("Failed to read initialization response body")
	}

	// Log response details
	logger.WithFields(logrus.Fields{
		"statusCode":    response.StatusCode,
		"headers":       response.Header,
		"responseBody": string(body),
		"contentType":  response.Header.Get("Content-Type"),
	}).Info("Received initialization response")

	if response.StatusCode != http.StatusOK {
		logger.WithFields(logrus.Fields{
			"statusCode": response.StatusCode,
			"body":      string(body),
		}).Fatal("Initialization request failed")
	}

	logger.Info("Device initialization successful, proceeding with data synchronization...")

	// 1. Code Data Sequence
	codeRequest := CodeRequest{
		Tin:       tin,
		BhfId:     bhfId,
		LastReqDt: time.Now().Format("20060102150405"), // Current timestamp
	}

	codeRequestBody, err := json.Marshal(codeRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal code request")
	}

	if err := makeRequest("/code/selectCodes", codeRequestBody, "code list"); err != nil {
		logger.WithError(err).Fatal("Failed to fetch code list")
	}

	// 2. Notice List
	noticeRequest := NoticeRequest{
		Tin:       tin,
		BhfId:     bhfId,
		LastReqDt: time.Now().Format("20060102150405"), // Current timestamp
	}

	noticeRequestBody, err := json.Marshal(noticeRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal notice request")
	}

	if err := makeRequest("/notices/selectNotices", noticeRequestBody, "notice list"); err != nil {
		logger.WithError(err).Fatal("Failed to fetch notice list")
	}

	// 3. Branch List
	branchRequest := BranchRequest{
		Tin:       tin,
		BhfId:     bhfId,
		LastReqDt: time.Now().Format("20060102150405"), // Current timestamp
	}

	branchRequestBody, err := json.Marshal(branchRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal branch request")
	}

	if err := makeRequest("/branches/selectBranches", branchRequestBody, "branch list"); err != nil {
		logger.WithError(err).Fatal("Failed to fetch branch list")
	}

	// 4. Import Items
	importRequest := ImportItemRequest{
		Tin:       tin,
		BhfId:     bhfId,
		LastReqDt: time.Now().Format("20060102150405"), // Current timestamp
	}

	importRequestBody, err := json.Marshal(importRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal import request")
	}

	if err := makeRequest("/imports/selectImportItems", importRequestBody, "import items"); err != nil {
		logger.WithError(err).Fatal("Failed to fetch import items")
	}

	// 5. Purchase Transactions
	purchaseRequest := PurchaseRequest{
		Tin:       tin,
		BhfId:     bhfId,
		LastReqDt: time.Now().Format("20060102150405"), // Current timestamp
	}

	purchaseRequestBody, err := json.Marshal(purchaseRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal purchase request")
	}

	if err := makeRequest("/trnsPurchase/selectTrnsPurchaseSales", purchaseRequestBody, "purchase transactions"); err != nil {
		logger.WithError(err).Fatal("Failed to fetch purchase transactions")
	}

	// 6. Stock Items
	stockRequest := StockRequest{
		Tin:       tin,
		BhfId:     bhfId,
		ItemCd:    "KE1NTXU0000006",
		RsdQty:    10,
		LastReqDt: time.Now().Format("20060102150405"), // Current timestamp
		RegrId:    "Admin",
		RegrNm:    "Admin",
		ModrId:    "Admin",
		ModrNm:    "Admin",
	}

	stockRequestBody, err := json.Marshal(stockRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal stock request")
	}

	if err := makeRequest("/stock/selectStockItems", stockRequestBody, "stock items"); err != nil {
		logger.WithError(err).Fatal("Failed to fetch stock items")
	}

	// 7. Item Classification List
	itemClassRequest := ItemClassRequest{
		Tin:       tin,
		BhfId:     bhfId,
		LastReqDt: time.Now().Format("20060102150405"), // Current timestamp
	}

	itemClassRequestBody, err := json.Marshal(itemClassRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal item classification request")
	}

	if err := makeRequest("/itemClass/selectItemsClass", itemClassRequestBody, "item classification list"); err != nil {
		logger.WithError(err).Fatal("Failed to fetch item classification list")
	}

	// 8. Customer List (PIN List)
	customerRequest := CustomerRequest{
		Tin:      tin,
		BhfId:    bhfId,
		CustmTin: "A123456789Z", // Must be between 9-15 characters
	}

	customerRequestBody, err := json.Marshal(customerRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal customer request")
	}

	if err := makeRequest("/customers/selectCustomer", customerRequestBody, "customer list"); err != nil {
		logger.WithError(err).Fatal("Failed to fetch customer list")
	}

	// 9. Send Branch Customer Information
	branchCustomerRequest := BranchCustomerRequest{
		Tin:     tin,
		BhfId:   bhfId,
		CustNo:  "CUST001",
		CustTin: "A123456789Z",
		CustNm:  "Test Customer",
		UseYn:   "Y",
		RegrNm:  "Admin",
		RegrId:  "Admin",
		ModrNm:  "Admin",
		ModrId:  "Admin",
	}

	branchCustomerRequestBody, err := json.Marshal(branchCustomerRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal branch customer request")
	}

	if err := makeRequest("/branches/saveBrancheCustomers", branchCustomerRequestBody, "branch customer information"); err != nil {
		logger.WithError(err).Fatal("Failed to send branch customer information")
	}

	// 10. Send Branch User Account
	branchUserRequest := BranchUserRequest{
		Tin:    tin,
		BhfId:  bhfId,
		UserId: "user001",
		UserNm: "Test User",
		Pwd:    "password123",
		RoleCd: "ADMIN",
		UseYn:  "Y",
		RegrId: "Admin",
		RegrNm: "Admin",
		ModrId: "Admin",
		ModrNm: "Admin",
	}

	branchUserRequestBody, err := json.Marshal(branchUserRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal branch user request")
	}

	if err := makeRequest("/branches/saveBrancheUsers", branchUserRequestBody, "branch user account"); err != nil {
		logger.WithError(err).Fatal("Failed to send branch user account")
	}

	// 11. Item Classification
	itemClassRequest = ItemClassRequest{
		Tin:       tin,
		BhfId:     bhfId,
		LastReqDt: time.Now().Format("20060102150405"), // Current timestamp
	}

	itemClassRequestBody, err = json.Marshal(itemClassRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal item classification request")
	}

	if err := makeRequest("/itemClass/selectItemsClass", itemClassRequestBody, "item classification"); err != nil {
		logger.WithError(err).Fatal("Failed to fetch item classification")
	}

	// 12. Branch Insurance
	branchInsuranceRequest := BranchInsuranceRequest{
		Tin:     tin,
		BhfId:   bhfId,
		IsrccCd: "ISRCC01",
		IsrccNm: "Sample Insurance",
		IsrcRt:  16.0,
		UseYn:   "Y",
		RegrNm:  "Admin",
		RegrId:  "Admin",
		ModrNm:  "Admin",
		ModrId:  "Admin",
	}

	branchInsuranceRequestBody, err := json.Marshal(branchInsuranceRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal branch insurance request")
	}

	if err := makeRequest("/branches/saveBrancheInsurances", branchInsuranceRequestBody, "branch insurance"); err != nil {
		logger.WithError(err).Fatal("Failed to send branch insurance")
	}

	// 13. Save Item
	itemRequest := ItemRequest{
		Tin:         tin,
		BhfId:       bhfId,
		ItemCd:      "KE1NTXU0000007",
		ItemClsCd:   "5022110801",
		ItemTyCd:    "1",
		ItemNm:      "Test Item",
		ItemStdNm:   "Standard Item",
		OrgnNatCd:   "KE",
		PkgUnitCd:   "NT",
		QtyUnitCd:   "U",
		TaxTyCd:     "B",
		BtchNo:      "",
		DftPrc:      1000.00,
		IsrcAplcbYn: "N",
		UseYn:       "Y",
		RegrNm:      "Admin",
		RegrId:      "Admin",
		ModrNm:      "Admin",
		ModrId:      "Admin",
	}

	itemRequestBody, err := json.Marshal(itemRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal item request")
	}

	if err := makeRequest("/items/saveItems", itemRequestBody, "item"); err != nil {
		logger.WithError(err).Fatal("Failed to send item")
	}

	// 14. Stock Master
	stockMasterRequest := StockMasterRequest{
		Tin:    tin,
		BhfId:  bhfId,
		ItemCd: "KE1NTXU0000007",
		RsdQty: 100,
		RegrId: "Admin",
		RegrNm: "Admin",
		ModrId: "Admin",
		ModrNm: "Admin",
	}

	stockMasterRequestBody, err := json.Marshal(stockMasterRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal stock master request")
	}

	if err := makeRequest("/stockMaster/saveStockMaster", stockMasterRequestBody, "stock master"); err != nil {
		logger.WithError(err).Fatal("Failed to send stock master")
	}

	// Test Customer Information
	customerInfoRequest := CustomerInfoRequest{
		Tin:        tin,
		CustmTin:   "A123456789Z",
		CustmBhfId: "00",
		LastReqDt:  time.Now().Format("20060102150405"), // Current timestamp
	}

	customerInfoRequestBody, err := json.Marshal(customerInfoRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal customer info request")
	}

	if err := makeRequest("/customers/selectCustomer", customerInfoRequestBody, "customer info"); err != nil {
		logger.WithError(err).Fatal("Failed to fetch customer info")
	}

	// Save Branch Customer
	branchCustomerRequest = BranchCustomerRequest{
		Tin:     tin,
		BhfId:   bhfId,
		CustNo:  "CUST001",
		CustTin: "A123456789B",
		CustNm:  "Test Customer",
		UseYn:   "Y",
		RegrId:  "Admin",
		RegrNm:  "Admin",
		ModrId:  "Admin",
		ModrNm:  "Admin",
	}

	branchCustomerRequestBody, err = json.Marshal(branchCustomerRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal branch customer request")
	}

	if err := makeRequest("/branches/saveBrancheCustomers", branchCustomerRequestBody, "branch customer"); err != nil {
		logger.WithError(err).Fatal("Failed to send branch customer")
	}

	// Save Branch User Account
	branchUserRequest = BranchUserRequest{
		Tin:    tin,
		BhfId:  bhfId,
		UserId: "USER001",
		UserNm: "Test User",
		Pwd:    "password123",
		RoleCd: "ADMIN",
		UseYn:  "Y",
		RegrId: "Admin",
		RegrNm: "Admin",
		ModrId: "Admin",
		ModrNm: "Admin",
	}

	branchUserRequestBody, err = json.Marshal(branchUserRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal branch user request")
	}

	if err := makeRequest("/branches/saveBrancheUsers", branchUserRequestBody, "branch user account"); err != nil {
		logger.WithError(err).Fatal("Failed to send branch user account")
	}

	// Save Item Composition
	itemCompositionRequest := ItemCompositionRequest{
		Tin:        tin,
		BhfId:      bhfId,
		ItemCd:     "KE1NTXU0000007", // Using the item we created earlier
		CpstItemCd: "KE1NTXU0000006", // Using existing item from stock request
		CpstQty:    1.0,
		CpstUnitCd: "U",
		RegrId:     "Admin",
		RegrNm:     "Admin",
		ModrId:     "Admin",
		ModrNm:     "Admin",
	}

	itemCompositionRequestBody, err := json.Marshal(itemCompositionRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal item composition request")
	}

	if err := makeRequest("/items/saveItemComposition", itemCompositionRequestBody, "item composition"); err != nil {
		logger.WithError(err).Fatal("Failed to send item composition")
	}

	// Stock In/Out
	stockInOutRequest := StockInOutRequest{
		Tin:    tin,
		BhfId:  bhfId,
		StockItems: []StockItem{
			{
				ItemCd:     "KE1NTXU0000007", // Using the item we created earlier
				ItemClsCd:  "5022110801",
				ItemNm:     "Test Item",
				PkgUnitCd:  "NT",
				QtyUnitCd:  "U",
				TaxTyCd:    "B",
				Bcd:        "",
				RegBhfId:   bhfId,
				Pkg:        1,
				Qty:        10,
				DcRt:       0,
				SupplrTin:  tin,           // Using our own TIN as supplier
				PchsTyCd:   "NS",          // Changed to NS (Normal Stock)
				OrgnNatCd:  "KE",          // Origin nation code
				ItemExprDt: "20241231",    // Item expiry date
				ItemSttsCd: "01",          // Item status code
				RegrId:     "Admin",
				RegrNm:     "Admin",
				ModrId:     "Admin",
				ModrNm:     "Admin",
			},
		},
		LastReqDt: time.Now().Format("20060102150405"), // Using current timestamp
		RegrId:    "Admin",
		RegrNm:    "Admin",
		ModrId:    "Admin",
		ModrNm:    "Admin",
	}

	stockInOutRequestBody, err := json.Marshal(stockInOutRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal stock in/out request")
	}

	if err := makeRequest("/stock/saveStockItems", stockInOutRequestBody, "stock in/out"); err != nil {
		logger.WithError(err).Fatal("Failed to send stock in/out")
	}

	// Get Item Information
	itemInfoRequest := GetItemRequest{
		Tin:       tin,
		BhfId:     bhfId,
		LastReqDt: time.Now().Format("20060102150405"), // Current timestamp
	}

	itemInfoRequestBody, err := json.Marshal(itemInfoRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal item info request")
	}

	if err := makeRequest("/items/selectItems", itemInfoRequestBody, "item information"); err != nil {
		logger.WithError(err).Fatal("Failed to fetch item information")
	}

	// Send Converted Import Item Information
	importUpdateRequest := struct {
		Tin            string  `json:"tin"`
		BhfId          string  `json:"bhfId"`
		TaskCd         string  `json:"taskCd"`
		DclDe          string  `json:"dclDe"`
		ItemSeq        int     `json:"itemSeq"`
		HsCd           string  `json:"hsCd"`
		ItemClsCd      string  `json:"itemClsCd"`
		ItemCd         string  `json:"itemCd"`
		ImptItemSttsCd string  `json:"imptItemSttsCd"`
		Remark         string  `json:"remark"`
		ModrNm         string  `json:"modrNm"`
		ModrId         string  `json:"modrId"`
	}{
		Tin:            tin,
		BhfId:          bhfId,
		TaskCd:         "2231943",
		DclDe:          "20240101",
		ItemSeq:        1,
		HsCd:           "1231531231",
		ItemClsCd:      "5022110801",
		ItemCd:         "KE1NTXU0000001",
		ImptItemSttsCd: "1",
		Remark:         "Import update",
		ModrNm:         "Admin",
		ModrId:         "Admin",
	}

	importUpdateRequestBody, err := json.Marshal(importUpdateRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal import update request")
	}

	if err := makeRequest("/imports/updateImportItems", importUpdateRequestBody, "import update"); err != nil {
		logger.WithError(err).Fatal("Failed to update import items")
	}

	// Sales Transaction
	salesTransactionRequest := SalesTransactionRequest{
		Tin:         tin,
		BhfId:       bhfId,
		SalesTyCd:   "NS",        // Changed to NS (Normal Sale)
		RcptTyCd:    "NR",        // Changed to NR (Normal Receipt)
		CustTin:     tin,         // Using our own TIN as customer
		CustNm:      "Test Customer",
		CustBhfId:   bhfId,
		SalesSttsCd: "02",        // 02: Completed
		CfmDt:       "20241211085127", // Current timestamp
		SaleItems: []SaleItem{
			{
				ItemSeq:    1,
				ItemCd:     "KE1NTXU0000007", // Using the item we created earlier
				ItemClsCd:  "5022110801",
				ItemNm:     "Test Item",
				PkgUnitCd:  "NT",
				QtyUnitCd:  "U",
				Pkg:        1,
				Qty:        1,
				PrcAmt:     1000.00,
				DcRt:       0,
				DcAmt:      0,
				TaxTyCd:    "B",
				TaxAmt:     160.00,
				TotAmt:     1160.00,
				ItemExprDt: "20241231",
			},
		},
		TotItemCnt:  1,
		TaxblAmtA:   0,
		TaxblAmtB:   1000.00,
		TaxblAmtC:   0,
		TaxblAmtD:   0,
		TaxblAmtE:   0,
		TaxRtA:      0,
		TaxRtB:      16,
		TaxRtC:      0,
		TaxRtD:      0,
		TaxRtE:      0,
		TaxAmtA:     0,
		TaxAmtB:     160.00,
		TaxAmtC:     0,
		TaxAmtD:     0,
		TaxAmtE:     0,
		TotTaxblAmt: 1000.00,
		TotTaxAmt:   160.00,
		TotAmt:      1160.00,
		PmtTyCd:     "01",       // 01: Cash
		RegrId:      "Admin",
		RegrNm:      "Admin",
		ModrId:      "Admin",
		ModrNm:      "Admin",
	}

	salesTransactionRequestBody, err := json.Marshal(salesTransactionRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal sales transaction request")
	}

	if err := makeRequest("/trnsSales/saveSales", salesTransactionRequestBody, "sales transaction"); err != nil {
		logger.WithError(err).Fatal("Failed to send sales transaction")
	}

	// Stock Movement
	stockMovementRequest := StockMovementRequest{
		Tin:       tin,
		BhfId:     bhfId,
		LastReqDt: time.Now().Format("20060102150405"), // Current timestamp
	}

	stockMovementRequestBody, err := json.Marshal(stockMovementRequest)
	if err != nil {
		logger.WithError(err).Fatal("Failed to marshal stock movement request")
	}

	if err := makeRequest("/stock/selectStockItems", stockMovementRequestBody, "stock movement"); err != nil {
		logger.WithError(err).Fatal("Failed to fetch stock movement")
	}

	logger.WithFields(logrus.Fields{
		"successful_requests": stats.successful,
		"failed_requests":     stats.failed,
	}).Info("Data synchronization completed")

	logger.Info("Data synchronization completed successfully")
}

func sendRequest(url string, headers map[string]string, requestBody []byte) (*http.Response, error) {
	log := logrus.WithFields(logrus.Fields{
		"url":     url,
		"headers": headers,
		"method":  "POST",
	})

	// Generate request ID for tracing
	requestID := fmt.Sprintf("%d", time.Now().UnixNano())
	log = log.WithField("request_id", requestID)

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Request-ID", requestID)
	req.Header.Set("CMC-KEY", cmcKey)
	req.Header.Set("User-Agent", "etims-client/1.0")

	// Set custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Log request details
	log.WithFields(logrus.Fields{
		"body_size":     len(requestBody),
		"body":          string(requestBody),
		"request_headers": req.Header,
	}).Info("Sending request")

	// Record start time for request duration
	start := time.Now()

	// Send request with longer timeout
	client := &http.Client{
		Timeout: time.Second * 360,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Error("Request failed")
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Calculate request duration
	duration := time.Since(start)

	// Log response details
	log.WithFields(logrus.Fields{
		"status_code":       resp.StatusCode,
		"duration_ms":       duration.Milliseconds(),
		"content_length":    resp.ContentLength,
		"response_headers":  resp.Header,
		"content_type":      resp.Header.Get("Content-Type"),
	}).Info("Received response")

	return resp, nil
}
