import { Server, ServerCredentials } from "@grpc/grpc-js";
import dotenv from "dotenv";
import { IMoviesService, moviesServiceDefinition } from "@jtoloui/proto-store";

import { logger as log } from "./logger/logger";
import { getMovieById, createMovie } from "./grpc/services";
import { connectDB } from "./db";

dotenv.config();

const server = new Server();

const service: IMoviesService = {
	getMovieById,
	createMovie,
};

connectDB();

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
