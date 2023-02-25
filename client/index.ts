import express, { Express } from "express";
import dotenv from "dotenv";

import { getMovieByTitle } from "./controllers";
import { winstonLogger } from "./middleware";
import { client as grpcClient } from "./grpcClient";

dotenv.config();

const client = grpcClient;

const app: Express = express();

app.use(winstonLogger);

app.get("/", getMovieByTitle(client));

// app.get("/test", (req: Request<requestParams>, res: Response) => {
// 	const { title } = req.query;

// 	fetch(`http://localhost:3000/?title=${title}`)
// 		.then((response) => response.json())
// 		.then((data) => res.json(data));
// });

app.listen(8080, () => {
	console.log("Server running on port 8080");
});
