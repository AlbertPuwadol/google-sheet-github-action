# OAuth Refresh Token Setup Guide

This guide shows you how to set up OAuth refresh token authentication for the Google Sheets Action. This is useful when you want to authenticate as a specific user rather than using a service account.

## Why Use OAuth Refresh Tokens?

- ✅ **Long-lived:** Refresh tokens don't expire (unless revoked)
- ✅ **User context:** Actions performed as a specific user
- ✅ **No service account needed:** Use your personal Google account
- ✅ **Automatic refresh:** The action automatically gets new access tokens when needed

## Prerequisites

- A Google account with access to the Google Sheet
- A Google Cloud project with Sheets API enabled

## Step 1: Create OAuth Client Credentials

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Select your project (or create a new one)
3. Navigate to **APIs & Services** > **Credentials**
4. Click **Create Credentials** > **OAuth client ID**
5. If prompted, configure the OAuth consent screen:
   - User Type: Choose "External" (for personal use) or "Internal" (for workspace)
   - Fill in the required fields (App name, user support email, developer contact)
   - Add scopes: `https://www.googleapis.com/auth/spreadsheets`
   - Add test users (your email) if using External
6. Back to Create OAuth client ID:
   - Application type: **Desktop app** (or Web application)
   - Name: "GitHub Actions Sheet Updater" (or any name)
   - Click **Create**
7. Download the JSON file or copy the **Client ID** and **Client Secret**

## Step 2: Get Your Refresh Token

### Option A: Using OAuth 2.0 Playground (Easiest)

