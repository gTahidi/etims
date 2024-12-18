package main

import (
	"log"

	"etims-test/scripts"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	
	if err := scripts.ProcessStockData(); err != nil {
		log.Fatalf("Error processing stock data: %v", err)
	}
	
	log.Println("Stock data processing completed successfully")
}
