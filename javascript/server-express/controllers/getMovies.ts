import { Metadata } from "@grpc/grpc-js";
import { MoviesServiceClient } from "@jtoloui/proto-store";
import { Request, Response } from "express";
import { ParamsDictionary } from "express-serve-static-core";

import { logger as log } from "../middleware";

const logger = log("getMovies");

interface paramsDictionary extends ParamsDictionary {
	page: string;
	per_page: string;
}

type queryParams = {
	page: string;
	per_page: string;
};

export const getMovies =
	(client: MoviesServiceClient) =>
	(req: Request<paramsDictionary, {}, {}, queryParams>, res: Response) => {
		let tracerId = req.get("x-tracer-id");

		const { page: reqPage, per_page: reqPerPage } = req.query;

		let page = 1;
		let perPage = 10;

		if (reqPage) {
			page = parseInt(reqPage as string, 10);
		}

		if (reqPerPage) {
			perPage = parseInt(reqPerPage as string, 10);
		}

		if (!tracerId) tracerId = "no-tracer-id";
		const metadata = new Metadata();
		metadata.add("X-Tracer-Id", tracerId as string);

		client.getMovies(
			{
				page,
				perPage,
			},
			metadata,
			(err, value) => {
				if (err) {
					logger.error(err.message);
					return res.status(500).json({ error: err.message });
				}

				if (value) {
					logger.info(`tracer response: ${tracerId}`, {
						movies: value.movies,
					});
					return res.status(200).json({
						movies: value.movies,
						total: value.total,
					});
				}
			}
		);
	};
