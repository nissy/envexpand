VERSION := 1.0.0
NAME := envexpand
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

ifeq ($(GOOS),windows)
dist/$(NAME)-$(VERSION)_$(GOOS)_$(GOARCH).zip:
	go build -ldflags="-X main.version=$(VERSION)" -o dist/$(NAME)-$(VERSION)/$(NAME).exe cmd/main.go
	zip -j dist/$(NAME)-$(VERSION)_$(GOOS)_$(GOARCH).zip dist/$(NAME)-$(VERSION)/$(NAME).exe
else
dist/$(NAME)-$(VERSION)_$(GOOS)_$(GOARCH).tar.gz:
	go build -ldflags="-X main.version=$(VERSION)" -o dist/$(NAME)-$(VERSION)/$(NAME) cmd/$(NAME)/main.go
	tar cfz dist/$(NAME)-$(VERSION)_$(GOOS)_$(GOARCH).tar.gz -C dist/$(NAME)-$(VERSION) $(NAME)
endif

.PHONY: clean
clean:
	rm -rf dist