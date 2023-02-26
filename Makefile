# Define the generic rule for running servers
run-%:
	@echo "Running server $*..."
	cd server-$* && npm run start:dev

# Define the targets that depend on the generic rule
grpc: run-grpc
express: run-express

run-all:
	@echo "Running all servers..."
	cd server-grpc && npm run start:dev &
	cd server-express && npm run start:dev

.PHONY: run-% grpc express run-all