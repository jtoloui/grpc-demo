import {
	Metadata,
	sendUnaryData,
	ServerUnaryCall,
	status,
} from "@grpc/grpc-js";
import {
	GetMovieByIdRequest,
	GetMovieByIdResponse,
} from "@jtoloui/proto-store";
import { logger as log } from "../../logger";
import { Movie } from "../../models";

export const getMovieById = (
	call: ServerUnaryCall<GetMovieByIdRequest, GetMovieByIdResponse>,
	callback: sendUnaryData<GetMovieByIdResponse>
) => {
	const logger = log("getMovieById");
	logger.info(`tracer call: ${call.metadata.get("X-Tracer-Id")[0]}`);

	if (call.request.id.length === 0) {
		const metadata = new Metadata();
		metadata.add("error", "invalid title");
		logger.error("invalid title");
		callback(
			{
				code: status.INVALID_ARGUMENT,
				details: "invalid title",
				metadata,
			},
			null,
			metadata
		);
	}

	const findMovie = Movie.findById(call.request.id);
	findMovie
		.then((movie) => {
			if (movie) {
				const response = GetMovieByIdResponse.create({
					movie: {
						title: movie.title,
						year: movie.year,
						director: movie.director,
					},
				});
				callback(null, response);
			}
		})
		.catch((err) => {
			logger.error(err.message);
			callback(
				{
					code: status.NOT_FOUND,
					details: `record not found for id: ${call.request.id}`,
				},
				null
			);
		});
};
