package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
			
			// Send sales request
			resp, err := client.SaveSales(salesReq)
			if err != nil {
				logger.WithError(err).Error("Failed to send sales request")
				return nil
			}

			logger.WithField("response", resp).Info("Sales request successful")

		case "stock.json":
			var stockReq models.StockRequest
			if err := json.Unmarshal(data, &stockReq); err != nil {
				logger.WithError(err).Error("Failed to parse stock request data")
				return nil
			}

			// Send stock request
			resp, err := client.SaveStock(stockReq)
			if err != nil {
				logger.WithError(err).Error("Failed to send stock request")
				return nil
			}

			logger.WithField("response", resp).Info("Stock request successful")

		case "item.json":
			var itemReq models.ItemRequest
			if err := json.Unmarshal(data, &itemReq); err != nil {
				logger.WithError(err).Error("Failed to parse item request data")
				return nil
			}

			// Send item request
			resp, err := client.SaveItem(itemReq)
			if err != nil {
				logger.WithError(err).Error("Failed to send item request")
				return nil
			}

			logger.WithField("response", resp).Info("Item request successful")

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
func initializeAndSync(client *client.VSCUClient, logger *logrus.Logger) error {
	// 1. Device initialization
	logger.Info("Initializing device")
	resp, err := client.InitializeDevice()
	if err != nil {
		return fmt.Errorf("device initialization failed: %w", err)
	}
	logger.WithField("response", resp).Info("Device initialization successful")

	// Get current timestamp for sync
	lastReqDt := client.GetCurrentTimestamp()

	// 2. Sync basic data
	logger.Info("Syncing basic data")

	// Get code list
	if _, err := client.GetCodeList(lastReqDt); err != nil {
		return fmt.Errorf("failed to sync code list: %w", err)
	}

	// Get item classification list
	if _, err := client.GetItemClassList(lastReqDt); err != nil {
		return fmt.Errorf("failed to sync item classifications: %w", err)
	}

	// Get branch list
	if _, err := client.GetBranchList(lastReqDt); err != nil {
		return fmt.Errorf("failed to sync branch list: %w", err)
	}

	// Get notices
	if _, err := client.GetNoticeList(lastReqDt); err != nil {
		return fmt.Errorf("failed to sync notices: %w", err)
	}

	// 3. Sync imported items (if head office)
	if client.BhfId == "00" {
		logger.Info("Syncing imported items")
		if _, err := client.GetImportedItems(lastReqDt); err != nil {
			return fmt.Errorf("failed to sync imported items: %w", err)
		}
	}

	logger.Info("Initialization and sync completed successfully")
	return nil
}
