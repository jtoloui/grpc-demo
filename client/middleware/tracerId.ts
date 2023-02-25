import { NextFunction, Request, Response } from "express";
import { v4 as uuid } from "uuid";

export const setTracerdHeader = (
	req: Request,
	res: Response,
	next: NextFunction
) => {
	const tracingId = uuid();
	req.headers["x-tracer-id"] = tracingId;
	res.setHeader("X-Tracer-Id", tracingId);
	next();
};
