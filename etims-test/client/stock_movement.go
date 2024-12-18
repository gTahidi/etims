package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	CustTin     *string          `json:"custTin"`
	CustNm      *string          `json:"custNm"`
	CustBhfId   *string          `json:"custBhfId"`
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

// CreateStockMovementFromSale creates a stock movement record from a sales transaction
func (c *VSCUClient) CreateStockMovementFromSale(salesData []byte) (*StockMovementRecord, error) {
	var sale models.SalesRequest
	if err := json.Unmarshal(salesData, &sale); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sales data: %w", err)
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
			return nil, fmt.Errorf("failed to parse stock data: %w", err)
		}
		if len(movements) > 0 {
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

	// Format current time for ocrnDt
	currentTime := time.Date(2024, 12, 18, 14, 53, 27, 0, time.UTC).Add(3 * time.Hour)
	ocrnDt := currentTime.Format("20060102")

	// Create new stock movement record
	stockMovement := &StockMovementRecord{
		Tin:         sale.Tin,
		BhfId:       sale.BhfId,
		SarNo:       lastSarNo + 1,
		OrgSarNo:    0,
		RegTyCd:     "M", // M for Movement
		CustTin:     custTin,
		CustNm:      custNm,
		CustBhfId:   custBhfId,
		SarTyCd:     "11", // 11 for Sales Out
		OcrnDt:      ocrnDt,
		StockRlsDt:  "", // This should be fetched from the original stock record
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
			Qty:        item.Qty,
			StockTyCd:  "11", // 11 for Sales Out
			ItemExprDt: "20251215", // This should be fetched from the original stock record
			Prc:        item.Prc,
			SplyAmt:    item.SplyAmt,
			DcAmt:      item.DcAmt,
			TotDcAmt:   item.DcAmt, // Using DcAmt as TotDcAmt
			TaxblAmt:   item.TaxblAmt,
			TaxTyCd:    item.TaxTyCd,
			TaxAmt:     item.TaxAmt,
			TotAmt:     item.TotAmt,
			Remark:     fmt.Sprintf("Stock out for sales invoice #%s", sale.InvcNo),
		}
		stockMovement.ItemList = append(stockMovement.ItemList, stockItem)
	}

	// Save the updated stock movements back to file
	var stockMovements []StockMovementRecord
	if len(stockData) > 0 {
		if err := json.Unmarshal(stockData, &stockMovements); err != nil {
			return nil, fmt.Errorf("failed to parse stock data: %w", err)
		}
	}
	stockMovements = append(stockMovements, *stockMovement)
	updatedStockData, err := json.MarshalIndent(stockMovements, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated stock data: %w", err)
	}

	if err := ioutil.WriteFile(stockFile, updatedStockData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write updated stock data: %w", err)
	}

	return stockMovement, nil
}

// SaveStockMovement saves the stock movement record to the stock.json file
func (c *VSCUClient) SaveStockMovement(movement *StockMovementRecord) error {
	// Load existing movements
	movements, err := c.LoadStockMovements()
	if err != nil {
		return err
	}

	// Add new movement
	movements = append(movements, *movement)

	// Marshal the movements array to JSON
	movementJSON, err := json.MarshalIndent(movements, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshaling stock movements: %v", err)
	}

	// Write to file
	stockFile := filepath.Join("testdata", "stock.json")
	if err := ioutil.WriteFile(stockFile, movementJSON, 0644); err != nil {
		return fmt.Errorf("error writing stock movements: %v", err)
	}

	return nil
}

// LoadStockMovements loads all stock movements from the stock.json file
func (c *VSCUClient) LoadStockMovements() ([]StockMovementRecord, error) {
	stockFile := filepath.Join("testdata", "stock.json")

	// Read existing file
	data, err := ioutil.ReadFile(stockFile)
	if err != nil {
		if os.IsNotExist(err) {
			// If file doesn't exist, return empty array
			return []StockMovementRecord{}, nil
		}
		return nil, fmt.Errorf("error reading stock movements: %v", err)
	}

	// Parse JSON array
	var movements []StockMovementRecord
	if err := json.Unmarshal(data, &movements); err != nil {
		return nil, fmt.Errorf("error parsing stock movements: %v", err)
	}

	return movements, nil
}
