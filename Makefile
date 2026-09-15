BINARY := bin/wtf

.PHONY: build run vet clean install

build:
	go build -o $(BINARY) .

run: build
	./$(BINARY)

install:
	go install .
	git config core.hooksPath hooks

vet:
	go vet ./...

clean:
	rm -f $(BINARY)
