package models

import (
	"fmt"
	_ "strings"
)

// Validator interface for request validation
type Validator interface {
	Validate() error
}

// BaseRequest contains common fields for all requests
type BaseRequest struct {
	Tin   string `json:"tin"`
	BhfId string `json:"bhfId"`
}

// APIResponse is the common response structure for all API calls
type APIResponse struct {
	ResultCd  string      `json:"resultCd"`
	ResultMsg string      `json:"resultMsg"`
	Data      interface{} `json:"data,omitempty"`
}

// CodeRequest for getting code lists
type CodeRequest struct {
	BaseRequest
	LastReqDt string `json:"lastReqDt"`
}

func (r CodeRequest) Validate() error {
	if r.Tin == "" {
		return fmt.Errorf("tin is required")
	}
	if r.BhfId == "" {
		return fmt.Errorf("bhfId is required")
	}
	if r.LastReqDt == "" {
		return fmt.Errorf("lastReqDt is required")
	}
	return nil
}

// ItemClassRequest for getting item classifications
type ItemClassRequest struct {
	BaseRequest
	LastReqDt string `json:"lastReqDt"`
}

func (r ItemClassRequest) Validate() error {
	if r.Tin == "" {
		return fmt.Errorf("tin is required")
	}
	if r.BhfId == "" {
		return fmt.Errorf("bhfId is required")
	}
	if r.LastReqDt == "" {
		return fmt.Errorf("lastReqDt is required")
	}
	return nil
}

// CustomerRequest for getting customer information
type CustomerRequest struct {
	BaseRequest
	CustmTin string `json:"custmTin"`
}

func (r CustomerRequest) Validate() error {
	if r.Tin == "" {
		return fmt.Errorf("tin is required")
	}
	if r.BhfId == "" {
		return fmt.Errorf("bhfId is required")
	}
	if r.CustmTin == "" {
		return fmt.Errorf("custmTin is required")
	}
	return nil
}

// BranchCustomerRequest for managing branch customers
type BranchCustomerRequest struct {
	BaseRequest
	CustNo  string `json:"custNo"`
	CustTin string `json:"custTin"`
	CustNm  string `json:"custNm"`
	Adrs    string `json:"adrs,omitempty"`
	TelNo   string `json:"telNo,omitempty"`
	Email   string `json:"email,omitempty"`
	FaxNo   string `json:"faxNo,omitempty"`
	UseYn   string `json:"useYn"`
	Remark  string `json:"remark,omitempty"`
	RegrNm  string `json:"regrNm"`
	RegrId  string `json:"regrId"`
	ModrNm  string `json:"modrNm"`
	ModrId  string `json:"modrId"`
}

func (r BranchCustomerRequest) Validate() error {
	if r.Tin == "" {
		return fmt.Errorf("tin is required")
	}
	if r.BhfId == "" {
		return fmt.Errorf("bhfId is required")
	}
	if r.CustNo == "" {
		return fmt.Errorf("custNo is required")
	}
	if r.CustTin == "" {
		return fmt.Errorf("custTin is required")
	}
	if r.CustNm == "" {
		return fmt.Errorf("custNm is required")
	}
	if r.UseYn == "" {
		return fmt.Errorf("useYn is required")
	}
	if r.RegrNm == "" {
		return fmt.Errorf("regrNm is required")
	}
	if r.RegrId == "" {
		return fmt.Errorf("regrId is required")
	}
	if r.ModrNm == "" {
		return fmt.Errorf("modrNm is required")
	}
	if r.ModrId == "" {
		return fmt.Errorf("modrId is required")
	}
	return nil
}

// ItemRequest for managing items
type ItemRequest struct {
	BaseRequest
	ItemCd      string  `json:"itemCd"`
	ItemClsCd   string  `json:"itemClsCd"`
	ItemTyCd    string  `json:"itemTyCd"`
	ItemNm      string  `json:"itemNm"`
	ItemStdNm   string  `json:"itemStdNm"`
	OrgnNatCd   string  `json:"orgnNatCd"`
	PkgUnitCd   string  `json:"pkgUnitCd"`
	QtyUnitCd   string  `json:"qtyUnitCd"`
	TaxTyCd     string  `json:"taxTyCd"`
	BcdRate     float64 `json:"bcdRate"`
	DtyPrc      float64 `json:"dtyPrc"`
	DftPrc      float64 `json:"dftPrc"`
	Regulator   string  `json:"regulator"`
	BatchNo     string  `json:"batchNo"`
	ItemExprDt  string  `json:"itemExprDt"`
	ItemSrno    string  `json:"itemSrno"`
	GrpPrcL1    float64 `json:"grpPrcL1"`
	GrpPrcL2    float64 `json:"grpPrcL2"`
	GrpPrcL3    float64 `json:"grpPrcL3"`
	GrpPrcL4    float64 `json:"grpPrcL4"`
	AddInfo     string  `json:"addInfo"`
	UseYn       string  `json:"useYn"`
	IsrcAplcbYn string  `json:"isrcAplcbYn"`
	RegrId      string  `json:"regrId"`
	RegrNm      string  `json:"regrNm"`
	ModrId      string  `json:"modrId"`
	ModrNm      string  `json:"modrNm"`
}

