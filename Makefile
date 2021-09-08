# Install a macOS universal binary.
darwin:
	mkdir -p .build
	GOOS=darwin GOARCH=amd64 go build -o .build/godoc-readme-gen.amd64
	GOOS=darwin GOARCH=arm64 go build -o .build/godoc-readme-gen.arm64
	lipo \
		.build/godoc-readme-gen.amd64 \
		.build/godoc-readme-gen.arm64 \
		-create \
		-output godoc-readme-gen
	mv godoc-readme-gen $(GOPATH)/bin
	rm -rf .build/
