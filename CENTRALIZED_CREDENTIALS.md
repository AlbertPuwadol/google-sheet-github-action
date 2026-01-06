# Centralized Credentials Setup

This guide explains how to set up credentials once and use them across multiple repositories.

## 🎯 Two Approaches

### Approach 1: Reusable Workflow (Recommended) ⭐

Store credentials in the action repository and call it as a reusable workflow from other repositories.

#### Setup in Action Repository

- **Add credentials to this repository:**
  - Go to this repo → Settings → Secrets and variables → Actions
  - Add environment variables (not secrets, for reusable workflows):
    - `GOOGLE_CREDENTIALS` - Your service account JSON
    - `SPREADSHEET_ID` - Your spreadsheet ID

#### Usage from Other Repositories

In any repository, create `.github/workflows/use-sheets-action.yml`:

```yaml
name: Update Google Sheet

on:
  push:
    branches: [main]

jobs:
  log-deployment:
    uses: AlbertPuwadol/google-sheet-github-action/.github/workflows/reusable-append-row.yml@main
    with:
      sheet_name: "Deployments"
      values: '["${{ github.sha }}", "${{ github.actor }}", "${{ github.event.head_commit.message }}"]'
      use_action_repo_credentials: true
```

**Benefits:**

- ✅ Credentials stored once in action repository
- ✅ No need to add secrets to every repository
- ✅ Centralized credential management
- ✅ Easy to rotate credentials (update in one place)

**Limitations:**

- The action repository must be accessible by caller repositories
- Requires checkout of action repository during workflow run

---

### Approach 2: Organization Secrets (Best for Teams)

If you're using GitHub Organizations, you can set organization-level secrets.

#### Setup

1. **Go to Organization Settings:**

   - Navigate to your GitHub Organization
   - Go to Settings → Secrets and variables → Actions
   - Click "New organization secret"

2. **Add secrets:**

   - Name: `ORG_GOOGLE_CREDENTIALS`
   - Value: Your service account JSON
   - Repository access: Select which repositories can use this secret

3. **Add spreadsheet ID:**
   - Name: `ORG_SPREADSHEET_ID`
   - Value: Your spreadsheet ID

#### Usage from Any Repository in Organization

```yaml
name: Update Sheet

on:
  push:
    branches: [main]

jobs:
  update-sheet:
    runs-on: ubuntu-latest
    steps:
      - uses: AlbertPuwadol/google-sheet-github-action@main
        with:
          spreadsheet_id: ${{ secrets.ORG_SPREADSHEET_ID }}
          credentials: ${{ secrets.ORG_GOOGLE_CREDENTIALS }}
          sheet_name: "Activity"
          values: '["${{ github.repository }}", "${{ github.actor }}"]'
```

**Benefits:**

- ✅ True centralized credentials
- ✅ Fine-grained access control per repository
- ✅ Audit logs
- ✅ No additional workflow files needed

**Limitations:**

- Only available for GitHub Organizations (not personal accounts)
- Requires organization admin access to set up

---

## 🔄 Comparison Table

| Approach             | Setup Complexity | Credential Management | Access Control | Best For                   |
| -------------------- | ---------------- | --------------------- | -------------- | -------------------------- |
| Reusable Workflow    | Medium           | Centralized           | Repository     | Personal accounts, teams   |
| Organization Secrets | Low              | Centralized           | Fine-grained   | Organizations, enterprises |
| Per-repo Secrets     | Low              | Distributed           | Per-repo       | Independent projects       |

## 📋 Step-by-Step: Reusable Workflow Setup

### In Calling Repository (Your Other Projects)

**Step 1:** Create workflow file

Create `.github/workflows/update-sheet.yml`:

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
      values: '["${{ github.repository }}", "${{ github.sha }}", "${{ github.actor }}", "${{ github.event.head_commit.message }}"]'
      use_action_repo_credentials: true
```

**Step 2:** Push and test

```bash
git add .github/workflows/update-sheet.yml
git commit -m "Add Google Sheets logging"
git push
```

Check the Actions tab to see it run!
