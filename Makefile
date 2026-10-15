.PHONY: test-api build run clean

test-api:
	./test-api.sh

build:
	go build -o pizzeria-api main.go

run: build
	./pizzeria-api

clean:
	docker compose down -v --remove-orphans || true
	rm -f pizzeria-api