package sales

import (
	"etims-test/internal/app/models/client"
)

// SalesService defines the interface for sales transaction operations.
type SalesService interface {
	SendSalesTransaction() error // Sends sales transaction data to the eTIMS API.
	//SendSalesInvoice() error       // Sends sales invoice data to the eTIMS API.
}

// salesService implements the SalesService interface.
type salesService struct {
	client *client.Client
}

// NewSalesService creates a new instance of salesService.
func NewSalesService(client *client.Client) SalesService {
	return &salesService{
		client: client,
	}
}

// SendSalesTransaction sends sales transaction data to the eTIMS API.
func (s *salesService) SendSalesTransaction() error {
	// Implement the logic to send sales transaction data.
	// Refer to etims-test/shemas/VSCU-Sales_Management/Send_sales_transaction_information_request_schema.json for the data structure.
	// Use s.client to make the API request.
	return nil
}

// // SendSalesInvoice sends sales invoice data to the eTIMS API.
// func (s *salesService) SendSalesInvoice() error {
// 	// Implement the logic to send sales invoice data.
// 	// Refer to etims-test/shemas/VSCU-Sales_Management/Send_sales_transaction_information_request_schema.json for the data structure.
// 	// Use s.client to make the API request.
// 	return nil
// }
