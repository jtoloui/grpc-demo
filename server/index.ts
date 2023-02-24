import {
	Metadata,
	sendUnaryData,
	Server,
	ServerCredentials,
	ServerUnaryCall,
	status,
} from "@grpc/grpc-js";

import {
	GetMovieByTitleRequest,
	GetMovieByTitleResponse,
	IMoviesService,
	Movie,
	moviesServiceDefinition,
} from "@jtoloui/proto-store";

// create unary call

const getMovieByTitle = (
	call: ServerUnaryCall<GetMovieByTitleRequest, GetMovieByTitleResponse>,
	callback: sendUnaryData<GetMovieByTitleResponse>
) => {
	console.log("hello");

	console.log(call.request);

	if (call.request.title.length === 0) {
		const metadata = new Metadata();
		metadata.add("error", "invalid title");
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

const server = new Server();

const service: IMoviesService = {
	getMovieByTitle,
};

server.addService(moviesServiceDefinition, service);
server.bindAsync(
	"0.0.0.0:50051",
	ServerCredentials.createInsecure(),
	(err: Error | null, port: number) => {
		if (err) {
			console.error(`Server error: ${err.message}`);
		} else {
			console.log(`Server bound on port: ${port}`);
			server.start();
		}
	}
);
