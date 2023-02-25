import winston from "winston";
import expressWinston from "express-winston";

export const winstonLogger = expressWinston.logger({
	transports: [new winston.transports.Console()],
	statusLevels: true,
	format: winston.format.combine(
		// winston.format.colorize(),
		winston.format.json()
	),
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
		if (res.statusCode === -401 || res.statusCode === 403) {
			level = "critical";
		}
		// No one should be using the old path, so always warn for those
		if (req.path === "/v1" && level === "info") {
			level = "warn";
		}
		return level;
	},
	meta: false,
	msg: "HTTP {{req.method}} {{req.url}}",
	expressFormat: true,
	colorize: false,
	ignoreRoute: function (req, res) {
		return false;
	},
});
