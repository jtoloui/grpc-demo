import express, { Express, Request, Response } from "express";
import dotenv from "dotenv";
import { ChannelCredentials } from "@grpc/grpc-js";
import {
	MoviesServiceClient,
	IMoviesServiceClient,
	GetMovieByTitleRequest,
	GetMovieByTitleResponse,
} from "@jtoloui/proto-store";
import axios from "axios";

dotenv.config();

const client = new MoviesServiceClient(
	"0.0.0.0:50051",
	ChannelCredentials.createInsecure(),
	{},
	{}
);

const deadline = new Date();
deadline.setSeconds(deadline.getSeconds() + 5);
client.waitForReady(deadline, (err) => {
	if (err) {
		console.log("error: ", err);
	}

	console.log("client is ready");
});

const app: Express = express();

type requestParams = {
	title: string;
};
app.get("/", (req: Request<requestParams>, res: Response) => {
	const { title } = req.query;

	if (!title) {
		res.status(400).json({ error: "invalid title" });
	}

	const reqTitle = title as string;

	const message = GetMovieByTitleRequest.create({
		title: reqTitle,
	});

	client.getMovieByTitle(message, (err, value) => {
		if (err) {
			if (err.code === 3) {
				res.status(400).json({ error: err.message });
			}

			res.status(500).json({ error: err.message });
		}

		if (value) {
			res.status(200).json(value.movie);
		}
	});
});

app.get("/test", (req: Request<requestParams>, res: Response) => {
	const { title } = req.query;

	fetch(`http://localhost:3000/?title=${title}`)
		.then((response) => response.json())
		.then((data) => res.json(data));
});

app.listen(8080, () => {
	console.log("Server running on port 8080");
});
