.PHONY: setup run build clean

# Create a local .env file from .env.example
setup:
	cp -n .env.example .env || echo ".env already exists"

# Run the bot with Go
run:
	go run ./cmd/main.go

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
