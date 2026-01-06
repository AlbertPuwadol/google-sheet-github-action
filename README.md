# Google Sheets Append Row - GitHub Action

A GitHub Action that appends a row to a Google Sheet using the official Google Sheets API. Written in Go using the [google-api-go-client](https://github.com/googleapis/google-api-go-client) library.

📖 **[Quick Reference](QUICK_REFERENCE.md)** | 🔐 **[Centralized Credentials](CENTRALIZED_CREDENTIALS.md)** | 🔑 **[OAuth Setup](OAUTH_SETUP.md)**

## Features

- ✅ Append rows to Google Sheets from GitHub Actions workflows
- ✅ Multiple authentication methods:
  - **Service Account** (recommended for automation)
  - **OAuth Access Token** (short-lived)
  - **OAuth Refresh Token** (long-lived, auto-refreshing)
- ✅ Multiple credential storage options:
  - **GitHub Secrets** (per-repository)
  - **Organization Secrets** (shared across repos)
  - **Reusable Workflow** (centralized credentials) 🆕
- ✅ Easy to use from any GitHub repository
- ✅ Built with Go for fast execution
- ✅ Outputs updated range and row count

## Usage

### 🎯 Quick Start: Reusable Workflow (Centralized Credentials)

**Store credentials once, use from any repository!**

In any repository, create `.github/workflows/log-to-sheet.yml`:

```yaml
name: Log to Google Sheet

on:
  push:
    branches: [main]

jobs:
  log:
    uses: AlbertPuwadol/google-sheet-github-action/.github/workflows/reusable-append-row.yml@main
    with:
      sheet_name: "Activity Log"
      values: '["${{ github.repository }}", "${{ github.sha }}", "${{ github.actor }}"]'
      use_action_repo_credentials: true
```

**📖 Setup Guide:** [CENTRALIZED_CREDENTIALS.md](CENTRALIZED_CREDENTIALS.md)

---

### Basic Example (Service Account)

```yaml
- name: Append row to Google Sheet
  uses: AlbertPuwadol/google-sheet-github-action@v1
  with:
    spreadsheet_id: "YOUR_SPREADSHEET_ID"
    sheet_name: "Sheet1"
    values: '["Column 1", "Column 2", "Column 3"]'
    auth_type: "service_account"
    credentials: ${{ secrets.GOOGLE_SERVICE_ACCOUNT_JSON }}
```

### OAuth Access Token Example

```yaml
- name: Append row to Google Sheet
  uses: AlbertPuwadol/google-sheet-github-action@v1
  with:
    spreadsheet_id: "YOUR_SPREADSHEET_ID"
    sheet_name: "Sheet1"
    values: '["Data 1", "Data 2", 123]'
    auth_type: "oauth"
    credentials: ${{ secrets.GOOGLE_OAUTH_TOKEN }}
```

### OAuth Refresh Token Example (Recommended for OAuth)

```yaml
- name: Append row to Google Sheet
  uses: AlbertPuwadol/google-sheet-github-action@v1
  with:
    spreadsheet_id: ${{ secrets.SPREADSHEET_ID }}
    sheet_name: "Sheet1"
    values: '["Data 1", "Data 2", 123]'
    auth_type: "oauth_refresh_token"
    credentials: ${{ secrets.GOOGLE_OAUTH_REFRESH_TOKEN }}
    oauth_client_id: ${{ secrets.GOOGLE_OAUTH_CLIENT_ID }}
    oauth_client_secret: ${{ secrets.GOOGLE_OAUTH_CLIENT_SECRET }}
```

## Inputs

| Input                 | Description                                                                                     | Required | Default           |
| --------------------- | ----------------------------------------------------------------------------------------------- | -------- | ----------------- |
| `spreadsheet_id`      | The ID of the Google Spreadsheet (found in the URL)                                             | Yes      | -                 |
| `sheet_name`          | The name of the sheet/tab within the spreadsheet                                                | No       | `Sheet1`          |
| `values`              | JSON array of values to append as a row                                                         | Yes      | -                 |
| `auth_type`           | Authentication type: `service_account`, `oauth`, or `oauth_refresh_token`                       | No       | `service_account` |
| `credentials`         | For service_account: Full JSON. For oauth: Access token. For oauth_refresh_token: Refresh token | Yes      | -                 |
| `oauth_client_id`     | OAuth Client ID (required only for `oauth_refresh_token`)                                       | No       | -                 |
| `oauth_client_secret` | OAuth Client Secret (required only for `oauth_refresh_token`)                                   | No       | -                 |

## Outputs

| Output          | Description                                   |
| --------------- | --------------------------------------------- |
| `updated_range` | The A1 notation of the range that was updated |
| `updated_rows`  | Number of rows that were updated              |

## Authentication Methods Comparison

Choose the authentication method that best fits your needs:

| Feature               | Service Account          | OAuth Access Token     | OAuth Refresh Token      |
| --------------------- | ------------------------ | ---------------------- | ------------------------ |
| **Setup Complexity**  | Medium                   | Low                    | Medium                   |
| **Expiration**        | Never                    | ~1 hour                | Never (auto-refreshes)   |
| **Best for**          | Automation, CI/CD        | Quick testing          | User-specific automation |
| **GitHub Secrets**    | 1 (JSON)                 | 1 (token)              | 3 (token, ID, secret)    |
| **User Context**      | Service account email    | User who granted token | User who granted token   |
| **Sheet Permissions** | Must share with SA email | Uses user's access     | Uses user's access       |
| **Token Refresh**     | N/A                      | Manual                 | Automatic                |
| **Recommended for**   | ✅ Production workflows  | ❌ Testing only        | ✅ Personal automation   |

**Quick Recommendation:**

- 🏢 **Team/Organization:** Use **Service Account**
- 👤 **Personal Projects:** Use **OAuth Refresh Token**
- 🧪 **Quick Testing:** Use **OAuth Access Token**

## Setup

### Option 1: Service Account Authentication (Recommended)

1. **Create a Google Cloud Project** (if you don't have one)

   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project or select an existing one

2. **Enable Google Sheets API**

   - Navigate to "APIs & Services" > "Library"
   - Search for "Google Sheets API"
   - Click "Enable"

3. **Create a Service Account**

   - Go to "APIs & Services" > "Credentials"
   - Click "Create Credentials" > "Service Account"
   - Fill in the service account details
   - Click "Create and Continue"
   - Skip granting roles (optional)
   - Click "Done"

4. **Create and Download JSON Key**

   - Click on the created service account
   - Go to "Keys" tab
   - Click "Add Key" > "Create new key"
   - Select "JSON" format
   - Download the JSON file

5. **Share Google Sheet with Service Account**

   - Open your Google Sheet
   - Click "Share" button
   - Add the service account email (found in the JSON file as `client_email`)
   - Give it "Editor" permission

6. **Add to GitHub Secrets**

   - Go to your GitHub repository
   - Navigate to Settings > Secrets and variables > Actions

   **Add Service Account Credentials:**

   - Click "New repository secret"
   - Name: `GOOGLE_SERVICE_ACCOUNT_JSON`
   - Value: Open the JSON file in a text editor, select ALL content (Ctrl+A/Cmd+A), and paste
   - ⚠️ **Important:** Make sure you copy the entire JSON including `{` and `}`, with no extra spaces at the beginning or end
   - Click "Add secret"

   **Add Spreadsheet ID:**

   - Click "New repository secret"
   - Name: `SPREADSHEET_ID`
   - Value: The ID from your Google Sheets URL (the long string between `/d/` and `/edit`)
   - Example: In `https://docs.google.com/spreadsheets/d/1abc...xyz/edit`, the ID is `1abc...xyz`
   - Click "Add secret"

### Option 2: OAuth Token Authentication

#### OAuth Access Token (Short-lived)

1. **Obtain an OAuth Access Token**

   - Use Google OAuth 2.0 Playground or your own OAuth flow
   - Request access to `https://www.googleapis.com/auth/spreadsheets` scope

2. **Add to GitHub Secrets**
   - Name: `GOOGLE_OAUTH_TOKEN`
   - Value: Your OAuth access token

> **Note**: OAuth access tokens typically expire after 1 hour. Not recommended for production workflows.

#### OAuth Refresh Token (Long-lived, Recommended)

OAuth refresh tokens don't expire and automatically refresh access tokens when needed.

**📖 See the detailed guide:** [OAUTH_SETUP.md](OAUTH_SETUP.md)

**Quick setup:**

1. **Create OAuth Client Credentials** in Google Cloud Console
2. **Get a Refresh Token** using OAuth 2.0 Playground or CLI
3. **Add these secrets to GitHub:**

   - `GOOGLE_OAUTH_REFRESH_TOKEN` - Your refresh token
   - `GOOGLE_OAUTH_CLIENT_ID` - OAuth client ID
   - `GOOGLE_OAUTH_CLIENT_SECRET` - OAuth client secret
   - `SPREADSHEET_ID` - Your spreadsheet ID

4. **Use in workflow:**

```yaml
- uses: AlbertPuwadol/google-sheet-github-action@v1
  with:
    spreadsheet_id: ${{ secrets.SPREADSHEET_ID }}
    sheet_name: "Sheet1"
    values: '["data"]'
    auth_type: "oauth_refresh_token"
    credentials: ${{ secrets.GOOGLE_OAUTH_REFRESH_TOKEN }}
    oauth_client_id: ${{ secrets.GOOGLE_OAUTH_CLIENT_ID }}
    oauth_client_secret: ${{ secrets.GOOGLE_OAUTH_CLIENT_SECRET }}
```

## Examples

### Log Deployment Information

```yaml
name: Log Deployment

on:
  push:
    branches: [main]

jobs:
  log-deployment:
    runs-on: ubuntu-latest
    steps:
      - name: Log to Google Sheet
        uses: AlbertPuwadol/google-sheet-github-action@v1
        with:
          spreadsheet_id: ${{ secrets.SPREADSHEET_ID }}
          sheet_name: "Deployments"
          values: '["${{ github.sha }}", "${{ github.actor }}", "${{ github.event.head_commit.message }}", "${{ github.event.repository.updated_at }}"]'
          auth_type: "service_account"
          credentials: ${{ secrets.GOOGLE_SERVICE_ACCOUNT_JSON }}
```

### Track Test Results

```yaml
name: Test and Track

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run tests
        id: test
        run: |
          npm test
          echo "status=passed" >> $GITHUB_OUTPUT
        continue-on-error: true

      - name: Log test results
        uses: AlbertPuwadol/google-sheet-github-action@v1
        with:
          spreadsheet_id: ${{ secrets.SPREADSHEET_ID }}
          sheet_name: "Test Results"
          values: '["${{ github.run_number }}", "${{ steps.test.outputs.status }}", "${{ github.ref }}", "${{ github.event.head_commit.timestamp }}"]'
          auth_type: "service_account"
          credentials: ${{ secrets.GOOGLE_SERVICE_ACCOUNT_JSON }}
```

### Using Outputs

```yaml
- name: Append row to Google Sheet
  id: append
  uses: AlbertPuwadol/google-sheet-github-action@v1
  with:
    spreadsheet_id: ${{ secrets.SPREADSHEET_ID }}
    sheet_name: "Sheet1"
    values: '["Data"]'
    auth_type: "service_account"
    credentials: ${{ secrets.GOOGLE_SERVICE_ACCOUNT_JSON }}

- name: Print result
  run: |
    echo "Updated range: ${{ steps.append.outputs.updated_range }}"
    echo "Updated rows: ${{ steps.append.outputs.updated_rows }}"
```

## Values Format

The `values` input expects a JSON array. Each element in the array represents a column value:

```yaml
# Simple values
values: '["Text", 123, true]'

# With variables
values: '["${{ github.actor }}", "${{ github.run_number }}", "success"]'

# Complex example
values: '["Deployment", "${{ github.sha }}", "${{ github.event.repository.full_name }}", "${{ github.event.head_commit.timestamp }}"]'
```

## Troubleshooting

### Permission Denied Error

- Make sure you've shared the Google Sheet with the service account email
- Verify the service account has "Editor" permission

### Invalid Credentials / Invalid Header Field Value

This error can occur if credentials are not properly formatted:

**For Service Account (JSON):**

- Open the downloaded JSON file in a text editor
- Copy the **entire** contents (including opening `{` and closing `}`)
- Paste directly into GitHub Secrets - do not add extra spaces or newlines
- The JSON should start with `{` and end with `}`
- Make sure there are no trailing spaces or formatting issues

**For OAuth Token:**

- Ensure the token is a single line with no spaces or newlines
- The token should be just the access token string, nothing else
- Common issue: copying the token with quotes or extra whitespace

**Quick fix:**

- Delete and recreate the secret in GitHub
- Use `cat credentials.json | pbcopy` (Mac) or `cat credentials.json | xclip` (Linux) to copy the exact file contents
- Paste directly into GitHub Secrets without any modification
- **Tip:** Run `./validate-credentials.sh your-credentials.json` to verify your credentials file before adding to GitHub

### Spreadsheet Not Found

- Verify the `spreadsheet_id` is correct (from the URL)
- Ensure the service account has access to the spreadsheet

### Sheet Not Found

- Check that the `sheet_name` matches exactly (case-sensitive)
- The sheet tab must exist in the spreadsheet

## Development

### Local Testing

1. Clone the repository:

```bash
git clone https://github.com/AlbertPuwadol/google-sheet-github-action.git
cd google-sheet-github-action
```

2. Install dependencies:

```bash
go mod download
```

3. Set environment variables:

```bash
export INPUT_SPREADSHEET_ID="your-spreadsheet-id"
export INPUT_SHEET_NAME="Sheet1"
export INPUT_VALUES='["Test", "Data", 123]'
export INPUT_AUTH_TYPE="service_account"
export INPUT_CREDENTIALS='{"type":"service_account",...}'
```

4. Run the program:

```bash
go run main.go
```

### Building Docker Image

```bash
docker build -t google-sheet-action .
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.

## Related Projects

- [google-api-go-client](https://github.com/googleapis/google-api-go-client) - Official Google APIs Client Library for Go
- [GitHub Actions Documentation](https://docs.github.com/en/actions)

## Support

If you encounter any issues or have questions, please [open an issue](https://github.com/AlbertPuwadol/google-sheet-github-action/issues).