func (r ItemRequest) Validate() error {
	if r.Tin == "" {
		return fmt.Errorf("tin is required")
	}
	if r.BhfId == "" {
		return fmt.Errorf("bhfId is required")
	}
	if r.ItemCd == "" {
		return fmt.Errorf("itemCd is required")
	}
	if r.ItemClsCd == "" {
		return fmt.Errorf("itemClsCd is required")
	}
	if r.ItemTyCd == "" {
		return fmt.Errorf("itemTyCd is required")
	}
	if r.ItemNm == "" {
		return fmt.Errorf("itemNm is required")
	}
	if r.UseYn == "" {
		return fmt.Errorf("useYn is required")
	}
	if r.IsrcAplcbYn == "" {
		return fmt.Errorf("isrcAplcbYn is required")
	}
	if r.RegrId == "" {
		return fmt.Errorf("regrId is required")
	}
	if r.RegrNm == "" {
		return fmt.Errorf("regrNm is required")
	}
	if r.ModrId == "" {
		return fmt.Errorf("modrId is required")
	}
	if r.ModrNm == "" {
		return fmt.Errorf("modrNm is required")
	}
	return nil
}

// StockItem represents an item in a stock transaction
type StockItem struct {
	ItemSeq    int     `json:"itemSeq"`
	ItemCd     string  `json:"itemCd"`
	ItemClsCd  string  `json:"itemClsCd"`
	ItemNm     string  `json:"itemNm"`
	PkgUnitCd  string  `json:"pkgUnitCd"`
	Pkg        int     `json:"pkg"`
	QtyUnitCd  string  `json:"qtyUnitCd"`
	Qty        int     `json:"qty"`
	Prc        float64 `json:"prc"`
	SplyAmt    float64 `json:"splyAmt"`
	DcRt       float64 `json:"dcRt"`
	DcAmt      float64 `json:"dcAmt"`
	TotDcAmt   float64 `json:"totDcAmt"`
	TaxTyCd    string  `json:"taxTyCd"`
	TaxblAmt   float64 `json:"taxblAmt"`
	TaxAmt     float64 `json:"taxAmt"`
	TotAmt     float64 `json:"totAmt"`
	StockTyCd  string  `json:"stockTyCd"`
	ItemExprDt string  `json:"itemExprDt"`
	Remark     string  `json:"remark"`
}

// StockRequest for managing stock
type StockRequest struct {
	BaseRequest
	SarNo       int         `json:"sarNo"`
	OrgSarNo    int         `json:"orgSarNo"`
	RegTyCd     string      `json:"regTyCd"`
	CustTin     string      `json:"custTin"`
	CustNm      string      `json:"custNm"`
	CustBhfId   string      `json:"custBhfId"`
	SarTyCd     string      `json:"sarTyCd"`
	OcrnDt      string      `json:"ocrnDt"`
	StockRlsDt  string      `json:"stockRlsDt"`
	TotItemCnt  int         `json:"totItemCnt"`
	TotTaxblAmt float64     `json:"totTaxblAmt"`
	TotTaxAmt   float64     `json:"totTaxAmt"`
	TotAmt      float64     `json:"totAmt"`
	Remark      string      `json:"remark"`
	RegrId      string      `json:"regrId"`
	RegrNm      string      `json:"regrNm"`
	ModrId      string      `json:"modrId"`
	ModrNm      string      `json:"modrNm"`
	ItemList    []StockItem `json:"itemList"`
}

