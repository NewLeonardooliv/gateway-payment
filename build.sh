#!/bin/bash

set -e

OUTPUT_DIR="./bin"
OUTPUT_NAME="app"
CMD_PATH="./cmd/app"

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
COMMIT=$(git rev-parse --short HEAD)
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS="-s -w -X 'main.version=$VERSION' -X 'main.commit=$COMMIT' -X 'main.date=$BUILD_DATE'"

mkdir -p "$OUTPUT_DIR"

echo "🔨 Compilando aplicação..."
go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/$OUTPUT_NAME" "$CMD_PATH"

echo "✅ Build finalizado com sucesso:"
echo "   ➤ Versão: $VERSION"
echo "   ➤ Commit: $COMMIT"
echo "   ➤ Data:   $BUILD_DATE"
echo "   ➤ Binário: $OUTPUT_DIR/$OUTPUT_NAME"
