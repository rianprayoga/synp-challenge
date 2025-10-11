run-inventories:
	@cd '$(CURDIR)/inventories-service' && go run ./cmd/api/
	

gen: 
	@protoc --go_out=. --go_opt=paths=source_relative \
	 --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	 inventories-rpc/inventories.proto