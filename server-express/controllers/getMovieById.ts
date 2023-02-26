import { GetMovieByIdRequest, MoviesServiceClient } from "@jtoloui/proto-store";
import { Request, Response } from "express";
import { Metadata } from "@grpc/grpc-js";
import { ParamsDictionary } from "express-serve-static-core";
import { logger as log } from "../middleware";

interface paramsDictionary extends ParamsDictionary {
	id: string;
}

type requestParams = {
	id: string;
};

const logger = log("getMovieByTitle");

export const getMovieById = (client: MoviesServiceClient) => {
	return (
		req: Request<paramsDictionary, {}, {}, requestParams>,
		res: Response
	) => {
		const { id } = req.query;
		let tracerId = req.get("x-tracer-id");

		if (!tracerId) tracerId = "no-tracer-id";
		const metadata = new Metadata();
		metadata.add("X-Tracer-Id", tracerId as string);

		if (!id) {
			res.status(400).json({ error: "invalid id" });
		}

		const reqTitle = id as string;

		const message = GetMovieByIdRequest.create({
			id: reqTitle,
		});

		client.getMovieById(message, metadata, (err, value) => {
			if (err) {
				switch (err.code) {
					case 3:
						logger.error(err.message);
						res.status(400).json({ error: err.message });
						break;
					case 5:
						logger.error(err.message);
						res.status(404).json({ error: err.details });
						break;
					default:
						logger.error(err.message);
						res.status(500).json({ error: err.message });
						break;
				}
			}

			if (value) {
				logger.info(`tracer response: ${tracerId}`, {
					movie: value.movie,
				});
				res.status(200).json(value.movie);
			}
		});
	};
};
