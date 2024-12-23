package client

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"etims-test/models"
)

type StockMovement struct {
	Tin       string `json:"tin"`
	BhfId     string `json:"bhfId"`
	ItemCd    string `json:"itemCd"`
	RsdQty    int    `json:"rsdQty"`
	LastReqDt string `json:"lastReqDt"`
}

func (c *VSCUClient) HandleStockMovement(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the sales request
	var salesReq models.SalesRequest
	if err := json.NewDecoder(r.Body).Decode(&salesReq); err != nil {
		logrus.WithError(err).Error("Failed to decode sales request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create stock movement from sale
	salesData, _ := json.Marshal(salesReq)
	stockMovement, err := c.CreateStockMovementFromSale(salesData)
	if err != nil {
		logrus.WithError(err).Error("Failed to create stock movement")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response with movement details
	response := models.APIResponse{
		ResultCd:  "200",
		ResultMsg: "Stock movement created successfully",
		Data:      stockMovement,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *VSCUClient) HandleGetStockMovements(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Load all stock movements
	movements, err := c.LoadStockMovements()
	if err != nil {
		logrus.WithError(err).Error("Failed to load stock movements")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return movements
	response := models.APIResponse{
		ResultCd:  "200",
		ResultMsg: "Stock movements retrieved successfully",
		Data:      movements,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RegisterStockHandlers registers all stock-related handlers
func RegisterStockHandlers(mux *http.ServeMux, client *VSCUClient) {
	mux.HandleFunc("/stock/movement", client.HandleStockMovement)
	mux.HandleFunc("/stock/movements", client.HandleGetStockMovements)
}
