package testmanager

import (
	"fmt"
	"math/rand"
	"time"
)

// DataGenerator handles the generation of test data
type DataGenerator struct {
	baseData map[string]interface{}
}

// NewDataGenerator creates a new instance of DataGenerator
func NewDataGenerator(baseData map[string]interface{}) *DataGenerator {
	return &DataGenerator{
		baseData: baseData,
	}
}

// GenerateCodeData creates test data for code lookup
func (g *DataGenerator) GenerateCodeData() map[string]interface{} {
	return map[string]interface{}{
		"tin":       g.baseData["tin"],
		"bhfId":     g.baseData["bhfId"],
		"lastReqDt": time.Now().Format("20060102150405"),
		"codes": []map[string]interface{}{
			{
				"groupCd": "15100000",
				"groupNm": "Fuels",
				"codeCd":  "01",
				"codeNm":  "Regular",
				"useYn":   "Y",
			},
		},
	}
}

// GenerateItemData creates test data for item management
func (g *DataGenerator) GenerateItemData() map[string]interface{} {
	now := time.Now()
	itemCd := fmt.Sprintf("KE1NTXU%07d", rand.Intn(1000000))

	return map[string]interface{}{
		"tin":   g.baseData["tin"],
		"bhfId": g.baseData["bhfId"],
		"items": []map[string]interface{}{
			{
				"itemCd":     itemCd,
				"itemClsCd":  "15101502",
				"itemTyCd":   "1",
				"itemNm":     "Test Fuel Item",
				"itemStdNm":  "Standard Test Fuel",
				"orgnNatCd":  "KE",
				"pkgUnitCd":  "JY",
				"qtyUnitCd":  "U",
				"taxTyCd":    "B",
				"bcdRate":    0,
				"dtyPrc":     150.75,
				"dftPrc":     150.75,
				"regulator":  "KRA",
				"batchNo":    fmt.Sprintf("BATCH%03d", rand.Intn(1000)),
				"itemExprDt": now.AddDate(1, 0, 0).Format("20060102"),
				"itemSrno":   fmt.Sprintf("SN%03d", rand.Intn(1000)),
			},
		},
	}
}

// GenerateBranchCustomerData creates test data for branch customer management
func (g *DataGenerator) GenerateBranchCustomerData() map[string]interface{} {
	custNo := fmt.Sprintf("CUST%03d", rand.Intn(1000))

	return map[string]interface{}{
		"tin":     g.baseData["tin"],
		"bhfId":   g.baseData["bhfId"],
		"custNo":  custNo,
		"custTin": fmt.Sprintf("A%09dX", rand.Intn(1000000000)),
		"custNm":  "Test Customer Ltd",
		"adrs":    "Test Address, Nairobi",
		"telNo":   fmt.Sprintf("+2547%08d", rand.Intn(100000000)),
		"email":   "test@example.com",
		"faxNo":   "",
		"useYn":   "Y",
		"remark":  "Test customer record",
		"regrNm":  "Admin",
		"regrId":  "Admin",
		"modrNm":  "Admin",
		"modrId":  "Admin",
	}
}

