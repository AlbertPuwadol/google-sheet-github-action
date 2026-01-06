# Quick Reference Guide

## 🚀 Ways to Use This Action

### 1️⃣ Repository Secrets (Standard)

**Add secrets to repository and use in workflows**

```yaml
- name: Append to sheet
  id: append
  uses: AlbertPuwadol/google-sheet-github-action@main
  with:
    spreadsheet_id: ${{ secrets.SPREADSHEET_ID }}
    credentials: ${{ secrets.GOOGLE_SERVICE_ACCOUNT_JSON }}
    sheet_name: "Logs"
    values: '["data"]'

- name: Show results
  run: |
    echo "Added to row: ${{ steps.append.outputs.row_number }}"
    echo "Total rows: ${{ steps.append.outputs.total_rows }}"
```

**Setup:** Settings → Secrets and variables → Actions

**The action automatically:**

- ✅ Reads existing data from the sheet
- ✅ Determines next row number
- ✅ Uses precise range notation
- ✅ Logs data length information

---

### 2️⃣ Organization Secrets (For Teams)

**Set up once in Organization → use in all repos**

```yaml
- uses: AlbertPuwadol/google-sheet-github-action@main
  with:
    spreadsheet_id: ${{ secrets.ORG_SPREADSHEET_ID }}
    credentials: ${{ secrets.ORG_GOOGLE_CREDENTIALS }}
    sheet_name: "Logs"
    values: '["data"]'
```

**Requires:** GitHub Organization plan

## 🎯 Recommendations

- **🏢 GitHub Organization:** Use **Organization Secrets**
- **📦 Single repository:** Use **Repository Secrets**
- **👥 Multiple repos:** Use **Organization Secrets** or copy secrets to each repo

## 🔐 Authentication Types

### Service Account (Recommended)

```json
{
  "type": "service_account",
  "project_id": "...",
  "private_key": "...",
  "client_email": "..."
}
```

**Setup:** [README.md - Service Account Section](README.md#option-1-service-account-authentication-recommended)

### OAuth Refresh Token

```json
{
  "auth_type": "oauth_refresh_token",
  "credentials": "1//refresh-token",
  "oauth_client_id": "client-id",
  "oauth_client_secret": "client-secret"
}
```

**Setup:** [OAUTH_SETUP.md](OAUTH_SETUP.md)

### OAuth Access Token (Testing Only)

```yaml
auth_type: "oauth"
credentials: "ya29.access-token"
```

⚠️ Expires in ~1 hour

---

## 📝 Common Use Cases

### Log Deployments

```yaml
name: Log Deployment

on:
  release:
    types: [published]

jobs:
  log:
    uses: AlbertPuwadol/google-sheet-github-action/.github/workflows/reusable-append-row.yml@main
    with:
      sheet_name: "Deployments"
      values: '["${{ github.repository }}", "${{ github.event.release.tag_name }}", "${{ github.actor }}", "${{ github.event.release.published_at }}"]'
```

### Track CI/CD Pipeline

```yaml
name: Track Build

on: [push, pull_request]

jobs:
  track:
    uses: AlbertPuwadol/google-sheet-github-action/.github/workflows/reusable-append-row.yml@main
    with:
      sheet_name: "CI Pipeline"
      values: '["${{ github.workflow }}", "${{ github.event_name }}", "${{ github.sha }}", "${{ job.status }}"]'
```

### Monitor Issue Activity

```yaml
name: Track Issues

on:
  issues:
    types: [opened, closed]

jobs:
  track:
    uses: AlbertPuwadol/google-sheet-github-action/.github/workflows/reusable-append-row.yml@main
    with:
      sheet_name: "Issues"
      values: '["${{ github.event.issue.number }}", "${{ github.event.action }}", "${{ github.event.issue.title }}", "${{ github.actor }}"]'
```

### Log Pull Requests

```yaml
name: Track PRs

on:
  pull_request:
    types: [opened, closed, merged]

jobs:
  track:
    uses: AlbertPuwadol/google-sheet-github-action/.github/workflows/reusable-append-row.yml@main
    with:
      sheet_name: "Pull Requests"
      values: '["${{ github.event.pull_request.number }}", "${{ github.event.action }}", "${{ github.event.pull_request.title }}", "${{ github.actor }}"]'
```

---

## 🆘 Quick Troubleshooting

### Error: "credentials is required"

**Solutions:**

1. Check secret name matches exactly (case-sensitive)
2. Verify secret is set in correct location

### Error: "Spreadsheet not found"

**Solutions:**

1. Verify spreadsheet ID is correct
2. Check service account has access to sheet
3. Share sheet with service account email

### Error: "Invalid header field value"

**Solutions:**

1. Recreate the secret (copy-paste issue)
2. Ensure no extra spaces or newlines
3. For JSON, copy entire file contents

### Data length not showing

**Solutions:**

1. Check action logs for "Reading existing data" message
2. Verify sheet permissions allow reading
3. Ensure outputs are accessed correctly in workflow

---

## 📊 Key Features

### Automatic Data Tracking

The action automatically provides:

- `existing_row_count` - Rows before append
- `row_number` - Where data was added
- `total_rows` - Total rows after append
- `updated_range` - A1 notation (e.g., `Sheet1!A11:C11`)

### Precise Range Notation

Instead of generic `A:Z`, uses exact position like `A11:Z11` based on existing data.

## 📚 Full Documentation

- **[README.md](README.md)** - Complete documentation
- **[OAUTH_SETUP.md](OAUTH_SETUP.md)** - OAuth authentication setup

---

## 💡 Pro Tips

1. **Use organization secrets for teams:**

   - Centralized management
   - Fine-grained access control
   - Audit logs

2. **Use row tracking for sequential IDs:**

   ```yaml
   - id: log
     uses: .../
   - run: echo "Logged as entry #${{ steps.log.outputs.row_number }}"
   ```

3. **Monitor sheet capacity:**

   ```yaml
   - id: append
     uses: .../
   - run: |
       if [ ${{ steps.append.outputs.total_rows }} -gt 1000 ]; then
         echo "⚠️ Sheet is getting full!"
       fi
   ```

4. **Add timestamps automatically:**
   ```yaml
   values: '["${{ github.sha }}", "${{ github.event.head_commit.timestamp }}", "data"]'
   ```

---

## 🔗 Links

- [GitHub Repository](https://github.com/AlbertPuwadol/google-sheet-github-action)
- [Report Issues](https://github.com/AlbertPuwadol/google-sheet-github-action/issues)
- [Google Sheets API](https://developers.google.com/sheets/api)

---

**Need help?** Open an issue with the `question` label!
