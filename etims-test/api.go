package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"etims-test/client"
	"etims-test/models"
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

	// Initialize device and sync data
	if err := initializeAndSync(client, logger); err != nil {
		logger.WithError(err).Fatal("Failed to initialize device and sync data")
	}

	// Process test data files
	err = filepath.Walk(*testDataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Process only JSON files
		if filepath.Ext(path) != ".json" {
			return nil
		}

		logger.WithField("file", path).Info("Processing test data file")

		// Read and process the test data file
		data, err := os.ReadFile(path)
		if err != nil {
			logger.WithError(err).Error("Failed to read test data file")
			return nil
		}

		// Determine the request type based on the filename
		switch filepath.Base(path) {
		case "sales.json":
			var salesReq models.SalesRequest
			if err := json.Unmarshal(data, &salesReq); err != nil {
				logger.WithError(err).Error("Failed to parse sales request data")
				return nil
			}

			// Create stock movement first to validate stock levels
			stockMovement, err := client.CreateStockMovementFromSale(data)
			if err != nil {
				logger.WithError(err).Error("Failed to create stock movement - possibly insufficient stock")
				return nil
			}

			// If stock validation passed, proceed with the sale
			resp, err := client.SaveSales(salesReq)
			if err != nil {
				logger.WithError(err).Error("Failed to send sales request")
				return nil
			}
			logger.WithField("response", resp).Info("Sales request successful")
			logger.WithField("stockMovement", stockMovement).Info("Stock movement created for sale")

		case "purchase.json":
			var purchaseReq models.PurchaseRequest
			if err := json.Unmarshal(data, &purchaseReq); err != nil {
				logger.WithError(err).Error("Failed to parse purchase request data")
				return nil
			}

			// First save the purchase
			resp, err := client.SavePurchase(purchaseReq)
			if err != nil {
				logger.WithError(err).Error("Failed to send purchase request")
				return nil
			}
			logger.WithField("response", resp).Info("Purchase request successful")

			// Then create stock movement for incoming stock
			stockMovement, err := client.CreateStockMovementFromPurchase(data)
			if err != nil {
				logger.WithError(err).Error("Failed to create stock movement for purchase")
				return nil
			}
			logger.WithField("stockMovement", stockMovement).Info("Stock movement created for purchase")

		case "stock.json":
			// First try to parse as array
			var stockReqArray []models.StockRequest
			if err := json.Unmarshal(data, &stockReqArray); err != nil {
				// If array parsing fails, try single object
				var singleStockReq models.StockRequest
				if err := json.Unmarshal(data, &singleStockReq); err != nil {
					logger.WithError(err).Error("Failed to parse stock request data as array or single object")
					return nil
				}
				stockReqArray = []models.StockRequest{singleStockReq}
			}

			// Process each stock request
			for _, req := range stockReqArray {
				// Ensure required fields are set for each item
				for i := range req.ItemList {
					// Initialize discount fields if not set
					if req.ItemList[i].DcRt == 0 {
						req.ItemList[i].DcRt = 0
					}
					if req.ItemList[i].DcAmt == 0 {
						req.ItemList[i].DcAmt = 0
					}
					if req.ItemList[i].TotDcAmt == 0 {
						req.ItemList[i].TotDcAmt = 0
					}

					// Set stock type based on SAR type
					if req.SarTyCd == "01" { // Stock In
						req.ItemList[i].StockTyCd = "01" // Stock In
					} else if req.SarTyCd == "02" { // Stock Out
						req.ItemList[i].StockTyCd = "02" // Stock Out
						// For stock out, make quantity negative
						req.ItemList[i].Qty = -math.Abs(req.ItemList[i].Qty)
					}
				}

				// Send stock request
				resp, err := client.SaveStock(req)
				if err != nil {
					logger.WithError(err).Error("Failed to send stock request")
					continue
				}
				logger.WithField("response", resp).Info("Stock request successful")
			}

		case "item.json":
			// Parse items list
			var itemsList struct {
				Tin   string               `json:"tin"`
				BhfId string               `json:"bhfId"`
				Items []models.ItemRequest `json:"items"`
			}
			if err := json.Unmarshal(data, &itemsList); err != nil {
				logger.WithError(err).Error("Failed to parse items list")
				return nil
			}

			// Register each item individually
			for _, item := range itemsList.Items {
				// Send item request
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

		case "branch_customer.json":
			var custReq models.BranchCustomerRequest
			if err := json.Unmarshal(data, &custReq); err != nil {
				logger.WithError(err).Error("Failed to parse branch customer request data")
				return nil
			}

			// Send branch customer request
			resp, err := client.SaveBranchCustomer(custReq)
			if err != nil {
				logger.WithError(err).Error("Failed to send branch customer request")
				return nil
			}

			logger.WithField("response", resp).Info("Branch customer request successful")
		}

		return nil
	})

	if err != nil {
		logger.WithError(err).Fatal("Failed to process test data files")
	}
}

