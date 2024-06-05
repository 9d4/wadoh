GO := go
PROTOC := protoc
BIN_OUT := wadoh

build: build-js proto-gen build-bin

build-bin:
	$(GO) buildx -o $(BIN_OUT) .

run:
	@echo "Running in development"
	@echo "Ensure to run 'make watch-js' so any frontend changes within html/ will be updated"
	$(GO) run -tags dev .

proto-gen:
	$(PROTOC) \
		--go_out=wadoh-be \
		--go-grpc_out=wadoh-be \
		--proto_path=wadoh-be/proto \
		wadoh.proto

clean-js: html/static/dist
	cd html && rm -rf .parcel-cache

build-js: clean-js
	cd html && yarn build 

watch-js: clean-js
	cd html && yarn watch 
