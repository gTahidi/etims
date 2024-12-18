package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"time"

	"etims-test/models"
)

// StockMovementRecord represents a stock movement transaction
type StockMovementRecord struct {
	Tin         string           `json:"tin"`
	BhfId       string           `json:"bhfId"`
	SarNo       int              `json:"sarNo"`
	OrgSarNo    int              `json:"orgSarNo"`
	RegTyCd     string           `json:"regTyCd"`
	CustTin     *string          `json:"custTin,omitempty"`
	CustNm      *string          `json:"custNm,omitempty"`
	CustBhfId   *string          `json:"custBhfId,omitempty"`
	SarTyCd     string           `json:"sarTyCd"`
	OcrnDt      string           `json:"ocrnDt"`
	StockRlsDt  string           `json:"stockRlsDt"`
	TotItemCnt  int              `json:"totItemCnt"`
	TotTaxblAmt float64          `json:"totTaxblAmt"`
	TotTaxAmt   float64          `json:"totTaxAmt"`
	TotAmt      float64          `json:"totAmt"`
	Remark      string           `json:"remark"`
	RegrId      string           `json:"regrId"`
	RegrNm      string           `json:"regrNm"`
	ModrId      string           `json:"modrId"`
	ModrNm      string           `json:"modrNm"`
	ItemList    []models.StockItem `json:"itemList"`
}

// StockInventory represents the current stock levels
type StockInventory struct {
	ItemCd string  `json:"itemCd"`
	Qty    float64 `json:"qty"`
}

// CreateStockMovementFromSale creates a stock movement record from a sales transaction
func (c *VSCUClient) CreateStockMovementFromSale(salesData []byte) (*StockMovementRecord, error) {
	var sale models.SalesRequest
	if err := json.Unmarshal(salesData, &sale); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sales data: %w", err)
	}

	// Validate stock availability before proceeding
	if err := c.ValidateStockAvailability(sale); err != nil {
		return nil, fmt.Errorf("stock validation failed: %w", err)
	}

	// Get the latest SarNo from stock.json
	stockFile := filepath.Join("testdata", "stock.json")
	stockData, err := ioutil.ReadFile(stockFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read stock file: %w", err)
	}

	var lastSarNo int
	if len(stockData) > 0 {
		var movements []StockMovementRecord
		if err := json.Unmarshal(stockData, &movements); err != nil {
			// Try unmarshaling as single record if array fails
			var singleMovement StockMovementRecord
			if err2 := json.Unmarshal(stockData, &singleMovement); err2 != nil {
				return nil, fmt.Errorf("failed to parse stock data: %w", err)
			}
			lastSarNo = singleMovement.SarNo
		} else if len(movements) > 0 {
			lastSarNo = movements[len(movements)-1].SarNo
		}
	}

	// Convert customer details
	var custTin, custNm, custBhfId *string
	if sale.CustTin != "" {
		custTin = &sale.CustTin
	}
	if sale.CustNm != "" {
		custNm = &sale.CustNm
	}
	// BhfId is not in the sales data, so we'll leave it as nil

	// Format current time for ocrnDt and stockRlsDt
	currentTime := time.Now()
	ocrnDt := currentTime.Format("20060102")
	stockRlsDt := currentTime.Format("20060102150405") // YYYYMMDDhhmmss format

	// Create new stock movement record
	stockMovement := &StockMovementRecord{
		Tin:         sale.Tin,
		BhfId:       sale.BhfId,
		SarNo:       lastSarNo + 1,
		OrgSarNo:    0,
		RegTyCd:     "M", // M for Movement (sales)
		CustTin:     custTin,
		CustNm:      custNm,
		CustBhfId:   custBhfId,
		SarTyCd:     "02", // 02 for Stock Out
		OcrnDt:      ocrnDt,
		StockRlsDt:  stockRlsDt,
		TotItemCnt:  sale.TotItemCnt,
		TotTaxblAmt: sale.TotTaxblAmt,
		TotTaxAmt:   sale.TotTaxAmt,
		TotAmt:      sale.TotAmt,
		Remark:      fmt.Sprintf("Stock out for sales invoice #%s", sale.InvcNo),
		RegrId:      "Admin",
		RegrNm:      "Admin",
		ModrId:      "Admin",
		ModrNm:      "Admin",
		ItemList:    make([]models.StockItem, 0),
	}

	// Convert sale items to stock items
	for _, item := range sale.ItemList {
		stockItem := models.StockItem{
			ItemSeq:    item.ItemSeq,
			ItemCd:     item.ItemCd,
			ItemClsCd:  item.ItemClsCd,
			ItemNm:     item.ItemNm,
			PkgUnitCd:  item.PkgUnitCd,
			Pkg:        item.Pkg,
			QtyUnitCd:  item.QtyUnitCd,
			Qty:        -math.Abs(float64(item.Qty)), // Make quantity negative for stock out
			StockTyCd:  "02", // 02 for Stock Out
			ItemExprDt: stockRlsDt[:8], // Use the same date as stockRlsDt
			Prc:        item.Prc,
			SplyAmt:    item.SplyAmt,
			DcAmt:      item.DcAmt,
			DcRt:       item.DcRt,
			TotDcAmt:   item.DcAmt, // Using DcAmt as TotDcAmt
			TaxblAmt:   item.TaxblAmt,
			TaxTyCd:    item.TaxTyCd,
			TaxAmt:     item.TaxAmt,
			TotAmt:     item.TotAmt,
			Remark:     fmt.Sprintf("Stock out for sales invoice #%s", sale.InvcNo),
		}
		stockMovement.ItemList = append(stockMovement.ItemList, stockItem)
	}

	// Save the stock movement
	if err := c.SaveStockMovement(stockMovement); err != nil {
		return nil, fmt.Errorf("failed to save stock movement: %w", err)
	}

	return stockMovement, nil
}

