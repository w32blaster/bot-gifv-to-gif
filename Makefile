PKG := github.com/w32blaster/bot-gifv-to-gif
PKG_LIST=$(shell go list ${PKG}/... | grep -v /vendor/)
FILES := $$(find . -name '*.go' | grep -vE 'vendor')

test:
	@go test -v -race ${PKG_LIST}

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags "-s -w" -o bot-gifv-to-gif  github.com/w32blaster/bot-gifv-to-gif

regenerate:
	@echo "Regenerate the code using bot-scaffolding"
	@bot-scaffolding -o generate

.POHNY: test build generateTest up