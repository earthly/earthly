.phony: update
update: fmt lic imports

.PHONY: lic
lic:
	go run github.com/elastic/go-licenser@latest

.PHONY: fmt
fmt:
	go run mvdan.cc/gofumpt@latest -w -l ./

.PHONY: imports
imports:
	go run golang.org/x/tools/cmd/goimports@latest -l -local github.com/elastic/go-sysinfo ./
