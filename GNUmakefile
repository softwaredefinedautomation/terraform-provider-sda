HOSTNAME=registry.terraform.io
NAMESPACE=sda
NAME=sda
#VERSION=0.1.0
VERSION=$(shell git describe --tags --always --match 'v*')
BINARY=terraform-provider-${NAME}_${VERSION}
PLATFORMS=linux_amd64 linux_arm64 darwin_amd64 darwin_arm64 windows_amd64

default: build-all

deps:
	go mod tidy

docs: deps
	tfplugindocs generate

# Build for all platforms
.PHONY: build-all
build-all: clean deps
	@for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d'_' -f1); \
		ARCH=$$(echo $$platform | cut -d'_' -f2); \
		OUTPUT="dist/$(BINARY)_$${OS}_$${ARCH}"; \
		if [ "$$OS" = "windows" ]; then OUTPUT="$$OUTPUT.exe"; fi; \
		echo "Building for $${OS}/$${ARCH}..."; \
		GOOS=$${OS} GOARCH=$${ARCH} go build -o $${OUTPUT} \
			-ldflags="-s -w -X main.version=$(VERSION)" .; \
	done

# Create archives
.PHONY: package
package: build-all
	@echo "Creating archives..."
	@cd dist && for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d'_' -f1); \
		ARCH=$$(echo $$platform | cut -d'_' -f2); \
		ARCHIVE="$(BINARY)_$${OS}_$${ARCH}.zip"; \
		if [ "$$OS" = "windows" ]; then \
			zip $$ARCHIVE $(BINARY)_$${OS}_$${ARCH}.exe; \
		else \
			zip $$ARCHIVE $(BINARY)_$${OS}_$${ARCH}; \
		fi; \
		echo "Created $${ARCHIVE}"; \
	done

# Generate checksums
.PHONY: checksums
checksums: package
	@echo "Generating checksums..."
	@cd dist && shasum -a 256 *.zip > $(BINARY)_SHA256SUMS

# Sign checksums
.PHONY: sign
sign: checksums
	@echo "Signing checksums..."
	@cd dist && gpg --detach-sign --armor $(BINARY)_SHA256SUMS

# Complete release build
.PHONY: release
release: clean build-all package checksums sign
	@echo "Release artifacts ready in dist/"

.PHONY: clean
clean:
	@rm -rf dist/
	@mkdir -p dist/
	
