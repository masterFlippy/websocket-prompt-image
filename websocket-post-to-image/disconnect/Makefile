.PHONY: build-DisconnectFunction package-DisconnectFunction

build-DisconnectFunction:
	GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
	chmod +x bootstrap
	cp bootstrap $(ARTIFACTS_DIR)/

package-DisconnectFunction:
	cd $(ARTIFACTS_DIR) && zip -r function.zip bootstrap