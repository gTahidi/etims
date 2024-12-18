package models

// StockMovementArray represents an array of stock movements
type StockMovementArray []StockRequest

// Validate validates all stock movements in the array
func (s StockMovementArray) Validate() error {
	for _, movement := range s {
		if err := movement.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// GetLatestSarNo returns the highest SarNo in the array
func (s StockMovementArray) GetLatestSarNo() int {
	maxSarNo := 0
	for _, movement := range s {
		if movement.SarNo > maxSarNo {
			maxSarNo = movement.SarNo
		}
	}
	return maxSarNo
}

// AddMovement adds a new movement to the array
func (s *StockMovementArray) AddMovement(movement StockRequest) {
	*s = append(*s, movement)
}
