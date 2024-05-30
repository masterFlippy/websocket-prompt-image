.PHONY: build-EventBridgeProxyFunction package-EventBridgeProxyFunction

build-EventBridgeProxyFunction:
	GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
	chmod +x bootstrap
	cp bootstrap $(ARTIFACTS_DIR)/

package-EventBridgeProxyFunction:
	cd $(ARTIFACTS_DIR) && zip -r function.zip bootstrap