import winston from "winston";
import expressWinston from "express-winston";

const logger = winston.createLogger({
	format: winston.format.combine(
		winston.format.ms(),
		winston.format.timestamp(),
		winston.format((info) => {
			const headers = info.meta?.req?.headers;
			const tracer = headers["x-tracer-id"];
			info?.meta ? delete info?.meta.req : null;
			console.log(info);

			return {
				...info,
				...(tracer && { "x-tracer-id": tracer }),
			};
		})(),
		winston.format.json()
	),
	transports: [new winston.transports.Console()],
});

export const winstonLogger = expressWinston.logger({
	winstonInstance: logger,
	meta: true,
});
