import express, { Express } from "express";
import dotenv from "dotenv";

import { createMovie, getMovieByTitle } from "./controllers";
import { logger as log, setTracerdHeader, winstonLogger } from "./middleware";
import { client as grpcClient } from "./grpcClient";
import { json } from "body-parser";

dotenv.config();

const logger = log("server-express");

const client = grpcClient;

const app: Express = express();

app.use(setTracerdHeader);
app.use(winstonLogger);
app.use(express.json());

app.get("/", getMovieByTitle(client));
app.post("/", createMovie);

// app.get("/test", (req: Request<requestParams>, res: Response) => {
// 	const { title } = req.query;

// 	fetch(`http://localhost:3000/?title=${title}`)
// 		.then((response) => response.json())
// 		.then((data) => res.json(data));
// });

app.listen(8080, () => {
	logger.info("Server running on port 8080");
});
