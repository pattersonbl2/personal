# Makefile for auto-creating a new Hugo post

newpost:
	@read -p "Enter post title: " POST_TITLE; \
	POST_SLUG=$$(echo $$POST_TITLE | tr '[:upper:]' '[:lower:]' | tr ' ' '-'); \
	POST_DATE=$$(date +"%Y-%m-%dT%H:%M:%S%z"); \
	POST_PATH=content/posts/$$POST_SLUG.md; \
	echo "---" > $$POST_PATH; \
	echo "title: \"$$POST_TITLE\"" >> $$POST_PATH; \
	echo "date: $$POST_DATE" >> $$POST_PATH; \
	echo "draft: true" >> $$POST_PATH; \
	echo "---" >> $$POST_PATH; \
	echo "Created new post at $$POST_PATH"

# Usage:
# make newpost POST_TITLE="Your Post Title"