package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func main() {
	// Get inputs from environment variables (GitHub Actions inputs)
	spreadsheetID := os.Getenv("INPUT_SPREADSHEET_ID")
	sheetName := os.Getenv("INPUT_SHEET_NAME")
	values := os.Getenv("INPUT_VALUES")
	authType := os.Getenv("INPUT_AUTH_TYPE")
	credentials := os.Getenv("INPUT_CREDENTIALS")

	// Validate inputs
	if spreadsheetID == "" {
		log.Fatal("spreadsheet_id is required")
	}
	if sheetName == "" {
		sheetName = "Sheet1" // Default to Sheet1
	}
	if values == "" {
		log.Fatal("values is required")
	}
	if authType == "" {
		authType = "service_account" // Default to service account
	}
	if credentials == "" {
		log.Fatal("credentials is required")
	}

	// Parse values (expect JSON array)
	var rowValues []interface{}
	if err := json.Unmarshal([]byte(values), &rowValues); err != nil {
		log.Fatalf("Failed to parse values as JSON array: %v", err)
	}

	// Create context
	ctx := context.Background()

	// Create Sheets service based on auth type
	srv, err := createSheetsService(ctx, authType, credentials)
	if err != nil {
		log.Fatalf("Failed to create Sheets service: %v", err)
	}

	// Prepare the range
	rangeNotation := fmt.Sprintf("%s!A:Z", sheetName)

	// Prepare the value range
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{rowValues},
	}

	// Append the row
	resp, err := srv.Spreadsheets.Values.Append(spreadsheetID, rangeNotation, valueRange).
		ValueInputOption("USER_ENTERED").
		InsertDataOption("INSERT_ROWS").
		Do()
	if err != nil {
		log.Fatalf("Failed to append row: %v", err)
	}

	// Output results
	fmt.Printf("Successfully appended row to %s\n", resp.TableRange)
	fmt.Printf("Updated range: %s\n", resp.Updates.UpdatedRange)
	fmt.Printf("Updated rows: %d\n", resp.Updates.UpdatedRows)
	fmt.Printf("Updated columns: %d\n", resp.Updates.UpdatedColumns)
	fmt.Printf("Updated cells: %d\n", resp.Updates.UpdatedCells)

	// Set output for GitHub Actions
	githubOutput := os.Getenv("GITHUB_OUTPUT")
	if githubOutput != "" {
		f, err := os.OpenFile(githubOutput, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Warning: Failed to open GITHUB_OUTPUT file: %v", err)
		} else {
			defer f.Close()
			fmt.Fprintf(f, "updated_range=%s\n", resp.Updates.UpdatedRange)
			fmt.Fprintf(f, "updated_rows=%d\n", resp.Updates.UpdatedRows)
		}
	}
}

func createSheetsService(ctx context.Context, authType, credentials string) (*sheets.Service, error) {
	var tokenSource oauth2.TokenSource

	switch strings.ToLower(authType) {
	case "service_account":
		// Service Account authentication
		creds, err := google.CredentialsFromJSON(ctx, []byte(credentials), sheets.SpreadsheetsScope)
		if err != nil {
			return nil, fmt.Errorf("failed to parse service account credentials: %w", err)
		}
		tokenSource = creds.TokenSource

	case "oauth":
		// OAuth token authentication
		token := &oauth2.Token{
			AccessToken: credentials,
		}
		tokenSource = oauth2.StaticTokenSource(token)

	default:
		return nil, fmt.Errorf("unsupported auth_type: %s (supported: service_account, oauth)", authType)
	}

	// Create the Sheets service
	srv, err := sheets.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("failed to create Sheets service: %w", err)
	}

	return srv, nil
}
