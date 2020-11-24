# Copyright (c) 2020 Kristian Rumberg (kristianrumberg@gmail.com)
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

.PHONY: all
all: podbook

.PHONY: check-md
check-md:
	@find *.md | xargs misspell

.PHONY: fmt-md
fmt-md: check-md
	@find *.md | xargs markdownfmt -w

fmt-go: build
	@gofumpt -w .
	@gofmt -s -w .

golangci-lint: build fmt-go
	@golangci-lint run --enable-all --disable=exhaustivestruct,godot,goerr113,gomnd,nlreturn,wrapcheck,wsl

build:
	@mkdir -p out
	@go mod tidy
	@go build -o out/ ./...

podbook: build golangci-lint fmt-md

.PHONY: clean
clean:
	rm -rf out
	rm -rf books
	rm -rf .podbook
