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

// StockQueryRequest for querying stock movements
type StockQueryRequest struct {
	models.BaseRequest
	LastReqDt string `json:"lastReqDt"`
}

// CreateStockMovementFromSale creates a stock movement record from a sales transaction
func (c *VSCUClient) CreateStockMovementFromSale(saleData []byte) (*StockMovementRecord, error) {
	var sale models.SalesRequest
	if err := json.Unmarshal(saleData, &sale); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sale data: %w", err)
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

	// Format current time for ocrnDt and stockRlsDt
	currentTime := time.Now()
	ocrnDt := currentTime.Format("20060102")
	stockRlsDt := currentTime.Format("20060102150405") // YYYYMMDDhhmmss format

	// Create new stock movement record for sale (stock out)
	stockMovement := &StockMovementRecord{
		Tin:         sale.Tin,
		BhfId:       sale.BhfId,
		SarNo:       lastSarNo + 1,
		OrgSarNo:    0,
		RegTyCd:     "M", // M for Movement
		CustTin:     &sale.Receipt.CustTin,
		CustNm:      &sale.Receipt.TrdeNm,
		CustBhfId:   nil,
		SarTyCd:     "02", // 02 for Stock Out from sale
		OcrnDt:      ocrnDt,
		StockRlsDt:  stockRlsDt,
		TotItemCnt:  sale.TotItemCnt,
		TotTaxblAmt: sale.TotTaxblAmt,
		TotTaxAmt:   sale.TotTaxAmt,
		TotAmt:      sale.TotAmt,
		Remark:      fmt.Sprintf("Stock out from sale invoice #%s", sale.InvcNo),
		RegrId:      sale.RegrId,
		RegrNm:      sale.RegrNm,
		ModrId:      sale.ModrId,
		ModrNm:      sale.ModrNm,
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
			Qty:        -math.Abs(float64(item.Qty)), // Always negative for sales
			StockTyCd:  "02", // 02 for Stock Out from sale
			ItemExprDt: stockRlsDt[:8], // Use the same date as stockRlsDt
			Prc:        item.Prc,
			SplyAmt:    item.SplyAmt,
			DcAmt:      item.DcAmt,
			DcRt:       item.DcRt,
			TotDcAmt:   item.DcAmt,
			TaxblAmt:   item.TaxblAmt,
			TaxTyCd:    item.TaxTyCd,
			TaxAmt:     item.TaxAmt,
			TotAmt:     item.TotAmt,
			Remark:     fmt.Sprintf("Stock out from sale invoice #%s", sale.InvcNo),
		}
		stockMovement.ItemList = append(stockMovement.ItemList, stockItem)
	}

	// Save the stock movement
	if err := c.SaveStockMovement(stockMovement); err != nil {
		return nil, fmt.Errorf("failed to save stock movement: %w", err)
	}

	// Update stock master
	if err := c.UpdateStockMaster(stockMovement); err != nil {
		return nil, fmt.Errorf("failed to update stock master: %w", err)
	}

	return stockMovement, nil
}

