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
			// Parse stock request
			var stockReq struct {
				Tin         string             `json:"tin"`
				BhfId       string             `json:"bhfId"`
				SarNo       int                `json:"sarNo"`
				OrgSarNo    int                `json:"orgSarNo"`
				RegTyCd     string             `json:"regTyCd"`
				CustTin     string             `json:"custTin"`
				CustNm      string             `json:"custNm"`
				CustBhfId   string             `json:"custBhfId"`
				SarTyCd     string             `json:"sarTyCd"`
				OcrnDt      string             `json:"ocrnDt"`
				StockRlsDt  string             `json:"stockRlsDt"`
				TotItemCnt  int                `json:"totItemCnt"`
				TotTaxblAmt float64            `json:"totTaxblAmt"`
				TotTaxAmt   float64            `json:"totTaxAmt"`
				TotAmt      float64            `json:"totAmt"`
				Remark      string             `json:"remark"`
				RegrId      string             `json:"regrId"`
				RegrNm      string             `json:"regrNm"`
				ModrId      string             `json:"modrId"`
				ModrNm      string             `json:"modrNm"`
				StockItems  []models.StockItem `json:"stockItems"`
			}
			if err := json.Unmarshal(data, &stockReq); err != nil {
				logger.WithError(err).Error("Failed to parse stock request data")
				return nil
			}

			// Create StockRequest with all fields
			req := models.StockRequest{
				BaseRequest: models.BaseRequest{
					Tin:   stockReq.Tin,
					BhfId: stockReq.BhfId,
				},
				SarNo:       stockReq.SarNo,
				OrgSarNo:    stockReq.OrgSarNo,
				RegTyCd:     stockReq.RegTyCd,
				CustTin:     stockReq.CustTin,
				CustNm:      stockReq.CustNm,
				CustBhfId:   stockReq.CustBhfId,
				SarTyCd:     stockReq.SarTyCd,
				OcrnDt:      stockReq.OcrnDt,
				StockRlsDt:  stockReq.StockRlsDt,
				TotItemCnt:  stockReq.TotItemCnt,
				TotTaxblAmt: stockReq.TotTaxblAmt,
				TotTaxAmt:   stockReq.TotTaxAmt,
				TotAmt:      stockReq.TotAmt,
				Remark:      stockReq.Remark,
				RegrId:      stockReq.RegrId,
				RegrNm:      stockReq.RegrNm,
				ModrId:      stockReq.ModrId,
				ModrNm:      stockReq.ModrNm,
				ItemList:    stockReq.StockItems,
			}

			// Send stock request
			resp, err := client.SaveStock(req)
			if err != nil {
				logger.WithError(err).Error("Failed to send stock request")
				return nil
			}

			logger.WithField("response", resp).Info("Stock request successful")

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
