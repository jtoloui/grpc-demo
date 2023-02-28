# Define the generic rule for running servers
run-%:
	@echo "Running server $*..."
	cd $(ROOT)/server-$* && $(EXECUTION)

# Define the targets that depend on the generic rule
js-server: 
	$(MAKE) run-grpc ROOT=javascript EXECUTION="npm run start:dev"
js-client:
	$(MAKE) run-express ROOT=javascript EXECUTION="npm run start:dev"

# Define the target all go servers

go-server:
	$(MAKE) run-grpc ROOT=go EXECUTION="go run main.go"
go-client:
	$(MAKE) run-gin ROOT=go EXECUTION="go run main.go"

go-compile-%:
	@echo "Building server $*..."
	cd go/server-$* && go build -o bin/server-$* main.go

go-build-grpc: go-compile-grpc
go-build-gin: go-compile-gin

go-run-grpc: 
	@echo "Running server grpc..."
	cd go/server-grpc && ./bin/server-grpc

go-run-gin:
	@echo "Running server gin..."
	cd go/server-gin && ./bin/server-gin

run-all:
	@echo "Running all servers..."
	cd javascript/server-grpc && npm run start:dev &
	cd javascript/server-express && npm run start:dev

.PHONY: run-% js-server js-client run-all go-server go-client go-build-grpc go-build-gin go-compile-%