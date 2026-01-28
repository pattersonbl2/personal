#!/bin/bash

# Cloudflare Pages build script for Hugo
# Uses CF_PAGES_URL for preview deployments, falls back to production URL

if [ -n "$CF_PAGES_URL" ]; then
  # Preview deployment - use the provided URL
  BASE_URL="$CF_PAGES_URL"
else
  # Production deployment - use the configured domain
  BASE_URL="https://ark31.info"
fi

# Ensure baseURL ends with /
if [[ ! "$BASE_URL" =~ /$ ]]; then
  BASE_URL="${BASE_URL}/"
fi

echo "Building Hugo site with baseURL: $BASE_URL"

# Build Hugo with the appropriate baseURL
hugo --baseURL "$BASE_URL" --minify
