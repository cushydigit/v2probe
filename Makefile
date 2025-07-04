.PHONY: build run

build:
	@echo "building..."
	@go build -o ./bin/v2probe ./

run: build
	@echo "running..."
	@./bin/v2probe
