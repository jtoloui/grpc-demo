import { sendUnaryData, ServerUnaryCall, status } from "@grpc/grpc-js";
import {
	CreateMovieRequest,
	CreateMovieResponse,
	Movie as MovieGrpc,
} from "@jtoloui/proto-store";

import { logger as log } from "../../logger";
import { Movie } from "../../models";

export const createMovie = async (
	call: ServerUnaryCall<CreateMovieRequest, CreateMovieResponse>,
	callback: sendUnaryData<CreateMovieResponse>
) => {
	const logger = log("createMovie");
	logger.info(`tracer call: ${call.metadata.get("X-Tracer-Id")[0]}`);
	const movieReq = call.request.movie;
	// validate request

	if (movieReq === undefined) {
		logger.error("invalid movie");
		callback(
			{
				code: status.INVALID_ARGUMENT,
				details: "invalid movie",
			},
			null
		);
	}

	if (movieReq?.title === undefined) {
		logger.error("invalid title");
		callback(
			{
				code: status.INVALID_ARGUMENT,
				details: "invalid title",
			},
			null
		);
	}

	if (movieReq?.year === undefined) {
		logger.error("invalid year");
		callback(
			{
				code: status.INVALID_ARGUMENT,
				details: "invalid year",
			},
			null
		);
	}

	if (movieReq?.director === undefined) {
		logger.error("invalid director");
		callback(
			{
				code: status.INVALID_ARGUMENT,
				details: "invalid director",
			},
			null
		);
	}

	const title = movieReq?.title as string;
	const director = movieReq?.director as string;
	const year = movieReq?.year as number;

	// create movie
	const movie = Movie.build({
		title,
		year,
		director,
	});

	await movie.save();

	const movieResponse = MovieGrpc.create({
		title: movie.title,
		year: movie.year,
		director: movie.director,
	});

	const response = CreateMovieResponse.create({
		movie: movieResponse,
		id: movie.id,
	});

	callback(null, response);
};
