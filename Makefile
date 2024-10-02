# Variables
BINFOLDER=bin
BINNAME=okteto-test-action

.PHONY: all build test cves clean

all: build test cves

build:
	go build -o bin/okteto-test -ldflags="-s -w" ./cmd/main.go

test:
	go test ./... -v

cves:
	@echo "Running Trivy for CVE scan..."
	@docker build -t $(BINNAME):latest .
	@docker save $(BINNAME):latest -o $(BINNAME).tar
	@docker run --rm -v $(PWD)/$(BINNAME).tar:/$(BINNAME).tar -v $(PWD)/trivy-cache:/root/.cache/ aquasec/trivy:latest image --severity HIGH,CRITICAL --input /$(BINNAME).tar
	@rm $(BINNAME).tar

clean:
	rm -rf $(BINFOLDER)
