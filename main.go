package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	baseURL = "http://20.94.40.140:8088"
	tin     = "P052248771P"
	bhfId   = "00"
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

type ItemRequest struct {
	Tin         string `json:"tin"`
	BhfId       string `json:"bhfId"`
	ItemCd      string `json:"itemCd"`
	ItemClsCd   string `json:"itemClsCd"`
	ItemTyCd    string `json:"itemTyCd"`
	ItemNm      string `json:"itemNm"`
	ItemStdNm   string `json:"itemStdNm"`
	OrgnNatCd   string `json:"orgnNatCd"`
	PkgUnitCd   string `json:"pkgUnitCd"`
	QtyUnitCd   string `json:"qtyUnitCd"`
	TaxTyCd     string `json:"taxTyCd"`
	BtchNo      string `json:"btchNo"`
	Bcd         string `json:"bcd"`
	DftPrc      int    `json:"dftPrc"`
	GrpPrcL1    int    `json:"grpPrcL1"`
	GrpPrcL2    int    `json:"grpPrcL2"`
	GrpPrcL3    int    `json:"grpPrcL3"`
	GrpPrcL4    int    `json:"grpPrcL4"`
	GrpPrcL5    int    `json:"grpPrcL5"`
	AddInfo     string `json:"addInfo"`
	SftyQty     string `json:"sftyQty"`
	IsrcAplcbYn string `json:"isrcAplcbYn"`
	UseYn       string `json:"useYn"`
	RegrNm      string `json:"regrNm"`
	RegrId      string `json:"regrId"`
	ModrNm      string `json:"modrNm"`
	ModrId      string `json:"modrId"`
}

type StockRequest struct {
	Tin    string `json:"tin"`
	BhfId  string `json:"bhfId"`
	ItemCd string `json:"itemCd"`
	RsdQty int    `json:"rsdQty"`
	RegrId string `json:"regrId"`
	RegrNm string `json:"regrNm"`
	ModrNm string `json:"modrNm"`
	ModrId string `json:"modrId"`
}

func main() {
	// Create request headers
	headers := map[string]string{
		"tin":   tin,
		"bhfId": bhfId,
	}

	// Create request bodies
	codeRequest := CodeRequest{
		Tin:       tin,
		BhfId:     bhfId,
		LastReqDt: "20230328000000",
	}

	itemClassRequest := ItemClassRequest{
		Tin:       tin,
		BhfId:     bhfId,
		LastReqDt: "20180523000000",
	}

	customerRequest := CustomerRequest{
		Tin:      tin,
		BhfId:    bhfId,
		CustmTin: "A987654321Z",
	}

	branchCustomerRequest := BranchCustomerRequest{
		Tin:     tin,
		BhfId:   bhfId,
		CustNo:  "999991113",
		CustTin: "A123456789Z",
		CustNm:  "MTEJA LIMITED",
		UseYn:   "Y",
		RegrNm:  "Admin",
		RegrId:  "Admin",
		ModrNm:  "Admin",
		ModrId:  "Admin",
	}

	itemRequest := ItemRequest{
		Tin:         tin,
		BhfId:       bhfId,
		ItemCd:      "KE1NTXU0000006",
		ItemClsCd:   "5059690800",
		ItemTyCd:    "1",
		ItemNm:      "test materialitem 3",
		ItemStdNm:   "",
		OrgnNatCd:   "KE",
		PkgUnitCd:   "NT",
		QtyUnitCd:   "U",
		TaxTyCd:     "B",
		BtchNo:      "",
		Bcd:         "",
		DftPrc:      3500,
		GrpPrcL1:    3500,
		GrpPrcL2:    3500,
		GrpPrcL3:    3500,
		GrpPrcL4:    3500,
		GrpPrcL5:    0,
		AddInfo:     "",
		SftyQty:     "",
		IsrcAplcbYn: "N",
		UseYn:       "Y",
		RegrNm:      "Admin",
		RegrId:      "Admin",
		ModrNm:      "Admin",
		ModrId:      "Admin",
	}

	stockRequest := StockRequest{
		Tin:    tin,
		BhfId:  bhfId,
		ItemCd: "KE1NTXU0000002",
		RsdQty: 10,
		RegrId: "Admin",
		RegrNm: "Admin",
		ModrNm: "Admin",
		ModrId: "Admin",
	}

	// Marshal request bodies to JSON
	codeRequestBody, err := json.Marshal(codeRequest)
	if err != nil {
		log.Fatal(err)
	}

	itemClassRequestBody, err := json.Marshal(itemClassRequest)
	if err != nil {
		log.Fatal(err)
	}

	customerRequestBody, err := json.Marshal(customerRequest)
	if err != nil {
		log.Fatal(err)
	}

	branchCustomerRequestBody, err := json.Marshal(branchCustomerRequest)
	if err != nil {
		log.Fatal(err)
	}

	itemRequestBody, err := json.Marshal(itemRequest)
	if err != nil {
		log.Fatal(err)
	}

	stockRequestBody, err := json.Marshal(stockRequest)
	if err != nil {
		log.Fatal(err)
	}

	// Send requests
	codeResponse, err := sendRequest(baseURL+"/code/selectCodes", headers, codeRequestBody)
	if err != nil {
		log.Fatal(err)
	}

	itemClassResponse, err := sendRequest(baseURL+"/itemClass/selectItemsClass", headers, itemClassRequestBody)
	if err != nil {
		log.Fatal(err)
	}

	customerResponse, err := sendRequest(baseURL+"/customers/selectCustomer", headers, customerRequestBody)
	if err != nil {
		log.Fatal(err)
	}

	branchCustomerResponse, err := sendRequest(baseURL+"/branches/saveBrancheCustomers", headers, branchCustomerRequestBody)
	if err != nil {
		log.Fatal(err)
	}

	itemResponse, err := sendRequest(baseURL+"/items/selectItems", headers, itemRequestBody)
	if err != nil {
		log.Fatal(err)
	}

	stockResponse, err := sendRequest(baseURL+"/stock/saveStock", headers, stockRequestBody)
	if err != nil {
		log.Fatal(err)
	}

	// Log responses
	log.Println(codeResponse)
	log.Println(itemClassResponse)
	log.Println(customerResponse)
	log.Println(branchCustomerResponse)
	log.Println(itemResponse)
	log.Println(stockResponse)
}

func sendRequest(url string, headers map[string]string, requestBody []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Print the response status and body
	log.Printf("Response Status: %d\n", resp.StatusCode)
	log.Printf("Response Body: %s\n", string(body))

	return resp, nil
}
