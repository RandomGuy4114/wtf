BINARY := bin/wtf

.PHONY: build run vet test clean install

build:
	go build -o $(BINARY) .

run: build
	./$(BINARY)

install:
	go install .
	git config core.hooksPath hooks

vet:
	go vet ./...

test:
	go test ./...

clean:
	rm -f $(BINARY)
