package main

import (
	"corporation-db/internal/infrastructure"
	"log"
)

func main() {
	log.Println("Examining CSV headers from gBizINFO data...")

	// Path to the downloaded ZIP file
	zipPath := "/Users/zoo/projects/corporatioin-db/data/gbiz_20250531_114825.zip"

	// Initialize gBizINFO client
	gbizClient := infrastructure.NewGBizClient()

	// Display CSV headers
	if err := gbizClient.DisplayCSVHeaders(zipPath); err != nil {
		log.Fatalf("Failed to display CSV headers: %v", err)
	}
}