// CreateStockMovementFromPurchase creates a stock movement record from a purchase transaction
func (c *VSCUClient) CreateStockMovementFromPurchase(purchaseData []byte) (*StockMovementRecord, error) {
	var purchase models.PurchaseRequest
	if err := json.Unmarshal(purchaseData, &purchase); err != nil {
		return nil, fmt.Errorf("failed to unmarshal purchase data: %w", err)
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

	// Format current time for ocrnDt and stockRlsDt
	currentTime := time.Now()
	ocrnDt := currentTime.Format("20060102")
	stockRlsDt := currentTime.Format("20060102150405") // YYYYMMDDhhmmss format

	// Create new stock movement record for purchase (stock in)
	stockMovement := &StockMovementRecord{
		Tin:         purchase.Tin,
		BhfId:       purchase.BhfId,
		SarNo:       lastSarNo + 1,
		OrgSarNo:    0,
		RegTyCd:     "M", // M for Movement
		CustTin:     &purchase.SpplrTin,
		CustNm:      &purchase.SpplrNm,
		CustBhfId:   &purchase.SpplrBhfId,
		SarTyCd:     "01", // 01 for Stock In from purchase
		OcrnDt:      ocrnDt,
		StockRlsDt:  stockRlsDt,
		TotItemCnt:  purchase.TotItemCnt,
		TotTaxblAmt: purchase.TotTaxblAmt,
		TotTaxAmt:   purchase.TotTaxAmt,
		TotAmt:      purchase.TotAmt,
		Remark:      fmt.Sprintf("Stock in from purchase invoice #%s", purchase.InvcNo),
		RegrId:      purchase.RegrId,
		RegrNm:      purchase.RegrNm,
		ModrId:      purchase.ModrId,
		ModrNm:      purchase.ModrNm,
		ItemList:    make([]models.StockItem, 0),
	}

	// Convert purchase items to stock items
	for _, item := range purchase.ItemList {
		stockItem := models.StockItem{
			ItemSeq:    item.ItemSeq,
			ItemCd:     item.ItemCd,
			ItemClsCd:  item.ItemClsCd,
			ItemNm:     item.ItemNm,
			PkgUnitCd:  item.PkgUnitCd,
			Pkg:        item.Pkg,
			QtyUnitCd:  item.QtyUnitCd,
			Qty:        math.Abs(float64(item.Qty)), // Always positive for purchases
			StockTyCd:  "01", // 01 for Stock In from purchase
			ItemExprDt: stockRlsDt[:8], // Use the same date as stockRlsDt
			Prc:        item.Prc,
			SplyAmt:    item.SplyAmt,
			DcAmt:      item.DcAmt,
			DcRt:       item.DcRt,
			TotDcAmt:   item.DcAmt,
			TaxblAmt:   item.TaxblAmt,
			TaxTyCd:    item.TaxTyCd,
			TaxAmt:     item.TaxAmt,
			TotAmt:     item.TotAmt,
			Remark:     fmt.Sprintf("Stock in from purchase invoice #%s", purchase.InvcNo),
		}
		stockMovement.ItemList = append(stockMovement.ItemList, stockItem)
	}

	// Save the stock movement
	if err := c.SaveStockMovement(stockMovement); err != nil {
		return nil, fmt.Errorf("failed to save stock movement: %w", err)
	}

	// Update stock master
	if err := c.UpdateStockMaster(stockMovement); err != nil {
		return nil, fmt.Errorf("failed to update stock master: %w", err)
	}

	return stockMovement, nil
}

// UpdateStockMaster updates the stock master record after a stock movement
func (c *VSCUClient) UpdateStockMaster(movement *StockMovementRecord) error {
    // For each item in the movement, send a separate stock master update
    for _, item := range movement.ItemList {
        stockMaster := struct {
            Tin     string  `json:"tin"`
            BhfId   string  `json:"bhfId"`
            ItemCd  string  `json:"itemCd"`
            RsdQty  float64 `json:"rsdQty"`
            RegrId  string  `json:"regrId"`
            RegrNm  string  `json:"regrNm"`
            ModrId  string  `json:"modrId"`
            ModrNm  string  `json:"modrNm"`
        }{
            Tin:    movement.Tin,
            BhfId:  movement.BhfId,
            ItemCd: item.ItemCd,
            RsdQty: item.Qty,  // This will be negative for stock out
            RegrId: movement.RegrId,
            RegrNm: movement.RegrNm,
            ModrId: movement.ModrId,
            ModrNm: movement.ModrNm,
        }

        // Send request to update stock master
        resp, err := c.SendRequest("POST", "/stockMaster/saveStockMaster", stockMaster, nil)
        if err != nil {
            return fmt.Errorf("failed to update stock master: %w", err)
        }

        if resp.ResultCd != "000" {
            return fmt.Errorf("failed to update stock master: %s", resp.ResultMsg)
        }
    }

    return nil
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
    // Query stock movements to check current levels
    stockQuery := StockQueryRequest{
        BaseRequest: models.BaseRequest{
            Tin:   sale.Tin,
            BhfId: sale.BhfId,
        },
        LastReqDt: time.Now().Format("20060102150405"),
    }

    // Get stock movements
    resp, err := c.SendRequest("POST", "/stock/selectStockItems", stockQuery, nil)
    if err != nil {
        return fmt.Errorf("failed to query stock movements: %w", err)
    }

    if resp.ResultCd != "000" {
        return fmt.Errorf("failed to query stock movements: %s", resp.ResultMsg)
    }

    // Parse the response to get current stock levels
    var stockResp struct {
        Data struct {
            StockList []struct {
                ItemCd string  `json:"itemCd"`
                Qty    float64 `json:"qty"`
            } `json:"stockList"`
        } `json:"data"`
    }
    
    // Marshal and unmarshal to handle interface{} data
    respBytes, err := json.Marshal(resp.Data)
    if err != nil {
        return fmt.Errorf("failed to marshal stock response: %w", err)
    }
    
    if err := json.Unmarshal(respBytes, &stockResp); err != nil {
        return fmt.Errorf("failed to parse stock response: %w", err)
    }

    // Create a map of current stock levels
    stockLevels := make(map[string]float64)
    for _, item := range stockResp.Data.StockList {
        stockLevels[item.ItemCd] = item.Qty
    }

    // Validate each item in the sale
    for _, item := range sale.ItemList {
        currentStock, exists := stockLevels[item.ItemCd]
        if !exists {
            return fmt.Errorf("item %s not found in stock", item.ItemCd)
        }
        if currentStock < item.Qty {
            return fmt.Errorf("insufficient stock for item %s: have %.2f, need %.2f", 
                item.ItemCd, currentStock, item.Qty)
        }
    }

    return nil
}
