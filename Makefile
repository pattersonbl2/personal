newpost:
	@read -p "Enter post title: " POST_TITLE; \
	POST_SLUG=$$(echo "$$POST_TITLE" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9_-]/-/g; s/-\{2,\}/-/g; s/^-//; s/-$$//'); \
	hugo new posts/"$$POST_SLUG".md
