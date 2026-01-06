# Quick Reference Guide

## 🚀 3 Ways to Use This Action

### 1️⃣ Reusable Workflow (Centralized Credentials) ⭐ EASIEST

**Store credentials once in action repo, use from anywhere!**

```yaml
jobs:
  log:
    uses: AlbertPuwadol/google-sheet-github-action/.github/workflows/reusable-append-row.yml@main
    with:
      sheet_name: "Logs"
      values: '["data1", "data2"]'
```

**No secrets needed in calling repository!**

📖 [Full Guide](CENTRALIZED_CREDENTIALS.md)

---

### 2️⃣ Organization Secrets (Best for Teams)

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

---

### 3️⃣ Repository Secrets (Standard)

**Add secrets to each repository individually**

```yaml
- uses: AlbertPuwadol/google-sheet-github-action@main
  with:
    spreadsheet_id: ${{ secrets.SPREADSHEET_ID }}
    credentials: ${{ secrets.GOOGLE_SERVICE_ACCOUNT_JSON }}
    sheet_name: "Logs"
    values: '["data"]'
```

**Setup:** Settings → Secrets and variables → Actions

---

## 📊 Comparison

| Method                | Setup Once | Use Everywhere | GitHub Org Required | Best For        |
| --------------------- | ---------- | -------------- | ------------------- | --------------- |
| Reusable Workflow     | ✅         | ✅             | ❌                  | Multiple repos  |
| Organization Secrets  | ✅         | ✅             | ✅                  | Teams           |
| Repository Secrets    | ❌         | ❌             | ❌                  | Single repo     |
| Environment Variables | ❌         | ❌             | ❌                  | Per-repo config |

## 🎯 Recommendations

- **👥 Multiple repos, no GitHub Org:** Use **Reusable Workflow** (#1)
- **🏢 GitHub Organization:** Use **Organization Secrets** (#2)
- **📦 Single repository:** Use **Repository Secrets** (#3)

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
3. For reusable workflow, ensure `use_action_repo_credentials: true`

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

### Reusable workflow not working

**Solutions:**

1. Check action repository is accessible
2. Verify config file exists in action repo
3. Ensure `@main` branch reference is correct

---

## 📚 Full Documentation

- **[README.md](README.md)** - Complete documentation
- **[CENTRALIZED_CREDENTIALS.md](CENTRALIZED_CREDENTIALS.md)** - Reusable workflows & organization secrets
- **[OAUTH_SETUP.md](OAUTH_SETUP.md)** - OAuth authentication setup

---

## 💡 Pro Tips

1. **Use reusable workflow for multiple repos:**

   - No secrets needed in each repo
   - Update credentials in one place
   - Perfect for microservices

2. **Use organization secrets for teams:**

   - Centralized management
   - Fine-grained access control
   - Audit logs

3. **Separate sheets per environment:**

   - Dev: `Dev-Logs`
   - Staging: `Staging-Logs`
   - Prod: `Production-Logs`

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
