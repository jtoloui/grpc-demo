import {
	Metadata,
	sendUnaryData,
	ServerUnaryCall,
	status,
} from "@grpc/grpc-js";
import {
	GetMovieByTitleRequest,
	GetMovieByTitleResponse,
	Movie,
} from "@jtoloui/proto-store";
import { logger as log } from "../../logger";

export const getMovieByTitle = (
	call: ServerUnaryCall<GetMovieByTitleRequest, GetMovieByTitleResponse>,
	callback: sendUnaryData<GetMovieByTitleResponse>
) => {
	const logger = log("getMovieByTitle");
	logger.info(`tracer call: ${call.metadata.get("X-Tracer-Id")[0]}`);

	if (call.request.title.length === 0) {
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

	const movieObj = Movie.create({
		title: call.request.title,
		year: 1972,
	});

	const response = GetMovieByTitleResponse.create({
		movie: movieObj,
	});
	callback(null, response);
};
