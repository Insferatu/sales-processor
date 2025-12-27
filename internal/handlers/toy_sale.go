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

// ToySaleEvent represents an incoming 3D toy sale event
type ToySaleEvent struct {
	Time        string `json:"time"`
	Item        string `json:"item"`
	Material    string `json:"material"`
	Price       string `json:"price"`
	PaymentType string `json:"paymentType"`
}

// ToySaleHandler handles incoming 3D toy sale events
type ToySaleHandler struct {
	sheets   *sheets.Client
	telegram *telegram.Client
}

// NewToySaleHandler creates a new ToySaleHandler with the given dependencies
func NewToySaleHandler(sheetsClient *sheets.Client, telegramClient *telegram.Client) *ToySaleHandler {
	return &ToySaleHandler{
		sheets:   sheetsClient,
		telegram: telegramClient,
	}
}

// HandleSale processes incoming 3D toy sale events
func (h *ToySaleHandler) HandleSale(w http.ResponseWriter, r *http.Request) {
	var event ToySaleEvent
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
		log.Printf("Processed toy sale with %d errors", len(errors))
		http.Error(w, "Partial failure: some operations failed", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully processed toy sale: %s - %s", event.Item, event.Material)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Sale processed successfully",
	})
}

// toSheetRow converts the sale event to a row format for Google Sheets
// Columns: Время продажи | Фигурка | Пластик | Цена | Тип оплаты
func (e *ToySaleEvent) toSheetRow() []interface{} {
	return []interface{}{
		e.Time,
		e.Item,
		e.Material,
		e.Price,
		e.PaymentType,
	}
}

// toTelegramMessage formats the sale event as a Telegram message
func (e *ToySaleEvent) toTelegramMessage() string {
	return fmt.Sprintf(
		"Фигурка: %s\nПластик: %s\nПродано за: %s\nТип оплаты: %s",
		e.Item,
		e.Material,
		e.Price,
		e.PaymentType,
	)
}
