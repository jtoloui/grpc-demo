import winston from "winston";
import expressWinston from "express-winston";

export const logger = (label: string) =>
	winston.createLogger({
		format: winston.format.combine(
			winston.format.ms(),
			winston.format.label({ label }),
			winston.format.timestamp(),
			winston.format((info) => {
				const headers = info.meta?.req?.headers;
				let tracer = null;
				if (headers) {
					tracer = headers["x-tracer-id"];
				}
				info?.meta ? delete info?.meta.req : null;
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
	winstonInstance: logger("express"),
	level: function (req, res) {
		var level = "";
		if (res.statusCode >= 100) {
			level = "info";
		}
		if (res.statusCode >= 400) {
			level = "warn";
		}
		if (res.statusCode >= 500) {
			level = "error";
		}
		// Ops is worried about hacking attempts so make Unauthorized and Forbidden critical
		if (res.statusCode === 401 || res.statusCode === 403) {
			level = "critical";
		}
		// No one should be using the old path, so always warn for those
		if (req.path === "/v1" && level === "info") {
			level = "warn";
		}
		return level;
	},
	meta: true,
});
