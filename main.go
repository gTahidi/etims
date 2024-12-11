package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
	Tin       string  `json:"tin"`
	BhfId     string  `json:"bhfId"`
	ItemCd    string  `json:"itemCd"`
	ItemClsCd string  `json:"itemClsCd"`
	ItemNm    string  `json:"itemNm"`
	ItemStdNm string  `json:"itemStdNm"`
	OrgnNatCd string  `json:"orgnNatCd"`
	PkgUnitCd string  `json:"pkgUnitCd"`
	QtyUnitCd string  `json:"qtyUnitCd"`
	TaxTyCd   string  `json:"taxTyCd"`
	BtchNo    string  `json:"btchNo"`
	DftPrc    float64 `json:"dftPrc"`
	RegrNm    string  `json:"regrNm"`
	RegrId    string  `json:"regrId"`
	ModrNm    string  `json:"modrNm"`
	ModrId    string  `json:"modrId"`
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
	CmpItemCd     string  `json:"cmpItemCd"`
	CmpItemQty    float64 `json:"cmpItemQty"`
	CmpItemUnitCd string  `json:"cmpItemUnitCd"`
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

func main() {
	// First, initialize the device
	initRequest := InitRequest{
		Tin:      tin,
		BhfId:    bhfId,
		DvcSrlNo: "7ba05e23-850a-44dd-b09a-2eac8405e592",
	}

	initRequestBody, err := json.Marshal(initRequest)
	if err != nil {
		log.Fatal(err)
	}

	// Log the request details
	log.Printf("Sending initialization request to: %s", baseURL+"/initializer/selectInitInfo")
	log.Printf("Request Headers:\ntin: %s\nbhfid: %s\ncmckey: %s", tin, bhfId, cmcKey)
	log.Printf("Request Body: %s", string(initRequestBody))

	initResponse, err := sendRequest(baseURL+"/initializer/selectInitInfo", nil, initRequestBody)
	if err != nil {
		log.Printf("Error making request: %v", err)
		log.Fatal(err)
	}
	defer initResponse.Body.Close()

	body, err := ioutil.ReadAll(initResponse.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		log.Fatal(err)
	}
	log.Printf("Init Response Status: %d", initResponse.StatusCode)
	log.Printf("Init Response Headers: %+v", initResponse.Header)
	log.Printf("Init Response Body: %s", string(body))

	// Only proceed with other requests if initialization is successful
	if initResponse.StatusCode != 200 {
		log.Printf("Device initialization failed with status code: %d", initResponse.StatusCode)
		log.Fatal("Device initialization failed")
	}

	log.Println("Device initialization successful, proceeding with data synchronization...")

	// 1. Code Data Sequence
	codeRequest := CodeRequest{
		Tin:       tin,
		BhfId:     bhfId,
		LastReqDt: "20240101000000",
	}

	codeRequestBody, err := json.Marshal(codeRequest)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Fetching code list...")
	codeResponse, err := sendRequest(baseURL+"/code/selectCodes", nil, codeRequestBody)
	if err != nil {
		log.Fatal(err)
	}
	defer codeResponse.Body.Close()

	codeBody, err := ioutil.ReadAll(codeResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Code List Response Status: %d", codeResponse.StatusCode)
	log.Printf("Code List Response Body: %s", string(codeBody))

	// 2. Notice List
	noticeRequest := NoticeRequest{
		Tin:       tin,
		BhfId:     bhfId,
		LastReqDt: "20240101000000",
	}

	noticeRequestBody, err := json.Marshal(noticeRequest)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Fetching notice list...")
	noticeResponse, err := sendRequest(baseURL+"/notices/selectNotices", nil, noticeRequestBody)
	if err != nil {
		log.Fatal(err)
	}
	defer noticeResponse.Body.Close()

	noticeBody, err := ioutil.ReadAll(noticeResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Notice List Response Status: %d", noticeResponse.StatusCode)
	log.Printf("Notice List Response Body: %s", string(noticeBody))

	// 3. Branch List
	branchRequest := BranchRequest{
		Tin:       tin,
		BhfId:     bhfId,
		LastReqDt: "20240101000000",
	}

	branchRequestBody, err := json.Marshal(branchRequest)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Fetching branch list...")
	branchResponse, err := sendRequest(baseURL+"/branches/selectBranches", nil, branchRequestBody)
	if err != nil {
		log.Fatal(err)
	}
	defer branchResponse.Body.Close()

	branchBody, err := ioutil.ReadAll(branchResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Branch List Response Status: %d", branchResponse.StatusCode)
	log.Printf("Branch List Response Body: %s", string(branchBody))

	// 4. Import Items
	importRequest := ImportItemRequest{
		Tin:       tin,
		BhfId:     bhfId,
		LastReqDt: "20240101000000",
	}

	importRequestBody, err := json.Marshal(importRequest)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Fetching import items...")
	importResponse, err := sendRequest(baseURL+"/imports/selectImportItems", nil, importRequestBody)
	if err != nil {
		log.Fatal(err)
	}
	defer importResponse.Body.Close()

	importBody, err := ioutil.ReadAll(importResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Import Items Response Status: %d", importResponse.StatusCode)
	log.Printf("Import Items Response Body: %s", string(importBody))

	// 5. Purchase Transactions
	purchaseRequest := PurchaseRequest{
		Tin:       tin,
		BhfId:     bhfId,
		LastReqDt: "20240101000000",
	}

	purchaseRequestBody, err := json.Marshal(purchaseRequest)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Fetching purchase transactions...")
	purchaseResponse, err := sendRequest(baseURL+"/trnsPurchase/selectTrnsPurchaseSales", nil, purchaseRequestBody)
	if err != nil {
		log.Fatal(err)
	}
	defer purchaseResponse.Body.Close()

	purchaseBody, err := ioutil.ReadAll(purchaseResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Purchase Transactions Response Status: %d", purchaseResponse.StatusCode)
	log.Printf("Purchase Transactions Response Body: %s", string(purchaseBody))

	// 6. Stock Items
	stockRequest := StockRequest{
		Tin:       tin,
		BhfId:     bhfId,
		ItemCd:    "KE1NTXU0000006",
		RsdQty:    10,
		LastReqDt: "20240101000000",
		RegrId:    "Admin",
		RegrNm:    "Admin",
		ModrId:    "Admin",
		ModrNm:    "Admin",
	}

	stockRequestBody, err := json.Marshal(stockRequest)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Fetching stock items...")
	stockResponse, err := sendRequest(baseURL+"/stock/selectStockItems", nil, stockRequestBody)
	if err != nil {
		log.Fatal(err)
	}
	defer stockResponse.Body.Close()

	stockBody, err := ioutil.ReadAll(stockResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Stock Items Response Status: %d", stockResponse.StatusCode)
	log.Printf("Stock Items Response Body: %s", string(stockBody))

	// 8. Item Classification
	itemClassRequest := ItemClassRequest{
		Tin:       tin,
		BhfId:     bhfId,
		LastReqDt: "20240101000000",
	}

	itemClassRequestBody, err := json.Marshal(itemClassRequest)
	if err != nil {
		log.Fatal(err)
	}

	itemClassResponse, err := sendRequest(baseURL+"/itemClass/selectItemsClass", nil, itemClassRequestBody)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Item Classification Response Status: %d", itemClassResponse.StatusCode)
	log.Printf("Item Classification Response Body: %s", itemClassResponse.Body)

	// 9. Branch Insurance
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
		log.Fatal(err)
	}

	branchInsuranceResponse, err := sendRequest(baseURL+"/branches/saveBrancheInsurances", nil, branchInsuranceRequestBody)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Branch Insurance Response Status: %d", branchInsuranceResponse.StatusCode)
	log.Printf("Branch Insurance Response Body: %s", branchInsuranceResponse.Body)

	// 10. Save Item
	itemRequest := ItemRequest{
		Tin:       tin,
		BhfId:     bhfId,
		ItemCd:    "KE1NTXU0000007",
		ItemClsCd: "5022110801",
		ItemNm:    "Test Item",
		ItemStdNm: "Standard Item",
		OrgnNatCd: "KE",
		PkgUnitCd: "EA",
		QtyUnitCd: "EA",
		TaxTyCd:   "B",
		BtchNo:    "BATCH001",
		DftPrc:    1000.00,
		RegrNm:    "Admin",
		RegrId:    "Admin",
		ModrNm:    "Admin",
		ModrId:    "Admin",
	}

	itemRequestBody, err := json.Marshal(itemRequest)
	if err != nil {
		log.Fatal(err)
	}

	itemResponse, err := sendRequest(baseURL+"/items/saveItems", nil, itemRequestBody)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Save Item Response Status: %d", itemResponse.StatusCode)
	log.Printf("Save Item Response Body: %s", itemResponse.Body)

	// 11. Stock Master
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
		log.Fatal(err)
	}

	stockMasterResponse, err := sendRequest(baseURL+"/stockMaster/saveStockMaster", nil, stockMasterRequestBody)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Stock Master Response Status: %d", stockMasterResponse.StatusCode)
	log.Printf("Stock Master Response Body: %s", stockMasterResponse.Body)

	// Test Customer Information
	customerInfoRequest := CustomerInfoRequest{
		Tin:        tin,
		CustmTin:   "A123456789Z",
		CustmBhfId: "00",
		LastReqDt:  "20240101000000",
	}

	customerInfoRequestBody, err := json.Marshal(customerInfoRequest)
	if err != nil {
		log.Fatal(err)
	}

	customerInfoResponse, err := sendRequest(baseURL+"/customers/selectCustomer", nil, customerInfoRequestBody)
	if err != nil {
		log.Fatal(err)
	}
	defer customerInfoResponse.Body.Close()
	log.Printf("Customer Info Response Status: %d", customerInfoResponse.StatusCode)
	log.Printf("Customer Info Response Body: %s", customerInfoResponse.Body)

	// Save Branch Customer
	branchCustomerRequest := BranchCustomerInfoRequest{
		Tin:     tin,
		BhfId:   bhfId,
		CustNo:  "CUST001",
		CustTin: "A123456789B",
		CustNm:  "Test Customer",
		TelNo:   "0712345678",
		Email:   "test@example.com",
		Fax:     "0712345679",
		RegrId:  "Admin",
		RegrNm:  "Admin",
		ModrId:  "Admin",
		ModrNm:  "Admin",
	}

	branchCustomerRequestBody, err := json.Marshal(branchCustomerRequest)
	if err != nil {
		log.Fatal(err)
	}

	branchCustomerResponse, err := sendRequest(baseURL+"/branches/saveBrancheCustomers", nil, branchCustomerRequestBody)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Branch Customer Response Status: %d", branchCustomerResponse.StatusCode)
	log.Printf("Branch Customer Response Body: %s", branchCustomerResponse.Body)

	// Save Branch User Account
	branchUserRequest := BranchUserRequest{
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

	branchUserRequestBody, err := json.Marshal(branchUserRequest)
	if err != nil {
		log.Fatal(err)
	}

	branchUserResponse, err := sendRequest(baseURL+"/branches/saveBrancheUsers", nil, branchUserRequestBody)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Branch User Response Status: %d", branchUserResponse.StatusCode)
	log.Printf("Branch User Response Body: %s", branchUserResponse.Body)

	// Save Item Composition
	itemCompositionRequest := ItemCompositionRequest{
		Tin:           tin,
		BhfId:         bhfId,
		ItemCd:        "KE1NTXU0000007",
		CmpItemCd:     "KE1NTXU0000008",
		CmpItemQty:    1.0,
		CmpItemUnitCd: "EA",
		RegrId:        "Admin",
		RegrNm:        "Admin",
		ModrId:        "Admin",
		ModrNm:        "Admin",
	}

	itemCompositionRequestBody, err := json.Marshal(itemCompositionRequest)
	if err != nil {
		log.Fatal(err)
	}

	itemCompositionResponse, err := sendRequest(baseURL+"/items/saveItemComposition", nil, itemCompositionRequestBody)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Item Composition Response Status: %d", itemCompositionResponse.StatusCode)
	log.Printf("Item Composition Response Body: %s", itemCompositionResponse.Body)

	// Save Stock In/Out
	stockInOutRequest := StockInOutRequest{
		Tin:   tin,
		BhfId: "00",
		StockItems: []StockItem{
			{
				ItemCd:     "ITEM001",
				ItemClsCd:  "01",
				ItemNm:     "Test Item",
				PkgUnitCd:  "EA",
				QtyUnitCd:  "EA",
				TaxTyCd:    "B",
				Bcd:        "123",
				RegBhfId:   "00",
				Pkg:        1,
				Qty:        1,
				DcRt:       0,
				SupplrTin:  "A123456789Z",
				PchsTyCd:   "01",
				OrgnNatCd:  "KE",
				ItemExprDt: "20241231",
				ItemSttsCd: "01",
				RegrId:     "SYSTEM",
				RegrNm:     "SYSTEM",
				ModrId:     "SYSTEM",
				ModrNm:     "SYSTEM",
			},
		},
	}

	stockInOutRequestBody, err := json.Marshal(stockInOutRequest)
	if err != nil {
		log.Fatal(err)
	}

	stockInOutResponse, err := sendRequest(baseURL+"/stock/saveStockItems", nil, stockInOutRequestBody)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Stock In/Out Response Status: %d", stockInOutResponse.StatusCode)
	log.Printf("Stock In/Out Response Body: %s", stockInOutResponse.Body)

	// Sales Transaction
	salesTransactionRequest := SalesTransactionRequest{
		Tin:         tin,
		BhfId:       "00",
		SalesTyCd:   "01",
		CustTin:     "A123456789Z",
		CustNm:      "Test Customer",
		CustBhfId:   "00",
		SalesSttsCd: "01",
		CfmDt:       time.Now().Format("20060102150405"),
		SaleItems: []SaleItem{
			{
				ItemSeq:    1,
				ItemCd:     "ITEM001",
				ItemClsCd:  "01",
				ItemNm:     "Test Item",
				PkgUnitCd:  "EA",
				QtyUnitCd:  "EA",
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
		PmtTyCd:     "01",
		RegrId:      "SYSTEM",
		RegrNm:      "SYSTEM",
		ModrId:      "SYSTEM",
		ModrNm:      "SYSTEM",
	}

	salesTransactionRequestBody, err := json.Marshal(salesTransactionRequest)
	if err != nil {
		log.Fatal(err)
	}

	salesTransactionResponse, err := sendRequest(baseURL+"/trnsSales/saveSales", nil, salesTransactionRequestBody)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Sales Transaction Response Status: %d", salesTransactionResponse.StatusCode)
	log.Printf("Sales Transaction Response Body: %s", salesTransactionResponse.Body)

	// Stock Movement
	stockMovementRequest := StockMovementRequest{
		Tin:       tin,
		BhfId:     bhfId,
		LastReqDt: time.Now().Format("20060102150405"),
	}

	stockMovementRequestBody, err := json.Marshal(stockMovementRequest)
	if err != nil {
		log.Fatal(err)
	}

	stockMovementResponse, err := sendRequest(baseURL+"/stock/selectStockItems", nil, stockMovementRequestBody)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Stock Movement Response Status: %d", stockMovementResponse.StatusCode)
	log.Printf("Stock Movement Response Body: %s", stockMovementResponse.Body)

	log.Println("Data synchronization completed successfully")
}

func sendRequest(url string, headers map[string]string, requestBody []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// Set common headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("tin", tin)
	req.Header.Set("bhfid", bhfId)
	req.Header.Set("cmckey", cmcKey)

	// Set additional headers if provided
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
