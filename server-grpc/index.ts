import { Server, ServerCredentials } from "@grpc/grpc-js";

import { IMoviesService, moviesServiceDefinition } from "@jtoloui/proto-store";

import { logger as log } from "./logger/logger";
import { getMovieByTitle } from "./grpc/services";

const server = new Server();

const service: IMoviesService = {
	getMovieByTitle,
};

server.addService(moviesServiceDefinition, service);
server.bindAsync(
	"localhost:50051",
	ServerCredentials.createInsecure(),
	(err: Error | null, port: number) => {
		const logger = log("server");
		if (err) {
			logger.error(`startup error: ${err.message}`);
		} else {
			logger.info(`started on 0.0.0.0:${port}`);
			server.start();
		}
	}
);
