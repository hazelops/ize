GIT_COMMIT=$$(git rev-parse --short HEAD)
GIT_IMPORT="github.com/hazelops/ize/internal/version"
GIT_DIRTY=$$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GOLDFLAGS="-s -w -X $(GIT_IMPORT).GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)"
CGO_ENABLED?=0

.PHONY: install
install: bin
		mv ./ize $(GOPATH)/bin/ize

.PHONY: bin
bin: 
	CGO_ENABLED=$(CGO_ENABLED) go build -ldflags $(GOLDFLAGS) -o ./ize ./cmd 