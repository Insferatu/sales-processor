package handlers

import (
	"strings"
	"testing"
)

func TestJewelrySaleEvent_toSheetRow(t *testing.T) {
	event := JewelrySaleEvent{
		Time:        "2025-11-18 22:45:48",
		Product:     "Браслет",
		Price:       "40",
		PaymentType: "Карта",
	}

	row := event.toSheetRow()

	if len(row) != 4 {
		t.Errorf("Expected 4 columns, got %d", len(row))
	}

	// Verify column order: Время продажи | Товар | Цена | Тип оплаты
	if row[0] != "2025-11-18 22:45:48" {
		t.Errorf("Expected time '2025-11-18 22:45:48', got '%v'", row[0])
	}

	if row[1] != "Браслет" {
		t.Errorf("Expected product 'Браслет', got '%v'", row[1])
	}

	if row[2] != "40" {
		t.Errorf("Expected price '40', got '%v'", row[2])
	}

	if row[3] != "Карта" {
		t.Errorf("Expected paymentType 'Карта', got '%v'", row[3])
	}
}

func TestJewelrySaleEvent_toTelegramMessage(t *testing.T) {
	event := JewelrySaleEvent{
		Time:        "2025-11-18 22:45:48",
		Product:     "Браслет",
		Price:       "40",
		PaymentType: "Карта",
	}

	message := event.toTelegramMessage()

	expectedMessage := "Товар: Браслет\nПродано за: 40\nТип оплаты: Карта"

	if message != expectedMessage {
		t.Errorf("Expected message:\n%s\n\nGot:\n%s", expectedMessage, message)
	}
}

func TestJewelrySaleEvent_toTelegramMessage_ContainsAllFields(t *testing.T) {
	event := JewelrySaleEvent{
		Time:        "2025-11-18 22:45:48",
		Product:     "Кольцо",
		Price:       "200",
		PaymentType: "Наличные",
	}

	message := event.toTelegramMessage()

	if !strings.Contains(message, "Товар: Кольцо") {
		t.Error("Message should contain 'Товар: Кольцо'")
	}

	if !strings.Contains(message, "Продано за: 200") {
		t.Error("Message should contain 'Продано за: 200'")
	}

	if !strings.Contains(message, "Тип оплаты: Наличные") {
		t.Error("Message should contain 'Тип оплаты: Наличные'")
	}
}
