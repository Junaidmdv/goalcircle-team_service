
VERSION := latest 
updateproto:
ifndef VERSION
	$(error VERSION is required. usage: make updateproto VERSION=v1.0.0)
endif
	go get github.com/Junaidmdv/goalcircle-protos@$(VERSION)
	go mod tidy
	@echo "✓ Updated goalcircle-protos to $(VERSION)"