#!/bin/bash
set -e

# Cloudflare Pages build script for Hugo
# Uses production URL for main branch, CF_PAGES_URL for preview deployments

if [ "$CF_PAGES_BRANCH" = "main" ]; then
  # Production deployment - use the configured domain
  BASE_URL="https://ark31.info"
elif [ -n "$CF_PAGES_URL" ]; then
  # Preview deployment - use the provided URL
  BASE_URL="$CF_PAGES_URL"
else
  # Local build fallback
  BASE_URL="https://ark31.info"
fi

# Ensure baseURL ends with /
if [[ ! "$BASE_URL" =~ /$ ]]; then
  BASE_URL="${BASE_URL}/"
fi

echo "Building Hugo site with baseURL: $BASE_URL"

# Build Hugo with the appropriate baseURL
hugo --baseURL "$BASE_URL" --minify
