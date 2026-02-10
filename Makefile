newpost:
	@read -p "Enter post title: " POST_TITLE; \
	POST_SLUG=$$(echo "$$POST_TITLE" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9_-]/-/g; s/-\{2,\}/-/g; s/^-//; s/-$$//'); \
	hugo new posts/"$$POST_SLUG".md

# Generate resume PDF. Requires: pandoc + BasicTeX (brew install pandoc basictex)
resume-pdf:
	pandoc content/resume.md -o backend/resume.pdf --metadata title="Resume" --pdf-engine=pdflatex
