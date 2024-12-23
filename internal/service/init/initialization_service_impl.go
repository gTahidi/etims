package init

import (
	"etims-test/internal/app/models"
	"etims-test/internal/client"
	"fmt"
)

// initializationServiceImpl implements the InitializationService interface.
type initializationServiceImpl struct {
	client client.VSCUClient
}

// NewInitializationService creates a new instance of initializationServiceImpl.
func NewInitializationService(client client.VSCUClient) InitializationService {
	return &initializationServiceImpl{
		client: client,
	}
}

// InitializeAndSync initializes the device and synchronizes basic data.
func (s *initializationServiceImpl) InitializeAndSync() error {
	// 1. Device initialization
	if _, err := s.InitializeDevice(); err != nil {
		return fmt.Errorf("device initialization failed: %w", err)
	}

	// Get current timestamp for sync
	lastReqDt := s.client.GetCurrentTimestamp()

	// 2. Sync basic data

	// Get code list
	if _, err := s.client.GetCodeList(lastReqDt); err != nil {
		return fmt.Errorf("failed to sync code list: %w", err)
	}

	// Get item classification list
	if _, err := s.client.GetItemClassList(lastReqDt); err != nil {
		return fmt.Errorf("failed to sync item classifications: %w", err)
	}

	// Get branch list
	if _, err := s.client.GetBranchList(lastReqDt); err != nil {
		return fmt.Errorf("failed to sync branch list: %w", err)
	}

	// Get notices
	if _, err := s.client.GetNoticeList(lastReqDt); err != nil {
		return fmt.Errorf("failed to sync notices: %w", err)
	}

	// 3. Sync imported items (if head office)
	if s.client.IsHeadOffice() {
		if _, err := s.client.GetImportedItems(lastReqDt); err != nil {
			return fmt.Errorf("failed to sync imported items: %w", err)
		}
	}

	return nil
}

func (s *initializationServiceImpl) InitializeDevice() (*models.InitResponse, error) {
	// 1. Device initialization
	initResp, err := s.client.InitializeDevice()
	if err != nil {
		// Check if error is APIError type and has code 902 (device already initialized)
		if e, ok := err.(interface{ Code() string }); ok && e.Code() == "902" {
			return initResp, nil
		} else {
			return nil, fmt.Errorf("device initialization failed: %w", err)
		}
	}

	return initResp, nil
}
