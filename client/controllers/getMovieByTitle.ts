import {
	GetMovieByTitleRequest,
	MoviesServiceClient,
} from "@jtoloui/proto-store";
import { Request, Response } from "express";
import { Metadata } from "@grpc/grpc-js";
import { ParamsDictionary } from "express-serve-static-core";
import { logger as log } from "../middleware";

interface paramsDictionary extends ParamsDictionary {
	title: string;
}

type requestParams = {
	title: string;
};

const logger = log("getMovieByTitle");

export const getMovieByTitle = (client: MoviesServiceClient) => {
	return (
		req: Request<paramsDictionary, {}, {}, requestParams>,
		res: Response
	) => {
		const { title } = req.query;
		let tracerId = req.get("x-tracer-id");

		if (!tracerId) tracerId = "no-tracer-id";
		const metadata = new Metadata();
		metadata.add("X-Tracer-Id", tracerId as string);

		if (!title) {
			res.status(400).json({ error: "invalid title" });
		}

		const reqTitle = title as string;

		const message = GetMovieByTitleRequest.create({
			title: reqTitle,
		});

		client.getMovieByTitle(message, metadata, (err, value) => {
			if (err) {
				if (err.code === 3) {
					logger.error(err.message);
					res.status(400).json({ error: err.message });
				}
				logger.error(err.message);
				res.status(500).json({ error: err.message });
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
