import express, { Express } from "express";
import dotenv from "dotenv";

import { createMovie, getMovieById, getMovies } from "./controllers";
import { logger as log, setTracerdHeader, winstonLogger } from "./middleware";
import { client as grpcClient } from "./grpcClient";

dotenv.config();

const logger = log("server-express");

const client = grpcClient;

const app: Express = express();

app.use(setTracerdHeader);
app.use(winstonLogger);
app.use(express.json());

app.get("/", getMovies(client));
app.post("/", createMovie(client));
app.get("/:id", getMovieById(client));

// app.get("/test", (req: Request<requestParams>, res: Response) => {
// 	const { title } = req.query;

// 	fetch(`http://localhost:3000/?title=${title}`)
// 		.then((response) => response.json())
// 		.then((data) => res.json(data));
// });

app.listen(8080, () => {
	logger.info("Server running on port 8080");
});