// SaveStockMovement saves the stock movement record to the stock.json file
func (c *VSCUClient) SaveStockMovement(movement *StockMovementRecord) error {
	stockFile := filepath.Join("testdata", "stock.json")
	
	// Load existing movements
	movements, err := c.LoadStockMovements()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load existing stock movements: %w", err)
	}

	// Append new movement
	movements = append(movements, *movement)

	// Save back to file
	stockData, err := json.MarshalIndent(movements, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal stock movements: %w", err)
	}

	if err := ioutil.WriteFile(stockFile, stockData, 0644); err != nil {
		return fmt.Errorf("failed to write stock file: %w", err)
	}

	return nil
}

// LoadStockMovements loads all stock movements from the stock.json file
func (c *VSCUClient) LoadStockMovements() ([]StockMovementRecord, error) {
	stockFile := filepath.Join("testdata", "stock.json")
	stockData, err := ioutil.ReadFile(stockFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []StockMovementRecord{}, nil
		}
		return nil, fmt.Errorf("failed to read stock file: %w", err)
	}

	if len(stockData) == 0 {
		return []StockMovementRecord{}, nil
	}

	// Try to unmarshal as array first
	var movements []StockMovementRecord
	if err := json.Unmarshal(stockData, &movements); err == nil {
		return movements, nil
	}

	// If array unmarshal fails, try as single object
	var singleMovement StockMovementRecord
	if err := json.Unmarshal(stockData, &singleMovement); err != nil {
		return nil, fmt.Errorf("failed to parse stock data: %w", err)
	}

	// Return single movement as slice
	return []StockMovementRecord{singleMovement}, nil
}

// GetCurrentStock gets the current stock level for an item
func (c *VSCUClient) GetCurrentStock(itemCd string) (float64, error) {
	movements, err := c.LoadStockMovements()
	if err != nil {
		return 0, fmt.Errorf("failed to load stock movements: %w", err)
	}

	var currentQty float64
	for _, mov := range movements {
		for _, item := range mov.ItemList {
			if item.ItemCd == itemCd {
				currentQty += item.Qty // Will add for stock in, subtract for stock out (due to negative quantities)
			}
		}
	}

	return currentQty, nil
}

// ValidateStockAvailability checks if we have enough stock for a sale
func (c *VSCUClient) ValidateStockAvailability(sale models.SalesRequest) error {
	for _, item := range sale.ItemList {
		currentStock, err := c.GetCurrentStock(item.ItemCd)
		if err != nil {
			return fmt.Errorf("failed to get current stock for item %s: %w", item.ItemCd, err)
		}

		if currentStock < float64(item.Qty) {
			return fmt.Errorf("insufficient stock for item %s: have %.2f, need %.2f", 
				item.ItemCd, currentStock, float64(item.Qty))
		}
	}
	return nil
}