1. Go to [Google OAuth 2.0 Playground](https://developers.google.com/oauthplayground/)

2. Click the gear icon (⚙️) in the top right

3. Check **"Use your own OAuth credentials"**

4. Enter your **OAuth Client ID** and **Client Secret** from Step 1

5. In the left sidebar under "Select & authorize APIs":

   - Scroll down or search for "Sheets API v4"
   - Select `https://www.googleapis.com/auth/spreadsheets`
   - Click **Authorize APIs**

6. Sign in with your Google account and grant permissions

7. You'll be redirected back to the playground

8. Click **Exchange authorization code for tokens**

9. Copy the **Refresh token** value - this is what you need!

### Option B: Using Command Line (Advanced)

1. Install Google's OAuth2 CLI tool:

```bash
pip install google-auth-oauthlib google-auth-httplib2
```

2. Create a file `get_refresh_token.py`:

```python
from google_auth_oauthlib.flow import InstalledAppFlow

SCOPES = ['https://www.googleapis.com/auth/spreadsheets']

# Replace with your client ID and secret
client_config = {
    "installed": {
        "client_id": "YOUR_CLIENT_ID",
        "client_secret": "YOUR_CLIENT_SECRET",
        "auth_uri": "https://accounts.google.com/o/oauth2/auth",
        "token_uri": "https://oauth2.googleapis.com/token",
        "redirect_uris": ["http://localhost"]
    }
}

flow = InstalledAppFlow.from_client_config(client_config, SCOPES)
creds = flow.run_local_server(port=0)

print("\n✅ Authentication successful!")
print(f"\nRefresh Token:\n{creds.refresh_token}")
print(f"\nAccess Token:\n{creds.token}")
```

3. Run the script:

```bash
python get_refresh_token.py
```

4. Your browser will open - sign in and grant permissions

5. Copy the refresh token from the output

## Step 3: Add Secrets to GitHub

Go to your GitHub repository → **Settings** → **Secrets and variables** → **Actions**

Add these three secrets:

### 1. `GOOGLE_OAUTH_REFRESH_TOKEN`

- Click "New repository secret"
- Name: `GOOGLE_OAUTH_REFRESH_TOKEN`
- Value: Paste your refresh token
- Click "Add secret"

### 2. `GOOGLE_OAUTH_CLIENT_ID`

- Click "New repository secret"
- Name: `GOOGLE_OAUTH_CLIENT_ID`
- Value: Paste your OAuth Client ID
- Click "Add secret"

### 3. `GOOGLE_OAUTH_CLIENT_SECRET`

- Click "New repository secret"
- Name: `GOOGLE_OAUTH_CLIENT_SECRET`
- Value: Paste your OAuth Client Secret
- Click "Add secret"

### 4. `SPREADSHEET_ID`

- Click "New repository secret"
- Name: `SPREADSHEET_ID`
- Value: Your Google Sheets ID (from the URL)
- Click "Add secret"

## Step 4: Use in Your Workflow

Create a workflow file `.github/workflows/update-sheet.yml`:

```yaml
name: Update Google Sheet

on:
  push:
    branches: [main]

jobs:
  update-sheet:
    runs-on: ubuntu-latest
    steps:
      - name: Append row to Google Sheet
        uses: AlbertPuwadol/google-sheet-github-action@main
        with:
          spreadsheet_id: ${{ secrets.SPREADSHEET_ID }}
          sheet_name: "Sheet1"
          values: '["${{ github.sha }}", "${{ github.actor }}", "${{ github.event.head_commit.message }}"]'
          auth_type: "oauth_refresh_token"
          credentials: ${{ secrets.GOOGLE_OAUTH_REFRESH_TOKEN }}
          oauth_client_id: ${{ secrets.GOOGLE_OAUTH_CLIENT_ID }}
          oauth_client_secret: ${{ secrets.GOOGLE_OAUTH_CLIENT_SECRET }}
```

## Step 5: Test It!

1. Commit and push your workflow file
2. Check the "Actions" tab in your GitHub repository
3. The workflow should run and append a row to your Google Sheet
4. Check your Google Sheet to verify the row was added

## Comparison: Service Account vs OAuth Refresh Token

| Feature              | Service Account                       | OAuth Refresh Token                        |
| -------------------- | ------------------------------------- | ------------------------------------------ |
| **Setup complexity** | Medium                                | Medium-High                                |
| **Expiration**       | Never (unless key rotated)            | Never (unless revoked)                     |
| **User context**     | Service account email                 | Your personal account                      |
| **Permissions**      | Must share sheet with service account | Uses your permissions                      |
| **Best for**         | Automated systems, shared resources   | Personal automation, user-specific actions |
| **Secrets needed**   | 1 (JSON file)                         | 3 (refresh token, client ID, secret)       |

## Troubleshooting

### ❌ "Invalid grant" error

- Your refresh token may have expired or been revoked
- Generate a new refresh token using the OAuth 2.0 Playground

### ❌ "Invalid client" error

- Check that your Client ID and Client Secret are correct
- Make sure you're using credentials from the same project

### ❌ "Access denied" error

- Verify you granted the `spreadsheets` scope when generating the refresh token
- Check that the user has access to the spreadsheet

### ❌ "Token has been expired or revoked"

- Refresh tokens can be revoked if:
  - You changed your Google password
  - You revoked access in Google Account settings
  - The OAuth consent screen is in "Testing" mode and 7 days have passed
- Solution: Generate a new refresh token

## Security Notes

⚠️ **Important Security Considerations:**

1. **Keep secrets secure:** Never commit refresh tokens, client IDs, or secrets to your repository
2. **Limit access:** Only add test users who need access in the OAuth consent screen
3. **Regular rotation:** Consider rotating credentials periodically
4. **Monitor usage:** Check Google Cloud Console for unusual API usage
5. **Revoke if compromised:** If credentials are exposed, immediately revoke them in Google Cloud Console

## Revoking Access

To revoke access:

1. Go to [Google Account Permissions](https://myaccount.google.com/permissions)
2. Find your app in the list
3. Click "Remove Access"
4. The refresh token will no longer work

Or revoke in Google Cloud Console:

1. Go to [API Credentials](https://console.cloud.google.com/apis/credentials)
2. Find your OAuth 2.0 Client ID
3. Delete or regenerate credentials

## Additional Resources

- [Google OAuth 2.0 Documentation](https://developers.google.com/identity/protocols/oauth2)
- [Google Sheets API Reference](https://developers.google.com/sheets/api)
- [OAuth 2.0 Playground](https://developers.google.com/oauthplayground/)
- [Google Cloud Console](https://console.cloud.google.com/)

## Need Help?

Open an issue on [GitHub](https://github.com/AlbertPuwadol/google-sheet-github-action/issues)
