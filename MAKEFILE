# ==== CONFIGURATION ====
PROJECT_ID := flatchecker
REGION := europe-north1
REPO := flatchecker-containers
IMAGE_NAME := flatchecker-scheduler
IMAGE_TAG := latest

# Construct full image name
IMAGE_URI := $(REGION)-docker.pkg.dev/$(PROJECT_ID)/$(REPO)/$(IMAGE_NAME):$(IMAGE_TAG)

# ==== COMMANDS ====

# Build the Docker image
build:
	docker build -t $(IMAGE_URI) .

# Push the image to Artifact Registry
push: build
	docker push $(IMAGE_URI)

# Clean local Docker image (optional)
clean:
	docker rmi $(IMAGE_URI)

# Full pipeline: build + push
deploy: push
	@echo "Image pushed: $(IMAGE_URI)"
