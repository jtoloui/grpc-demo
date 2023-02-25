runServer:
	@echo "Running server..."
	cd server && npm run start:dev
.PHONY: runServer

runClient:
	@echo "Running client..."
	cd client && npm run start:dev
.PHONY: runClient

runAll: 
	@echo "Running server and client..."
	cd server && npm run start:dev &
	cd client && npm run start:dev
.PHONY: runAll