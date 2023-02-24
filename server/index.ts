import {
	Metadata,
	sendUnaryData,
	Server,
	ServerCredentials,
	ServerUnaryCall,
} from "@grpc/grpc-js";
import { Status } from "@grpc/grpc-js/build/src/constants";

// create unary call

const getMovieByTitle = (
	call: ServerUnaryCall<GetMovieByTitleRequest, GetMovieByTitleResponse>,
	callback: sendUnaryData<GetMovieByTitleResponse>
) => {
	console.log(call.request);
	const movie = {
		title: "The Godfather",
		year: 1972,
	};

	if (call.request.title !== "dd") {
		// add metadata and error

		const metadata = new Metadata();
		metadata.add("error", "movie not found");
		callback(
			{
				code: Status.NOT_FOUND,
				details: "movie not found",
				metadata,
			},
			null,
			metadata
		);
	}

	const movieObj = Movie.create(movie);

	const response = GetMovieByTitleResponse.create({
		movie: movieObj,
	});
	callback(null, response);
};

import {
	GetMovieByTitleRequest,
	GetMovieByTitleResponse,
	IMoviesService,
	Movie,
	moviesServiceDefinition,
} from "@jtoloui/proto-store";

const server = new Server();

const svr: IMoviesService = {
	getMovieByTitle,
};

server.addService(moviesServiceDefinition, svr);
server.bindAsync(
	"0.0.0.0:50051",
	ServerCredentials.createInsecure(),
	(err, port) => {
		if (err) {
			throw err;
		}
		server.start();
		console.log(`Server running on port ${port}`);
	}
);