// GeneratePurchaseData creates test data for purchase transactions
func (g *DataGenerator) GeneratePurchaseData() map[string]interface{} {
	now := time.Now()
	invcNo := fmt.Sprintf("INV%d", now.Unix())

	return map[string]interface{}{
		"tin":         g.baseData["tin"],
		"bhfId":       g.baseData["bhfId"],
		"invcNo":      invcNo,
		"spplrTin":    fmt.Sprintf("S%010d", rand.Intn(1000000000)),
		"spplrBhfId":  "00",
		"spplrNm":     "Test Supplier",
		"pchsDt":      now.Format("20060102150405"),
		"wrhsDt":      now.Format("20060102150405"),
		"regTyCd":     "A",
		"pchsTyCd":    "N",
		"rcptTyCd":    "P",
		"pmtTyCd":     "C",
		"pchsSttsCd":  "02",
		"totItemCnt":  2,
		"taxblAmtA":   0,
		"taxblAmtB":   1000.00,
		"taxRtA":      0,
		"taxRtB":      16,
		"taxAmtA":     0,
		"taxAmtB":     160.00,
		"totTaxblAmt": 1000.00,
		"totTaxAmt":   160.00,
		"totAmt":      1160.00,
		"items": []map[string]interface{}{
			{
				"itemSeq":   1,
				"itemCd":    fmt.Sprintf("ITEM%d", rand.Intn(1000)),
				"itemClsCd": "15101502",
				"itemNm":    "Test Kerosene Item",
				"qty":       1,
				"prc":       500.00,
				"splyAmt":   500.00,
				"taxblAmt":  500.00,
				"taxAmt":    80.00,
				"totAmt":    580.00,
			},
			{
				"itemSeq":   2,
				"itemCd":    fmt.Sprintf("ITEM%d", rand.Intn(1000)),
				"itemClsCd": "15101502",
				"itemNm":    "Test Petrol Item 2",
				"qty":       1,
				"prc":       500.00,
				"splyAmt":   500.00,
				"taxblAmt":  500.00,
				"taxAmt":    80.00,
				"totAmt":    580.00,
			},
		},
	}
}

// GenerateStockData creates test data for stock movements
func (g *DataGenerator) GenerateStockData(purchaseData map[string]interface{}) map[string]interface{} {
	now := time.Now()
	sarNo := fmt.Sprintf("STK%d", now.Unix())

	purchaseItems, _ := purchaseData["items"].([]map[string]interface{})
	stockItems := make([]map[string]interface{}, 0)

	for _, item := range purchaseItems {
		stockItems = append(stockItems, map[string]interface{}{
			"itemCd":    item["itemCd"],
			"itemClsCd": item["itemClsCd"],
			"itemNm":    item["itemNm"],
			"qty":       item["qty"],
			"unitCd":    "U",
			"pkgUnitCd": "JY",
			"pkg":       1,
			"qtyUnitCd": "U",
			"prc":       item["prc"],
			"splyAmt":   item["splyAmt"],
			"totDcAmt":  0,
			"taxblAmt":  item["taxblAmt"],
			"taxTyCd":   "B",
			"taxAmt":    item["taxAmt"],
			"totAmt":    item["totAmt"],
		})
	}

	return map[string]interface{}{
		"tin":         g.baseData["tin"],
		"bhfId":       g.baseData["bhfId"],
		"sarNo":       sarNo,
		"orgSarNo":    "",
		"sarTyCd":     "01", // Stock increase from purchase
		"ocrnDt":      now.Format("20060102150405"),
		"totItemCnt":  len(stockItems),
		"totTaxblAmt": purchaseData["totTaxblAmt"],
		"totTaxAmt":   purchaseData["totTaxAmt"],
		"totAmt":      purchaseData["totAmt"],
		"items":       stockItems,
	}
}

// GenerateStockMasterData creates test data for stock master information
func (g *DataGenerator) GenerateStockMasterData() map[string]interface{} {
	now := time.Now()

	return map[string]interface{}{
		"tin":    g.baseData["tin"],
		"bhfId":  g.baseData["bhfId"],
		"itemCd": fmt.Sprintf("ITEM%d", rand.Intn(1000)),
		"rsdQty": 1000,
		"updDt":  now.Format("20060102150405"),
		"regrNm": "Admin",
		"regrId": "Admin",
		"modrNm": "Admin",
		"modrId": "Admin",
	}
}

// GenerateTestDataSet creates a complete set of related test data
func (g *DataGenerator) GenerateTestDataSet() map[string]interface{} {
	codeData := g.GenerateCodeData()
	itemData := g.GenerateItemData()
	customerData := g.GenerateBranchCustomerData()
	purchaseData := g.GeneratePurchaseData()
	stockData := g.GenerateStockData(purchaseData)
	stockMasterData := g.GenerateStockMasterData()

	return map[string]interface{}{
		"code":         codeData,
		"item":         itemData,
		"customer":     customerData,
		"purchase":     purchaseData,
		"stock":        stockData,
		"stock_master": stockMasterData,
	}
}