// Initialize device and sync basic data
// Initialize device and sync basic data
func initializeAndSync(client *client.VSCUClient, logger *logrus.Logger) error {
	// 1. Device initialization
	logger.Info("Initializing device")
	_, err := client.InitializeDevice()
	if err != nil {
		// Check if error is APIError type and has code 902 (device already initialized)
		if e, ok := err.(interface{ Code() string }); ok && e.Code() == "902" {
			logger.Info("Device is already initialized, proceeding with sync")
		} else {
			return fmt.Errorf("device initialization failed: %w", err)
		}
	}

	// Get current timestamp for sync
	lastReqDt := client.GetCurrentTimestamp()

	// 2. Sync basic data
	logger.Info("Syncing basic data")

	// Get code list
	codeResp, err := client.GetCodeList(lastReqDt)
	if err != nil {
		return fmt.Errorf("failed to sync code list: %w", err)
	}
	if codeResp.ResultCd == "001" {
		logger.Info("No code list found")
	} else if codeResp.ResultCd != "000" {
		return fmt.Errorf("failed to sync code list: API error: %s (code: %s)", codeResp.ResultMsg, codeResp.ResultCd)
	}

	// Get item classification list
	itemClassResp, err := client.GetItemClassList(lastReqDt)
	if err != nil {
		return fmt.Errorf("failed to sync item classifications: %w", err)
	}
	if itemClassResp.ResultCd == "001" {
		logger.Info("No item classification list found")
	} else if itemClassResp.ResultCd != "000" {
		return fmt.Errorf("failed to sync item classifications: API error: %s (code: %s)", itemClassResp.ResultMsg, itemClassResp.ResultCd)
	}

	// Get branch list
	branchResp, err := client.GetBranchList(lastReqDt)
	if err != nil {
		return fmt.Errorf("failed to sync branch list: %w", err)
	}
	if branchResp.ResultCd == "001" {
		logger.Info("No branch list found")
	} else if branchResp.ResultCd != "000" {
		return fmt.Errorf("failed to sync branch list: API error: %s (code: %s)", branchResp.ResultMsg, branchResp.ResultCd)
	}

	// Get notices
	noticeResp, err := client.GetNoticeList(lastReqDt)
	if err != nil {
		return fmt.Errorf("failed to sync notices: %w", err)
	}
	if noticeResp.ResultCd == "001" {
		logger.Info("No notices found")
	} else if noticeResp.ResultCd != "000" {
		return fmt.Errorf("failed to sync notices: API error: %s (code: %s)", noticeResp.ResultMsg, noticeResp.ResultCd)
	}

	// 3. Sync imported items (if head office)
	if client.BhfId == "00" {
		logger.Info("Syncing imported items")
		importedItemsResp, err := client.GetImportedItems(lastReqDt)
		if err != nil {
			return fmt.Errorf("failed to sync imported items: %w", err)
		}
		if importedItemsResp.ResultCd == "001" {
			logger.Info("No imported items found")
		} else if importedItemsResp.ResultCd != "000" {
			return fmt.Errorf("failed to sync imported items: API error: %s (code: %s)", importedItemsResp.ResultMsg, importedItemsResp.ResultCd)
		}
	}

	logger.Info("Initialization and sync completed successfully")
	return nil
}
