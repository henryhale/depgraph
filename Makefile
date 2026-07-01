tidy:
	go mod tidy

run:
	go run depgraph.go

build:
	chmod +x ./scripts/build.sh
	./scripts/build.sh

lint:
	chmod +x ./scripts/lint.sh
	./scripts/lint.sh

test:
	go test ./...

clean:
	rm -f depgraph

.PHONY: tidy run build lint clean test
