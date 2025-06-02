package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"corporation-db/internal/infrastructure"
)

func main() {
	// Parse command line flags
	var (
		outputPath = flag.String("output", "", "Output path for downloaded finance ZIP file (required)")
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

	log.Println("Starting gBizINFO Finance Data Download")
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
	client := infrastructure.NewGBizClient()

	// Download finance data
	log.Println("Downloading finance data from gBizINFO...")
	financeReader, err := client.DownloadFinanceData(ctx)
	if err != nil {
		log.Printf("Error downloading finance data: %v", err)
		os.Exit(1)
	}
	defer func() {
		if closer, ok := financeReader.(io.Closer); ok {
			closer.Close()
		}
	}()

	// Create output file
	outputFile, err := os.Create(*outputPath)
	if err != nil {
		log.Printf("Error creating output file: %v", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	// Copy data to output file
	log.Println("Writing data to output file...")
	bytesWritten, err := io.Copy(outputFile, financeReader)
	if err != nil {
		log.Printf("Error writing to output file: %v", err)
		os.Exit(1)
	}

	log.Printf("Finance data download completed successfully")
	log.Printf("Downloaded %d bytes to %s", bytesWritten, *outputPath)
}

func showHelp() {
	fmt.Println("gBizINFO Finance Data Downloader")
	fmt.Println()
	fmt.Println("Downloads finance data ZIP file from gBizINFO API")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Printf("  %s [options]\n", os.Args[0])
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -output string    Output path for downloaded ZIP file (required)")
	fmt.Println("  -help            Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Printf("  %s -output data/finance_latest.zip\n", os.Args[0])
	fmt.Printf("  %s -output /tmp/finance_data.zip\n", os.Args[0])
}
