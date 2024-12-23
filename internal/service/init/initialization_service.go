package init

import "etims-test/internal/app/models"

// InitializationService defines the interface for device initialization and data synchronization.
type InitializationService interface {
	InitializeAndSync() error
	InitializeDevice() (*models.InitResponse, error)
}