func (r StockRequest) Validate() error {
	if r.Tin == "" {
		return fmt.Errorf("tin is required")
	}
	if r.BhfId == "" {
		return fmt.Errorf("bhfId is required")
	}
	if len(r.ItemList) == 0 {
		return fmt.Errorf("at least one stock item is required")
	}
	for i, item := range r.ItemList {
		if item.ItemCd == "" {
			return fmt.Errorf("itemCd is required for stock item %d", i+1)
		}
		if item.ItemNm == "" {
			return fmt.Errorf("itemNm is required for stock item %d", i+1)
		}
		if item.StockTyCd == "" {
			return fmt.Errorf("stockTyCd is required for stock item %d", i+1)
		}
	}
	return nil
}

// SalesRequest for managing sales
type SalesRequest struct {
	BaseRequest
	InvcNo       string      `json:"invcNo"`
	OrgInvcNo    string      `json:"orgInvcNo"`
	CustTin      string      `json:"custTin"`
	CustNm       string      `json:"custNm"`
	SalesTyCd    string      `json:"salesTyCd"`
	RcptTyCd     string      `json:"rcptTyCd"`
	PmtTyCd      string      `json:"pmtTyCd"`
	SalesSttsCd  string      `json:"salesSttsCd"`
	CfmDt        string      `json:"cfmDt"`
	SalesDt      string      `json:"salesDt"`
	StockRlsDt   string      `json:"stockRlsDt"`
	CnclReqDt    string      `json:"cnclReqDt"`
	CnclDt       string      `json:"cnclDt"`
	RfdDt        string      `json:"rfdDt"`
	RfdRsnCd     string      `json:"rfdRsnCd"`
	TotItemCnt   int         `json:"totItemCnt"`
	TaxblAmtA    float64     `json:"taxblAmtA"`
	TaxblAmtB    float64     `json:"taxblAmtB"`
	TaxblAmtC    float64     `json:"taxblAmtC"`
	TaxblAmtD    float64     `json:"taxblAmtD"`
	TaxblAmtE    float64     `json:"taxblAmtE"`
	TaxRtA       float64     `json:"taxRtA"`
	TaxRtB       float64     `json:"taxRtB"`
	TaxRtC       float64     `json:"taxRtC"`
	TaxRtD       float64     `json:"taxRtD"`
	TaxRtE       float64     `json:"taxRtE"`
	TaxAmtA      float64     `json:"taxAmtA"`
	TaxAmtB      float64     `json:"taxAmtB"`
	TaxAmtC      float64     `json:"taxAmtC"`
	TaxAmtD      float64     `json:"taxAmtD"`
	TaxAmtE      float64     `json:"taxAmtE"`
	TotTaxblAmt  float64     `json:"totTaxblAmt"`
	TotTaxAmt    float64     `json:"totTaxAmt"`
	TotAmt       float64     `json:"totAmt"`
	PrchrAcptcYn string      `json:"prchrAcptcYn"`
	Remark       string      `json:"remark"`
	RegrId       string      `json:"regrId"`
	RegrNm       string      `json:"regrNm"`
	ModrId       string      `json:"modrId"`
	ModrNm       string      `json:"modrNm"`
	Receipt      Receipt     `json:"receipt"`
	ItemList     []SalesItem `json:"itemList"`
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
	Bcd       string  `json:"bcd"`
	PkgUnitCd string  `json:"pkgUnitCd"`
	Pkg       int     `json:"pkg"`
	QtyUnitCd string  `json:"qtyUnitCd"`
	Qty       int     `json:"qty"`
	Prc       float64 `json:"prc"`
	SplyAmt   float64 `json:"splyAmt"`
	DcRt      float64 `json:"dcRt"`
	DcAmt     float64 `json:"dcAmt"`
	IsrccCd   string  `json:"isrccCd"`
	IsrccNm   string  `json:"isrccNm"`
	IsrcRt    float64 `json:"isrcRt"`
	IsrcAmt   float64 `json:"isrcAmt"`
	TaxTyCd   string  `json:"taxTyCd"`
	TaxblAmt  float64 `json:"taxblAmt"`
	TaxAmt    float64 `json:"taxAmt"`
	TotAmt    float64 `json:"totAmt"`
}

func (r SalesRequest) Validate() error {
	if r.Tin == "" {
		return fmt.Errorf("tin is required")
	}
	if r.BhfId == "" {
		return fmt.Errorf("bhfId is required")
	}
	if r.SalesTyCd == "" {
		return fmt.Errorf("salesTyCd is required")
	}
	if r.RcptTyCd == "" {
		return fmt.Errorf("rcptTyCd is required")
	}
	if r.PmtTyCd == "" {
		return fmt.Errorf("pmtTyCd is required")
	}
	if r.SalesDt == "" {
		return fmt.Errorf("salesDt is required")
	}
	if len(r.ItemList) == 0 {
		return fmt.Errorf("at least one item is required")
	}
	return nil
}

