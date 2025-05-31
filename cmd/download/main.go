package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"corporation-db/internal/infrastructure"
)

func main() {
	// Parse command line flags
	var (
		outputPath = flag.String("output", "", "Output path for downloaded ZIP file (required)")
		help       = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		showHelp()
		os.Exit(0)
	}

	if *outputPath == "" {
		fmt.Println("Error: -output flag is required")
		showHelp()
		os.Exit(1)
	}

	log.Println("Starting gBizINFO Data Download")
	log.Printf("Output path: %s", *outputPath)

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		log.Printf("Received signal %v, initiating graceful shutdown...", sig)
		cancel()
	}()

	// Initialize gBizINFO client
	gbizClient := infrastructure.NewGBizClient()

	// Download CSV ZIP file
	log.Println("Downloading basic corporation information from gBizINFO...")
	zipPath, err := gbizClient.DownloadBasicInfoCSV(ctx)
	if err != nil {
		log.Fatalf("Failed to download CSV from gBizINFO: %v", err)
	}
	defer gbizClient.Cleanup(zipPath)

	log.Printf("Downloaded ZIP file: %s", zipPath)

	// Move the downloaded file to the specified output path
	err = os.Rename(zipPath, *outputPath)
	if err != nil {
		log.Fatalf("Failed to move downloaded file to output path: %v", err)
	}

	log.Printf("Successfully downloaded and saved to: %s", *outputPath)

	// Display CSV headers for verification
	log.Println("Displaying CSV headers...")
	err = gbizClient.DisplayCSVHeaders(*outputPath)
	if err != nil {
		log.Printf("Warning: Failed to display CSV headers: %v", err)
	}

	log.Println("Download completed successfully")
}

func showHelp() {
	fmt.Println(`gBizINFO Data Download Tool

This tool downloads the latest basic corporation information ZIP file from gBizINFO API.

Usage:
  download -output <path>

Options:
  -output string    Output path for downloaded ZIP file (required)
  -help            Show this help message

Examples:
  # Download to current directory
  ./download -output ./gbiz_data.zip

  # Download to specific directory
  ./download -output /data/gbiz_$(date +%Y%m%d).zip

The tool will:
1. Download the latest basic corporation information ZIP file from gBizINFO
2. Save it to the specified output path
3. Display CSV headers for verification

Note: The download may take several minutes depending on network speed and data size.`)
}
