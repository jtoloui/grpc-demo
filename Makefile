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

run-all:
	@echo "Running all servers..."
	cd javascript/server-grpc && npm run start:dev &
	cd javascript/server-express && npm run start:dev

.PHONY: run-% js-server js-client run-all go-server go-client