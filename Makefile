tidy:
	go mod tidy

install-fmt:
	go install github.com/fsgo/go_fmt/cmd/gorgeous@latest
fmt:
	gorgeous ./...