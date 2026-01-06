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
	oauthClientID := os.Getenv("INPUT_OAUTH_CLIENT_ID")
	oauthClientSecret := os.Getenv("INPUT_OAUTH_CLIENT_SECRET")

	// Validate inputs
	if credentials == "" {
		log.Fatal("credentials is required")
	}
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

	// Validate OAuth refresh token inputs
	if strings.ToLower(authType) == "oauth_refresh_token" {
		if oauthClientID == "" {
			log.Fatal("oauth_client_id is required for oauth_refresh_token auth type")
		}
		if oauthClientSecret == "" {
			log.Fatal("oauth_client_secret is required for oauth_refresh_token auth type")
		}
	}

	// Log configuration (sanitized)
	log.Printf("Spreadsheet ID: %s", spreadsheetID)
	log.Printf("Sheet Name: %s", sheetName)
	log.Printf("Auth Type: %s", authType)
	log.Printf("Credentials length: %d bytes", len(credentials))

	// Parse values (expect JSON array)
	var rowValues []interface{}
	if err := json.Unmarshal([]byte(values), &rowValues); err != nil {
		log.Fatalf("Failed to parse values as JSON array: %v", err)
	}

	// Create context
	ctx := context.Background()

	// Create Sheets service based on auth type
	srv, err := createSheetsService(ctx, authType, credentials, oauthClientID, oauthClientSecret)
	if err != nil {
		log.Fatalf("Failed to create Sheets service: %v", err)
	}

	var existingRowCount int
	var nextRowNumber int
	var rangeNotation string

	// Read existing data to get current row count
	log.Println("📖 Reading existing data from sheet...")
	readRange := fmt.Sprintf("%s!A:Z", sheetName)

	readResp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Printf("⚠️  Warning: Failed to read existing data: %v", err)
		log.Println("Continuing with standard append...")
		rangeNotation = fmt.Sprintf("%s!A:Z", sheetName)
	} else {
		existingRowCount = len(readResp.Values)
		nextRowNumber = existingRowCount + 1

		log.Printf("📊 Existing data:")
		log.Printf("  Total rows: %d", existingRowCount)
		log.Printf("  Next row number: %d", nextRowNumber)

		if existingRowCount > 0 && len(readResp.Values[0]) > 0 {
			log.Printf("  Columns in first row: %d", len(readResp.Values[0]))
		}

		// Use specific range notation based on existing data
		rangeNotation = fmt.Sprintf("%s!A%d:Z%d", sheetName, nextRowNumber, nextRowNumber)
		log.Printf("  Appending to: %s", rangeNotation)
	}

	// Prepare the value range
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{rowValues},
	}

	// Append the row
	log.Println("📝 Appending row...")
	resp, err := srv.Spreadsheets.Values.Append(spreadsheetID, rangeNotation, valueRange).
		ValueInputOption("USER_ENTERED").
		InsertDataOption("INSERT_ROWS").
		Do()
	if err != nil {
		log.Fatalf("Failed to append row: %v", err)
	}

	// Calculate total rows after append
	totalRows := existingRowCount + int(resp.Updates.UpdatedRows)

	// Output results
	fmt.Printf("\n✅ Successfully appended row!\n")
	fmt.Printf("  Updated range: %s\n", resp.Updates.UpdatedRange)
	fmt.Printf("  Updated rows: %d\n", resp.Updates.UpdatedRows)
	fmt.Printf("  Updated columns: %d\n", resp.Updates.UpdatedColumns)
	fmt.Printf("  Updated cells: %d\n", resp.Updates.UpdatedCells)

	fmt.Printf("\n📊 Data length information:\n")
	fmt.Printf("  Existing rows before append: %d\n", existingRowCount)
	fmt.Printf("  Row number where data added: %d\n", nextRowNumber)
	fmt.Printf("  Total rows after append: %d\n", totalRows)

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

func createSheetsService(ctx context.Context, authType, credentials, clientID, clientSecret string) (*sheets.Service, error) {
	var tokenSource oauth2.TokenSource

	// Trim whitespace from credentials
	credentials = strings.TrimSpace(credentials)

	switch strings.ToLower(authType) {
	case "service_account":
		// Service Account authentication
		creds, err := google.CredentialsFromJSON(ctx, []byte(credentials), sheets.SpreadsheetsScope)
		if err != nil {
			return nil, fmt.Errorf("failed to parse service account credentials: %w", err)
		}
		tokenSource = creds.TokenSource

	case "oauth":
		// OAuth access token authentication
		accessToken := strings.TrimSpace(credentials)
		if accessToken == "" {
			return nil, fmt.Errorf("OAuth access token cannot be empty")
		}

		token := &oauth2.Token{
			AccessToken: accessToken,
		}
		tokenSource = oauth2.StaticTokenSource(token)

	case "oauth_refresh_token":
		// OAuth refresh token authentication
		refreshToken := strings.TrimSpace(credentials)
		clientID = strings.TrimSpace(clientID)
		clientSecret = strings.TrimSpace(clientSecret)

		if refreshToken == "" {
			return nil, fmt.Errorf("OAuth refresh token cannot be empty")
		}
		if clientID == "" {
			return nil, fmt.Errorf("OAuth client ID is required for refresh token authentication")
		}
		if clientSecret == "" {
			return nil, fmt.Errorf("OAuth client secret is required for refresh token authentication")
		}

		// Configure OAuth2
		config := &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     google.Endpoint,
			Scopes:       []string{sheets.SpreadsheetsScope},
		}

		// Create token with refresh token
		token := &oauth2.Token{
			RefreshToken: refreshToken,
		}

		// This will automatically refresh the token when needed
		tokenSource = config.TokenSource(ctx, token)

		log.Println("Using OAuth refresh token authentication")

	default:
		return nil, fmt.Errorf("unsupported auth_type: %s (supported: service_account, oauth, oauth_refresh_token)", authType)
	}

	// Create the Sheets service
	srv, err := sheets.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("failed to create Sheets service: %w", err)
	}

	return srv, nil
}
