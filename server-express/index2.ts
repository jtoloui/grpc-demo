import express, { Express, Request, Response } from "express";
import dotenv from "dotenv";
const app = express();

dotenv.config();

app.get("/", (req: Request, res: Response) => {
	if (!req.query.title) {
		res.status(400).json({ error: "invalid title" });
	}

	res.status(200).json({ title: req.query.title, year: 1972, director: "" });
});

app.listen(3000, () => {
	console.log("server is listening on port 3000");
});
