import { MoviesServiceClient } from "@jtoloui/proto-store";
import { Request, Response } from "express";
import { Movie } from "../models";

type requestBody = {
	title: string;
	year: number;
	director: string;
};

export const createMovie = async (
	req: Request<{}, {}, requestBody>,
	res: Response
) => {
	const movie = Movie.build({
		title: req.body.title,
		year: req.body.year,
		director: req.body.director,
	});

	await movie.save();

	return res.status(200).json({
		id: movie.id,
		title: movie.title,
		year: movie.year,
		director: movie.director,
	});

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
