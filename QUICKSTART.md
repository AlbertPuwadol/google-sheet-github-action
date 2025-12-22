# Quick Start Guide

Get up and running with the Google Sheet GitHub Action in 5 minutes!

## Step 1: Setup Your Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project (or use an existing one)
3. Enable the **Google Sheets API**:
   - Navigate to "APIs & Services" > "Library"
   - Search for "Google Sheets API"
   - Click "Enable"

## Step 2: Create Service Account

1. Go to "APIs & Services" > "Credentials"
2. Click "Create Credentials" > "Service Account"
3. Enter a name (e.g., "github-actions-sheets")
4. Click "Create and Continue"
5. Skip the optional steps and click "Done"

## Step 3: Download JSON Credentials

1. Click on your newly created service account
2. Go to the "Keys" tab
3. Click "Add Key" > "Create new key"
4. Select "JSON" format
5. Click "Create" - the JSON file will download automatically

## Step 4: Prepare Your Google Sheet

1. Open (or create) your Google Sheet
2. Note the **Spreadsheet ID** from the URL:
   ```
   https://docs.google.com/spreadsheets/d/SPREADSHEET_ID_HERE/edit
   ```
3. Click the "Share" button
4. Add the service account email (found in the JSON file as `client_email`)
   - It looks like: `github-actions-sheets@project-id.iam.gserviceaccount.com`
5. Give it "Editor" permission
6. Click "Send"

## Step 5: Add GitHub Secrets

1. Go to your GitHub repository
2. Navigate to **Settings** > **Secrets and variables** > **Actions**
3. Add these secrets:

   **Secret 1:** `GOOGLE_SERVICE_ACCOUNT_JSON`
   - Click "New repository secret"
   - Name: `GOOGLE_SERVICE_ACCOUNT_JSON`
   - Value: Paste the **entire contents** of the downloaded JSON file
   - Click "Add secret"

   **Secret 2:** `SPREADSHEET_ID`
   - Click "New repository secret"
   - Name: `SPREADSHEET_ID`
   - Value: The ID from your Google Sheets URL
   - Click "Add secret"

## Step 6: Use the Action in Your Workflow

Create a file `.github/workflows/log-to-sheet.yml` in your repository:

```yaml
name: Log to Google Sheet

on:
  push:
    branches: [ main ]

jobs:
  log:
    runs-on: ubuntu-latest
    steps:
      - name: Append to Google Sheet
        uses: AlbertPuwadol/google-sheet-github-action@v1
        with:
          spreadsheet_id: ${{ secrets.SPREADSHEET_ID }}
          sheet_name: 'Sheet1'
          values: '["${{ github.sha }}", "${{ github.actor }}", "${{ github.event.head_commit.message }}"]'
          auth_type: 'service_account'
          credentials: ${{ secrets.GOOGLE_SERVICE_ACCOUNT_JSON }}
```

## Step 7: Test It!

1. Commit and push the workflow file
2. Check the "Actions" tab in your GitHub repository
3. The workflow should run and append a row to your Google Sheet

## Troubleshooting

### ❌ "Permission denied" error
- Make sure you shared the sheet with the service account email
- Check that the service account has "Editor" permission

### ❌ "Spreadsheet not found" error
- Verify the spreadsheet ID is correct
- Ensure the service account has access

### ❌ "Sheet not found" error
- Check that `sheet_name` matches exactly (case-sensitive)
- The sheet tab must exist

## Next Steps

- Read the full [README.md](README.md) for more examples
- Check out the example workflows in `.github/workflows/`
- Customize the `values` to log different information

## Need Help?

Open an issue on [GitHub](https://github.com/AlbertPuwadol/google-sheet-github-action/issues)

