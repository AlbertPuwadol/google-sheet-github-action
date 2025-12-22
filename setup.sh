#!/bin/bash

# Setup script for Google Sheet GitHub Action

echo "🚀 Setting up Google Sheet GitHub Action..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.23 or later."
    echo "   Visit: https://golang.org/doc/install"
    exit 1
fi

echo "✅ Go is installed: $(go version)"

# Download dependencies
echo "📦 Downloading dependencies..."
go mod download
if [ $? -ne 0 ]; then
    echo "❌ Failed to download dependencies"
    exit 1
fi

# Tidy up go.mod and create go.sum
echo "🧹 Running go mod tidy..."
go mod tidy
if [ $? -ne 0 ]; then
    echo "❌ Failed to tidy modules"
    exit 1
fi

# Build the binary
echo "🔨 Building the binary..."
go build -o append-sheet .
if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Set up Google Cloud Service Account (see README.md)"
echo "2. Share your Google Sheet with the service account email"
echo "3. Add secrets to your GitHub repository:"
echo "   - GOOGLE_SERVICE_ACCOUNT_JSON"
echo "   - SPREADSHEET_ID"
echo "4. Push this repository to GitHub"
echo "5. Use the action in your workflows!"
echo ""
echo "📖 For detailed instructions, see README.md"

