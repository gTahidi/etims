package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"etims-test/internal/app/models"
	"etims-test/internal/client"
)

type Config struct {
	API struct {
		BaseURL string `json:"baseURL"`
		Tin     string `json:"tin"`
		BhfId   string `json:"bhfId"`
		CmcKey  string `json:"cmcKey"`
	} `json:"api"`
	Logging struct {
		Level  string `json:"level"`
		Format string `json:"format"`
	} `json:"logging"`
}

// loadConfig reads a JSON configuration file and returns a Config object.
// The function returns an error if the file cannot be read or if the JSON
// cannot be parsed into a Config object.
func loadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return &config, nil
}

// setupLogger configures a logger using the settings in the given configuration object.
//
// If the logging level in the configuration is invalid, it defaults to the InfoLevel.
// The logging format is set to JSON if the configuration specifies so, otherwise it
// defaults to a text format with timestamps.
func setupLogger(config *Config) *logrus.Logger {
	logger := logrus.New()

	// Set logging level
	level, err := logrus.ParseLevel(config.Logging.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set format
	if config.Logging.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	return logger
}

// Main is the entry point of the application.
//
// It sets up a logger and VSCU client using the configuration file, then
// processes test data files in order. Each file is processed according to its
// type, and the response is logged.
//
// The supported file types are:
//
// - branch_customer.json: SaveBranchCustomer
// - item.json: SaveItem
// - stock.json: SaveStock
// - purchase.json: SavePurchase and CreateStockMovementFromPurchase
// - sales.json: SaveSales and CreateStockMovementFromSale
func main() {
	// Parse command line flags
	configPath := flag.String("config", "config/config.json", "path to configuration file")
	testDataDir := flag.String("testdata", "testdata", "directory containing test data files")
	flag.Parse()

	// Load configuration
	config, err := loadConfig(*configPath)
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Setup logger
	logger := setupLogger(config)

	// Create VSCU client
	client := client.NewVSCUClient(
		config.API.BaseURL,
		config.API.Tin,
		config.API.BhfId,
		config.API.CmcKey,
		logger,
	)

	// Define the order of processing
	fileOrder := []string{
		"branch_customer.json",
		"item.json",
		"stock.json",
		"purchase.json",
		"sales.json",
	}

	// Process files in order
	for _, filename := range fileOrder {
		filepath := filepath.Join(*testDataDir, filename)

		// Skip if file doesn't exist
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			logger.WithField("file", filepath).Info("File does not exist, skipping")
			continue
		}

		logger.WithField("file", filepath).Info("Processing test data file")

		// Read and process the test data file
		data, err := os.ReadFile(filepath)
		if err != nil {
			logger.WithError(err).Error("Failed to read test data file")
			continue
		}

		// Process based on file type
		switch filename {
		case "branch_customer.json":
			var custReq models.BranchCustomerRequest
			if err := json.Unmarshal(data, &custReq); err != nil {
				logger.WithError(err).Error("Failed to parse branch customer request data")
				continue
			}

			resp, err := client.SaveBranchCustomer(custReq)
			if err != nil {
				logger.WithError(err).Error("Failed to send branch customer request")
				continue
			}
			logger.WithField("response", resp).Info("Branch customer request successful")

		case "item.json":
			var itemsList struct {
				Tin   string               `json:"tin"`
				BhfId string               `json:"bhfId"`
				Items []models.ItemRequest `json:"items"`
			}
			if err := json.Unmarshal(data, &itemsList); err != nil {
				logger.WithError(err).Error("Failed to parse items list")
				continue
			}

			for _, item := range itemsList.Items {
				resp, err := client.SaveItem(item)
				if err != nil {
					logger.WithError(err).Error("Failed to send item request")
					continue
				}
				logger.WithFields(logrus.Fields{
					"itemCd":   item.ItemCd,
					"response": resp,
				}).Info("Item registration successful")
			}

		case "stock.json":
			var stockReqArray []models.StockRequest
			if err := json.Unmarshal(data, &stockReqArray); err != nil {
				var singleStockReq models.StockRequest
				if err := json.Unmarshal(data, &singleStockReq); err != nil {
					logger.WithError(err).Error("Failed to parse stock request data")
					continue
				}
				stockReqArray = []models.StockRequest{singleStockReq}
			}

			for _, req := range stockReqArray {
				for i := range req.ItemList {
					if req.ItemList[i].DcRt == 0 {
						req.ItemList[i].DcRt = 0
					}
					if req.ItemList[i].DcAmt == 0 {
						req.ItemList[i].DcAmt = 0
					}
					if req.ItemList[i].TotDcAmt == 0 {
						req.ItemList[i].TotDcAmt = 0
					}

					if req.SarTyCd == "01" {
						req.ItemList[i].StockTyCd = "01"
					} else if req.SarTyCd == "02" {
						req.ItemList[i].StockTyCd = "02"
						req.ItemList[i].Qty = -math.Abs(req.ItemList[i].Qty)
					}
				}

				resp, err := client.SaveStock(req)
				if err != nil {
					logger.WithError(err).Error("Failed to send stock request")
					continue
				}
				logger.WithField("response", resp).Info("Stock request successful")
			}

		case "purchase.json":
			var purchaseReq models.PurchaseRequest
			if err := json.Unmarshal(data, &purchaseReq); err != nil {
				logger.WithError(err).Error("Failed to parse purchase request data")
				continue
			}

			resp, err := client.SavePurchase(purchaseReq)
			if err != nil {
				logger.WithError(err).Error("Failed to send purchase request")
				continue
			}
			logger.WithField("response", resp).Info("Purchase request successful")

			stockMovement, err := client.CreateStockMovementFromPurchase(data)
			if err != nil {
				logger.WithError(err).Error("Failed to create stock movement for purchase")
				continue
			}
			logger.WithField("stockMovement", stockMovement).Info("Stock movement created for purchase")

		case "sales.json":
			var salesReq models.SalesRequest
			if err := json.Unmarshal(data, &salesReq); err != nil {
				logger.WithError(err).Error("Failed to parse sales request data")
				continue
			}

			stockMovement, err := client.CreateStockMovementFromSale(data)
			if err != nil {
				logger.WithError(err).Error("Failed to create stock movement - possibly insufficient stock")
				continue
			}

			resp, err := client.SaveSales(salesReq)
			if err != nil {
				logger.WithError(err).Error("Failed to send sales request")
				continue
			}
			logger.WithField("response", resp).Info("Sales request successful")
			logger.WithField("stockMovement", stockMovement).Info("Stock movement created for sale")
		}
	}
}
