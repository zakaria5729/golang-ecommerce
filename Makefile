build:
	@echo "Building easy_com_api..."
	go build -o bin/easy_com_api ./cmd

run: build
	@echo "Starting easy_com_api..."
	./bin/easy_com_api

clean:
	rm -rf bin/

.PHONY: build run clean
