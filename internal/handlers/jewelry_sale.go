package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/sales-processor/internal/sheets"
	"github.com/sales-processor/internal/telegram"
)

// JewelrySaleEvent represents an incoming jewelry sale event
type JewelrySaleEvent struct {
	Time        string `json:"time"`
	Item        string `json:"item"`
	Price       string `json:"price"`
	PaymentType string `json:"paymentType"`
}

// JewelrySaleHandler handles incoming jewelry sale events
type JewelrySaleHandler struct {
	sheets   *sheets.Client
	telegram *telegram.Client
}

// NewJewelrySaleHandler creates a new JewelrySaleHandler with the given dependencies
func NewJewelrySaleHandler(sheetsClient *sheets.Client, telegramClient *telegram.Client) *JewelrySaleHandler {
	return &JewelrySaleHandler{
		sheets:   sheetsClient,
		telegram: telegramClient,
	}
}

// HandleSale processes incoming jewelry sale events
func (h *JewelrySaleHandler) HandleSale(w http.ResponseWriter, r *http.Request) {
	var event JewelrySaleEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate required fields
	if event.Item == "" {
		http.Error(w, "Missing required field: item", http.StatusBadRequest)
		return
	}

	// Set default timestamp if not provided
	if event.Time == "" {
		event.Time = time.Now().UTC().Format("2006-01-02 15:04:05")
	}

	// Process in parallel: send to Google Sheets and Telegram concurrently
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(2)

	// Send to Google Sheets
	go func() {
		defer wg.Done()
		if err := h.sheets.AppendRow(r.Context(), event.toSheetRow()); err != nil {
			log.Printf("Error appending to Google Sheets: %v", err)
			errChan <- err
		}
	}()

	// Send to Telegram
	go func() {
		defer wg.Done()
		message := event.toTelegramMessage()
		if err := h.telegram.SendMessage(message); err != nil {
			log.Printf("Error sending Telegram message: %v", err)
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	// Check for errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		log.Printf("Processed jewelry sale with %d errors", len(errors))
		http.Error(w, "Partial failure: some operations failed", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully processed jewelry sale: %s", event.Item)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Sale processed successfully",
	})
}

// toSheetRow converts the sale event to a row format for Google Sheets
// Columns: Время продажи | Товар | Цена | Тип оплаты
func (e *JewelrySaleEvent) toSheetRow() []interface{} {
	return []interface{}{
		e.Time,
		e.Item,
		e.Price,
		e.PaymentType,
	}
}

// toTelegramMessage formats the sale event as a Telegram message
func (e *JewelrySaleEvent) toTelegramMessage() string {
	return fmt.Sprintf(
		"Товар: %s\nПродано за: %s\nТип оплаты: %s",
		e.Item,
		e.Price,
		e.PaymentType,
	)
}
