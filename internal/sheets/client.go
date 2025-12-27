package sheets

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Config holds the configuration for a Google Sheets client
type Config struct {
	SpreadsheetID string
	SheetName     string
	ColumnRange   string // e.g., "A:E" for 5 columns, "A:D" for 4 columns
}

// Client wraps the Google Sheets API client
type Client struct {
	service     *sheets.Service
	spreadsheet string
	sheetName   string
	columnRange string
}

// NewClient creates a new Google Sheets client with the given configuration
// Expects GOOGLE_APPLICATION_CREDENTIALS environment variable pointing to service account JSON
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.SpreadsheetID == "" {
		return nil, fmt.Errorf("SpreadsheetID is required")
	}

	if cfg.SheetName == "" {
		cfg.SheetName = "Sheet1"
	}

	if cfg.ColumnRange == "" {
		cfg.ColumnRange = "A:E"
	}

	var service *sheets.Service
	var err error

	// Check for service account credentials file
	credsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credsFile != "" {
		service, err = sheets.NewService(ctx, option.WithCredentialsFile(credsFile))
		if err != nil {
			return nil, fmt.Errorf("failed to create sheets service with credentials file: %w", err)
		}
	} else {
		// Fall back to Application Default Credentials
		service, err = sheets.NewService(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create sheets service: %w", err)
		}
	}

	return &Client{
		service:     service,
		spreadsheet: cfg.SpreadsheetID,
		sheetName:   cfg.SheetName,
		columnRange: cfg.ColumnRange,
	}, nil
}

// AppendRow appends a row of data to the Google Sheet
func (c *Client) AppendRow(ctx context.Context, row []interface{}) error {
	rangeNotation := fmt.Sprintf("%s!%s", c.sheetName, c.columnRange)

	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{row},
	}

	_, err := c.service.Spreadsheets.Values.Append(
		c.spreadsheet,
		rangeNotation,
		valueRange,
	).ValueInputOption("USER_ENTERED").
		InsertDataOption("INSERT_ROWS").
		Context(ctx).
		Do()

	if err != nil {
		return fmt.Errorf("failed to append row to sheet: %w", err)
	}

	return nil
}
