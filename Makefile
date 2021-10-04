.PHONY: install
install: bin
		cp ./ize $(GOPATH)/bin/ize

.PHONY: bin
bin: 
	go build -o ./ize ./cmd 