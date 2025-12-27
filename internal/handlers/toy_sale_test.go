package handlers

import (
	"strings"
	"testing"
)

func TestToySaleEvent_toSheetRow(t *testing.T) {
	event := ToySaleEvent{
		Time:        "2025-11-18 22:45:48",
		Item:        "Марк",
		Material:    "Золотой",
		Price:       "40",
		PaymentType: "Карта",
	}

	row := event.toSheetRow()

	if len(row) != 5 {
		t.Errorf("Expected 5 columns, got %d", len(row))
	}

	// Verify column order: Время продажи | Фигурка | Пластик | Цена | Тип оплаты
	if row[0] != "2025-11-18 22:45:48" {
		t.Errorf("Expected time '2025-11-18 22:45:48', got '%v'", row[0])
	}

	if row[1] != "Марк" {
		t.Errorf("Expected item 'Марк', got '%v'", row[1])
	}

	if row[2] != "Золотой" {
		t.Errorf("Expected material 'Золотой', got '%v'", row[2])
	}

	if row[3] != "40" {
		t.Errorf("Expected price '40', got '%v'", row[3])
	}

	if row[4] != "Карта" {
		t.Errorf("Expected paymentType 'Карта', got '%v'", row[4])
	}
}

func TestToySaleEvent_toTelegramMessage(t *testing.T) {
	event := ToySaleEvent{
		Time:        "2025-11-18 22:45:48",
		Item:        "Марк",
		Material:    "Золотой",
		Price:       "40",
		PaymentType: "Карта",
	}

	message := event.toTelegramMessage()

	expectedMessage := "Фигурка: Марк\nПластик: Золотой\nПродано за: 40\nТип оплаты: Карта"

	if message != expectedMessage {
		t.Errorf("Expected message:\n%s\n\nGot:\n%s", expectedMessage, message)
	}
}

func TestToySaleEvent_toTelegramMessage_ContainsAllFields(t *testing.T) {
	event := ToySaleEvent{
		Time:        "2025-11-18 22:45:48",
		Item:        "Тест",
		Material:    "Серебряный",
		Price:       "100",
		PaymentType: "Наличные",
	}

	message := event.toTelegramMessage()

	if !strings.Contains(message, "Фигурка: Тест") {
		t.Error("Message should contain 'Фигурка: Тест'")
	}

	if !strings.Contains(message, "Пластик: Серебряный") {
		t.Error("Message should contain 'Пластик: Серебряный'")
	}

	if !strings.Contains(message, "Продано за: 100") {
		t.Error("Message should contain 'Продано за: 100'")
	}

	if !strings.Contains(message, "Тип оплаты: Наличные") {
		t.Error("Message should contain 'Тип оплаты: Наличные'")
	}
}
