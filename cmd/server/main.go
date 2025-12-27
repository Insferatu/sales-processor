package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sales-processor/internal/handlers"
	"github.com/sales-processor/internal/sheets"
	"github.com/sales-processor/internal/telegram"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()

	// Initialize Telegram client (shared between handlers)
	telegramClient, err := telegram.NewClient()
	if err != nil {
		log.Fatalf("Failed to initialize Telegram client: %v", err)
	}

	// Initialize 3D Toy Sale sheets client
	toySheetsClient, err := sheets.NewClient(ctx, sheets.Config{
		SpreadsheetID: os.Getenv("TOY_SPREADSHEET_ID"),
		SheetName:     os.Getenv("TOY_SHEET_NAME"),
		ColumnRange:   "A:E", // 5 columns: Время продажи | Фигурка | Пластик | Цена | Тип оплаты
	})
	if err != nil {
		log.Fatalf("Failed to initialize Toy Google Sheets client: %v", err)
	}

	// Initialize Jewelry Sale sheets client
	jewelrySheetsClient, err := sheets.NewClient(ctx, sheets.Config{
		SpreadsheetID: os.Getenv("JEWELRY_SPREADSHEET_ID"),
		SheetName:     os.Getenv("JEWELRY_SHEET_NAME"),
		ColumnRange:   "A:D", // 4 columns: Время продажи | Товар | Цена | Тип оплаты
	})
	if err != nil {
		log.Fatalf("Failed to initialize Jewelry Google Sheets client: %v", err)
	}

	// Create handlers
	toySaleHandler := handlers.NewToySaleHandler(toySheetsClient, telegramClient)
	jewelrySaleHandler := handlers.NewJewelrySaleHandler(jewelrySheetsClient, telegramClient)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.HealthCheck)
	mux.HandleFunc("POST /3d-toy-sale/", toySaleHandler.HandleSale)
	mux.HandleFunc("POST /jewelry-sale/", jewelrySaleHandler.HandleSale)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Server starting on port %s", port)
	log.Printf("Endpoints:")
	log.Printf("  POST /3d-toy-sale/")
	log.Printf("  POST /jewelry-sale/")
	log.Printf("  GET  /health")

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
	log.Println("Server stopped")
}
