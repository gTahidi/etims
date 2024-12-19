package testmanager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// DataSnapshot represents a point-in-time snapshot of test data
type DataSnapshot struct {
	TestCase    string
	Data        map[string]interface{}
	Timestamp   time.Time
	ValidatedBy []string
}

// ValidationFunc defines a function type for data validation
type ValidationFunc func(data interface{}) error

// DataManager handles test data lifecycle and validation
type DataManager struct {
	mu          sync.RWMutex
	currentData map[string]interface{}
	history     []DataSnapshot
	validations map[string][]ValidationFunc
	logger      *logrus.Logger
	dataDir     string
}

// NewDataManager creates a new instance of DataManager
func NewDataManager(logger *logrus.Logger, dataDir string) *DataManager {
	return &DataManager{
		currentData: make(map[string]interface{}),
		history:     make([]DataSnapshot, 0),
		validations: make(map[string][]ValidationFunc),
		logger:      logger,
		dataDir:     dataDir,
	}
}

// RegisterValidation adds a validation function for a specific test case
func (dm *DataManager) RegisterValidation(testCase string, validation ValidationFunc) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.validations[testCase] == nil {
		dm.validations[testCase] = make([]ValidationFunc, 0)
	}
	dm.validations[testCase] = append(dm.validations[testCase], validation)
}

// PrepareTestData generates or retrieves test data for a specific test case
func (dm *DataManager) PrepareTestData(testCase string, baseData interface{}) (interface{}, error) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// Create a new data instance based on the test case
	data, err := dm.generateTestData(testCase, baseData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate test data: %w", err)
	}

	// Validate the data
	if err := dm.validateData(testCase, data); err != nil {
		return nil, fmt.Errorf("data validation failed: %w", err)
	}

	// Store the current state
	dm.currentData[testCase] = data
	dm.saveSnapshot(testCase, data)

	return data, nil
}

// generateTestData creates test data based on the test case and base data
func (dm *DataManager) generateTestData(testCase string, baseData interface{}) (interface{}, error) {
	// Clone the base data to avoid modifications
	var clonedData interface{}
	dataBytes, err := json.Marshal(baseData)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(dataBytes, &clonedData); err != nil {
		return nil, err
	}

	// Apply test case specific modifications
	switch testCase {
	case "SAVE_PURCHASE":
		return dm.modifyPurchaseData(clonedData)
	case "SAVE_STOCK":
		return dm.modifyStockData(clonedData)
	default:
		return clonedData, nil
	}
}

// validateData runs all registered validations for a test case
func (dm *DataManager) validateData(testCase string, data interface{}) error {
	validations, exists := dm.validations[testCase]
	if !exists {
		return nil
	}

	for _, validate := range validations {
		if err := validate(data); err != nil {
			return err
		}
	}
	return nil
}

// saveSnapshot creates and stores a snapshot of the current data state
func (dm *DataManager) saveSnapshot(testCase string, data interface{}) {
	snapshot := DataSnapshot{
		TestCase:  testCase,
		Data:      map[string]interface{}{"data": data},
		Timestamp: time.Now(),
	}
	dm.history = append(dm.history, snapshot)
}

// GetTestData retrieves the current data for a test case
func (dm *DataManager) GetTestData(testCase string) (interface{}, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	data, exists := dm.currentData[testCase]
	if !exists {
		return nil, fmt.Errorf("no data found for test case: %s", testCase)
	}
	return data, nil
}

// GetHistory returns the history of data changes
func (dm *DataManager) GetHistory() []DataSnapshot {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	
	return append([]DataSnapshot{}, dm.history...)
}

// GetData retrieves data from a JSON file
func (dm *DataManager) GetData(key string) ([]byte, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	filePath := filepath.Join(dm.dataDir, key+".json")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return data, nil
}

// SaveData saves data to a JSON file
func (dm *DataManager) SaveData(key string, data []byte) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	filePath := filepath.Join(dm.dataDir, key+".json")
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	return nil
}

// modifyPurchaseData applies purchase-specific modifications
func (dm *DataManager) modifyPurchaseData(data interface{}) (interface{}, error) {
	// TODO: Implement purchase data modifications
	return data, nil
}

// modifyStockData applies stock-specific modifications
func (dm *DataManager) modifyStockData(data interface{}) (interface{}, error) {
	// TODO: Implement stock data modifications
	return data, nil
}