type CodeListRequest struct {
	Tin       string `json:"tin"`
	BhfId     string `json:"bhfId"`
	LastReqDt string `json:"lastReqDt"`
}

// InitRequest represents a device initialization request
type InitRequest struct {
	Tin      string `json:"tin"`
	BhfId    string `json:"bhfId"`
	DvcSrlNo string `json:"dvcSrlNo"` // Device Serial Number
}

// BranchListRequest represents a request to get branch list
type BranchListRequest struct {
	Tin       string `json:"tin"`
	BhfId     string `json:"bhfId"`
	LastReqDt string `json:"lastReqDt"`
}

// NoticeRequest represents a request to get notices
type NoticeRequest struct {
	Tin       string `json:"tin"`
	BhfId     string `json:"bhfId"`
	LastReqDt string `json:"lastReqDt"`
}

// ItemCompositionRequest represents a request to save item composition
type ItemCompositionRequest struct {
	Tin       string `json:"tin"`
	BhfId     string `json:"bhfId"`
	TaskCd    string `json:"taskCd"`    // Task Code
	ItemCd    string `json:"itemCd"`    // Item Code
	ItemClsCd string `json:"itemClsCd"` // Item Classification Code
	HsCd      string `json:"hsCd"`      // HS Code
	ItemSeq   int    `json:"itemSeq"`   // Item Sequence
	Remark    string `json:"remark,omitempty"`
	ModrNm    string `json:"modrNm"` // Modifier Name
	ModrId    string `json:"modrId"` // Modifier ID
}

// ImportedItemRequest represents a request to get imported items
type ImportedItemRequest struct {
	Tin       string `json:"tin"`
	BhfId     string `json:"bhfId"`
	LastReqDt string `json:"lastReqDt"`
}

// ImportedItemSaveRequest represents a request to save imported items
type ImportedItemSaveRequest struct {
	Tin            string `json:"tin"`
	BhfId          string `json:"bhfId"`
	TaskCd         string `json:"taskCd"`         // Task Code
	DclDe          string `json:"dclDe"`          // Declaration Date
	ItemSeq        int    `json:"itemSeq"`        // Item Sequence
	HsCd           string `json:"hsCd"`           // HS Code
	ItemClsCd      string `json:"itemClsCd"`      // Item Classification Code
	ItemCd         string `json:"itemCd"`         // Item Code
	ImptItemSttsCd string `json:"imptItemSttsCd"` // Import Item Status Code
	Remark         string `json:"remark,omitempty"`
	ModrNm         string `json:"modrNm"` // Modifier Name
	ModrId         string `json:"modrId"` // Modifier ID
}

// PurchaseConfirmationRequest represents a request to confirm purchase transactions
type PurchaseConfirmationRequest struct {
	Tin        string `json:"tin"`
	BhfId      string `json:"bhfId"`
	PurchaseNo string `json:"purchaseNo"` // Purchase Number
	ConfirmDt  string `json:"confirmDt"`  // Confirmation Date
	Remark     string `json:"remark,omitempty"`
	ModrNm     string `json:"modrNm"` // Modifier Name
	ModrId     string `json:"modrId"` // Modifier ID
}

// CodeSaveRequest represents a request to save codes
type CodeSaveRequest struct {
	Tin   string `json:"tin"`
	BhfId string `json:"bhfId"`
	Codes []Code `json:"codes"`
}

// Code represents a code entry
type Code struct {
	GroupCd string `json:"groupCd"`
	GroupNm string `json:"groupNm"`
	CodeCd  string `json:"codeCd"`
	CodeNm  string `json:"codeNm"`
	UseYn   string `json:"useYn"`
}

// ItemClassSaveRequest represents a request to save item classifications
type ItemClassSaveRequest struct {
	Tin         string      `json:"tin"`
	BhfId       string      `json:"bhfId"`
	ItemClasses []ItemClass `json:"itemClasses"`
}

// ItemClass represents an item classification
type ItemClass struct {
	ItemClsCd  string `json:"itemClsCd"`
	ItemClsNm  string `json:"itemClsNm"`
	ItemClsLvl int    `json:"itemClsLvl"`
	UseYn      string `json:"useYn"`
	VatYn      string `json:"vatYn"`
}
