APP_NAME=ngoclam-zmp-be
MAIN_FILE=cmd/main.go
APP_FILE=bin/ngoclam-zmp-be



# Lệnh generate swagger docs
dev: 
	go run -mod=mod -tags=dev $(MAIN_FILE)

# Build binary
build: 
	go build -o bin/$(APP_NAME) $(MAIN_FILE)

# Run trực tiếp
run: 
	./$(APP_FILE)

diff:
	@read -p "Enter migration name: " name; \
	atlas migrate diff $$name --env gorm

migrate:
	atlas migrate apply --env gorm

status:
	atlas migrate status --env gorm

# Clean build + docs
clean:
	rm -rf bin/ docs/docs.go docs/swagger.json docs/swagger.yaml

.PHONY: swag build run clean
