newpost:
	@read -p "Enter post title: " POST_TITLE; \
	hugo new posts/"$${POST_TITLE// /-}".md
