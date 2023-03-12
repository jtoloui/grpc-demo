import { sendUnaryData, ServerUnaryCall, status } from "@grpc/grpc-js";
import { GetMoviesRequest, GetMoviesResponse } from "@jtoloui/proto-store";
import { logger as log } from "../../logger";
import { Movie } from "../../models";

export const getMovies = (
	call: ServerUnaryCall<GetMoviesRequest, GetMoviesResponse>,
	callback: sendUnaryData<GetMoviesResponse>
) => {
	const logger = log("getMovies");
	logger.info(`tracer call: ${call.metadata.get("X-Tracer-Id")[0]}`);

	const page = call.request.page || 1;
	const perPage = call.request.perPage || 10;

	const movies = Movie.find({})
		.lean()
		.skip((page - 1) * perPage)
		.limit(perPage)
		.exec();

	const totalRes = Movie.countDocuments().exec();

	Promise.all([movies, totalRes])
		.then((res) => {
			const [movies, total] = res;
			const response = GetMoviesResponse.create({
				movies: movies.map((movie) => ({
					title: movie.title,
					year: movie.year,
					director: movie.director,
					id: movie._id,
				})),
				total,
			});
			callback(null, response);
		})
		.catch((err) => {
			logger.error(err.message);
			callback(
				{
					code: status.INTERNAL,
					details: "internal server error",
				},
				null
			);
		});
};
