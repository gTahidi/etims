package scripts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"etims-test/models"
)

// StockData represents the structure for stock data
type StockData struct {
	Tin         string            `json:"tin"`
	BhfId       string            `json:"bhfId"`
	SarNo       int               `json:"sarNo"`
	OrgSarNo    int               `json:"orgSarNo"`
	RegTyCd     string            `json:"regTyCd"`
	CustTin     string            `json:"custTin"`
	CustNm      string            `json:"custNm"`
	CustBhfId   string            `json:"custBhfId"`
	SarTyCd     string            `json:"sarTyCd"`
	OcrnDt      string            `json:"ocrnDt"`
	StockRlsDt  string            `json:"stockRlsDt"`
	TotItemCnt  int               `json:"totItemCnt"`
	TotTaxblAmt float64           `json:"totTaxblAmt"`
	TotTaxAmt   float64           `json:"totTaxAmt"`
	TotAmt      float64           `json:"totAmt"`
	Remark      string            `json:"remark"`
	RegrId      string            `json:"regrId"`
	RegrNm      string            `json:"regrNm"`
	ModrId      string            `json:"modrId"`
	ModrNm      string            `json:"modrNm"`
	StockItems  []models.StockItem `json:"itemList"`
}

// ProcessStockData reads and processes the stock.json file
func ProcessStockData() error {
	// Read stock.json file
	stockFile := filepath.Join("testdata", "stock.json")
	data, err := ioutil.ReadFile(stockFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Stock file does not exist: %s", stockFile)
			return nil
		}
		return fmt.Errorf("failed to read stock file: %w", err)
	}

	// Try to unmarshal as array first
	var stockArray []StockData
	if err := json.Unmarshal(data, &stockArray); err == nil {
		log.Printf("Successfully parsed %d stock records", len(stockArray))
		// Process the array of stock data
		return processStockArray(stockArray)
	}

	// If array unmarshal fails, try single object
	var singleStock StockData
	if err := json.Unmarshal(data, &singleStock); err != nil {
		return fmt.Errorf("failed to parse stock data: %w", err)
	}

	// Process single stock record
	return processStockRecord(&singleStock)
}

// processStockArray processes an array of stock records
func processStockArray(stocks []StockData) error {
	for i, stock := range stocks {
		if err := processStockRecord(&stock); err != nil {
			return fmt.Errorf("error processing stock record %d: %w", i+1, err)
		}
	}
	return nil
}

// processStockRecord processes a single stock record
func processStockRecord(stock *StockData) error {
	// Validate the stock record
	if err := validateStockRecord(stock); err != nil {
		return fmt.Errorf("stock validation failed: %w", err)
	}

	// Here you can add any additional processing logic
	// For example, saving to database, updating other files, etc.
	log.Printf("Processed stock record: SarNo=%d, Tin=%s", stock.SarNo, stock.Tin)
	return nil
}

// validateStockRecord validates a single stock record
func validateStockRecord(stock *StockData) error {
	if stock.Tin == "" {
		return fmt.Errorf("tin is required")
	}
	if stock.BhfId == "" {
		return fmt.Errorf("bhfId is required")
	}
	if stock.SarNo <= 0 {
		return fmt.Errorf("invalid sarNo")
	}
	if len(stock.StockItems) == 0 {
		return fmt.Errorf("stock items list cannot be empty")
	}
	return nil
}
