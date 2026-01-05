#!/bin/bash

# Script to validate Google Service Account JSON credentials

echo "🔍 Google Service Account Credentials Validator"
echo "================================================"
echo ""

if [ $# -eq 0 ]; then
    echo "Usage: $0 <path-to-credentials.json>"
    echo ""
    echo "Example:"
    echo "  $0 service-account-key.json"
    exit 1
fi

CREDS_FILE="$1"

# Check if file exists
if [ ! -f "$CREDS_FILE" ]; then
    echo "❌ Error: File '$CREDS_FILE' not found"
    exit 1
fi

echo "📁 File: $CREDS_FILE"
echo ""

# Check if it's valid JSON
if ! jq empty "$CREDS_FILE" 2>/dev/null; then
    echo "❌ Error: Invalid JSON format"
    echo ""
    echo "The file is not valid JSON. Please check:"
    echo "  - File is complete and not truncated"
    echo "  - No extra characters or formatting issues"
    exit 1
fi

echo "✅ Valid JSON format"

# Check required fields for Google Service Account
REQUIRED_FIELDS=("type" "project_id" "private_key_id" "private_key" "client_email" "client_id")
MISSING_FIELDS=()

for field in "${REQUIRED_FIELDS[@]}"; do
    if ! jq -e ".$field" "$CREDS_FILE" > /dev/null 2>&1; then
        MISSING_FIELDS+=("$field")
    fi
done

if [ ${#MISSING_FIELDS[@]} -ne 0 ]; then
    echo "❌ Error: Missing required fields:"
    for field in "${MISSING_FIELDS[@]}"; do
        echo "  - $field"
    done
    exit 1
fi

echo "✅ All required fields present"

# Display service account info
TYPE=$(jq -r '.type' "$CREDS_FILE")
PROJECT_ID=$(jq -r '.project_id' "$CREDS_FILE")
CLIENT_EMAIL=$(jq -r '.client_email' "$CREDS_FILE")

echo ""
echo "📋 Service Account Details:"
echo "  Type: $TYPE"
echo "  Project ID: $PROJECT_ID"
echo "  Email: $CLIENT_EMAIL"

# Check for common issues
if [ "$TYPE" != "service_account" ]; then
    echo ""
    echo "⚠️  Warning: Type is '$TYPE' but should be 'service_account'"
fi

# Check for whitespace issues
FILE_SIZE=$(wc -c < "$CREDS_FILE" | tr -d ' ')
CONTENT=$(cat "$CREDS_FILE")
CONTENT_SIZE=${#CONTENT}

if [ "$FILE_SIZE" -ne "$CONTENT_SIZE" ]; then
    echo ""
    echo "⚠️  Warning: File may contain unusual whitespace or control characters"
fi

echo ""
echo "✅ Credentials appear valid!"
echo ""
echo "📝 Next steps:"
echo "  1. Copy the entire file contents:"
echo "     cat $CREDS_FILE | pbcopy    (on macOS)"
echo "     cat $CREDS_FILE | xclip -selection clipboard    (on Linux)"
echo ""
echo "  2. In GitHub, go to: Settings > Secrets and variables > Actions"
echo "  3. Click 'New repository secret'"
echo "  4. Name: GOOGLE_SERVICE_ACCOUNT_JSON"
echo "  5. Value: Paste the copied content (Cmd+V or Ctrl+V)"
echo "  6. Click 'Add secret'"
echo ""
echo "  Remember to share your Google Sheet with: $CLIENT_EMAIL"

