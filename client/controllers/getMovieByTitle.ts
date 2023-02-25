import {
	GetMovieByTitleRequest,
	MoviesServiceClient,
} from "@jtoloui/proto-store";
import { Request, Response } from "express";
import { Metadata } from "@grpc/grpc-js";
import { IncomingHttpHeaders } from "http";
import { ParamsDictionary } from "express-serve-static-core";

interface paramsDictionary extends ParamsDictionary {
	title: string;
}

type requestParams = {
	title: string;
};

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
					res.status(400).json({ error: err.message });
				}

				res.status(500).json({ error: err.message });
			}

			if (value) {
				res.status(200).json(value.movie);
			}
		});
	};
};
