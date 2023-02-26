import { MoviesServiceClient } from "@jtoloui/proto-store";
import { Request, Response } from "express";

type requestBody = {
	title: string;
	year: number;
	director: string;
};

export const createMovie = async (
	req: Request<{}, {}, requestBody>,
	res: Response
) => {
	const { title, year, director } = req.body;
	res.json({ title, year, director });

	// const movie = Movie.create({
	// 	title: req.body.title,
	// 	year: req.body.year,
	// 	director: req.body.director,
	// });
	// movie
	// 	.then((value) => {
	// 		res.status(200).json({
	// 			id: value.id,
	// 			title: value.title,
	// 			year: value.year,
	// 			director: value.director,
	// 		});
	// 	})
	// 	.catch((err) => {
	// res.status(500).json({ error: "err.message" });
	// 	});
};
