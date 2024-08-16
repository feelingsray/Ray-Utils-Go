tidy:
	go mod tidy

install-gofumpt:
	go install mvdan.cc/gofumpt@latest
fmt:
	gofumpt -l -w .