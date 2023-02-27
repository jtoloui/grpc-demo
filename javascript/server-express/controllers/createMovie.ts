import { Metadata } from "@grpc/grpc-js";
import { CreateMovieRequest, IMoviesServiceClient } from "@jtoloui/proto-store";
import { Request, Response } from "express";

import { logger as log } from "../middleware";

type requestBody = {
	title: string;
	year: number;
	director: string;
};

export const createMovie =
	(client: IMoviesServiceClient) =>
	async (req: Request<{}, {}, requestBody>, res: Response) => {
		const { title, year, director } = req.body;
		let tracerId = req.get("x-tracer-id");

		const logger = log("createMovie");
		if (!tracerId) tracerId = "no-tracer-id";
		const metadata = new Metadata();
		metadata.add("X-Tracer-Id", tracerId as string);

		if (!title && !year && !director) {
			logger.error("invalid title");
			return res.status(400).json({ error: "invalid title" });
		}

		const movieRequest = CreateMovieRequest.create({
			movie: {
				title,
				year,
				director,
			},
		});

		client.createMovie(movieRequest, metadata, (err, response) => {
			if (err) {
				logger.error(err.message);
				return res.status(500).json({ error: err.message });
			}

			logger.info(`tracer response: ${tracerId}`);
			return res.status(200).json({
				id: response?.id,
				title: response?.movie?.title,
				year: response?.movie?.year,
				director: response?.movie?.director,
			});
		});
	};
