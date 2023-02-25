import {
	GetMovieByTitleRequest,
	MoviesServiceClient,
} from "@jtoloui/proto-store";
import { Request, Response } from "express";

type requestParams = {
	title: string;
};

export const getMovieByTitle = (client: MoviesServiceClient) => {
	return (req: Request<requestParams>, res: Response) => {
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
	};
};
