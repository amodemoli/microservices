#!/bin/bash

BASE_URL="https://raw.githubusercontent.com/amodemoli/response-codes/main"
RESPONSE_FILE="response.codes.go"
STATUS_FILE="status.codes.go"

RESPONSE_DIR="../internal/helpers/codes/response"
STATUS_DIR="../internal/helpers/codes/status"

# create directory's if not exists
mkdir -p "${RESPONSE_DIR}"
mkdir -p "${STATUS_DIR}"

echo "Downloading ${RESPONSE_FILE}..."
curl -s -o "$RESPONSE_DIR/$RESPONSE_FILE" "$BASE_URL/$RESPONSE_FILE"
if [ $? -eq 0 ]; then
    echo "> $RESPONSE_FILE saved to $RESPONSE_DIR/"
else
    echo "> Failed to download $RESPONSE_FILE"
fi

echo "Downloading ${STATUS_FILE}..."
curl -s -o "$STATUS_DIR/$STATUS_FILE" "$BASE_URL/$STATUS_FILE"
if [ $? -eq 0 ]; then
    echo "> $STATUS_FILE saved to $STATUS_DIR/"
else
    echo "> Failed to download $STATUS_FILE"
fi

echo "Done."