package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"corporation-db/internal/infrastructure"
	"corporation-db/internal/usecase"
	"corporation-db/internal/utils"

	_ "github.com/lib/pq"
)

const (
	defaultDBHost     = "localhost"
	defaultDBPort     = "5432"
	defaultDBUser     = "postgres"
	defaultDBPassword = "password"
	defaultDBName     = "corporation_db"
)

func main() {
	// Parse command line flags
	var (
		dbHost     = flag.String("db-host", getEnv("DB_HOST", defaultDBHost), "Database host")
		dbPort     = flag.String("db-port", getEnv("DB_PORT", defaultDBPort), "Database port")
		dbUser     = flag.String("db-user", getEnv("DB_USER", defaultDBUser), "Database user")
		dbPassword = flag.String("db-password", getEnv("DB_PASSWORD", defaultDBPassword), "Database password")
		dbName     = flag.String("db-name", getEnv("DB_NAME", defaultDBName), "Database name")
		sslMode    = flag.String("ssl-mode", getEnv("DB_SSL_MODE", "disable"), "Database SSL mode")
		dryRun     = flag.Bool("dry-run", false, "Perform a dry run without actually importing data")
		testCSV    = flag.String("test-csv", "", "Use a local test CSV file instead of downloading from gBizINFO")
		help       = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		showHelp()
		os.Exit(0)
	}

	log.Println("Starting gBizINFO Corporation Data Batch Import")
	log.Printf("Target Database: %s:%s/%s", *dbHost, *dbPort, *dbName)

	if *dryRun {
		log.Println("DRY RUN MODE: No data will be imported")
	}

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

	// Connect to database
	db, err := connectToDatabase(*dbHost, *dbPort, *dbUser, *dbPassword, *dbName, *sslMode)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to database")

	if *dryRun {
		log.Println("Dry run completed successfully")
		return
	}

	// Initialize repositories and use cases
	corporationRepo := infrastructure.NewCorporationRepository(db)
	gbizClient := infrastructure.NewGBizClient()
	textConverter := utils.NewTextConverter()
	corporationUsecase := usecase.NewCorporationUsecase(corporationRepo, gbizClient, textConverter)

	// Start import process with progress tracking
	startTime := time.Now()

	if *testCSV != "" {
		log.Printf("Using test CSV file: %s", *testCSV)
		err = corporationUsecase.ImportFromTestCSV(ctx, *testCSV)
	} else {
		err = corporationUsecase.ImportFromGBizINFOWithProgress(ctx, func(stage string, progress float64) {
			log.Printf("Progress: %s (%.1f%%)", stage, progress)
		})
	}

	if err != nil {
		log.Fatalf("Import failed: %v", err)
	}

	duration := time.Since(startTime)
	log.Printf("Import completed successfully in %v", duration)

	// Show import statistics
	stats, err := corporationUsecase.GetImportStats(ctx)
	if err != nil {
		log.Printf("Warning: Failed to get import stats: %v", err)
	} else {
		log.Printf("Import Statistics:")
		log.Printf("  Total Corporations: %d", stats.TotalCorporations)
		log.Printf("  Status Breakdown:")
		for status, count := range stats.StatusCounts {
			log.Printf("    %s: %d", status, count)
		}
	}

	log.Println("Batch import completed successfully")
}

func connectToDatabase(host, port, user, password, dbname, sslmode string) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func showHelp() {
	fmt.Println(`gBizINFO Corporation Data Batch Import

This tool downloads and imports basic corporation information from gBizINFO API
into the local PostgreSQL database.

Usage:
  batch [options]

Options:
  -db-host string      Database host (default: localhost, env: DB_HOST)
  -db-port string      Database port (default: 5432, env: DB_PORT)
  -db-user string      Database user (default: postgres, env: DB_USER)
  -db-password string  Database password (default: password, env: DB_PASSWORD)
  -db-name string      Database name (default: corporation_db, env: DB_NAME)
  -ssl-mode string     Database SSL mode (default: disable, env: DB_SSL_MODE)
  -dry-run             Perform a dry run without importing data
  -test-csv string     Use a local test CSV file instead of downloading from gBizINFO
  -help                Show this help message

Environment Variables:
  DB_HOST              Database host
  DB_PORT              Database port
  DB_USER              Database user
  DB_PASSWORD          Database password
  DB_NAME              Database name
  DB_SSL_MODE          Database SSL mode

Examples:
  # Basic usage with default settings
  ./batch

  # Using environment variables
  DB_HOST=production.db.com DB_USER=admin ./batch

  # Dry run to test configuration
  ./batch -dry-run

  # Custom database connection
  ./batch -db-host mydb.com -db-port 5433 -db-user admin

  # Test with local CSV file
  ./batch -test-csv ./test_corporations.csv

The tool will:
1. Download the latest basic corporation information ZIP file from gBizINFO
   (or use a local CSV file if -test-csv is specified)
2. Extract and parse the CSV data
3. Perform bulk upsert operations to update the database
4. Show progress and statistics during the import process

Note: The import process may take several minutes to hours depending on
the data size and network/database performance.`)
}
