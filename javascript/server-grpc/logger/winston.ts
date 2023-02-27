import winston, { createLogger, transports, format } from "winston";

const myFormat = format.printf((info) => {
	return `${info.timestamp} [${info.label}] ${info.level}: ${info.message} ${info.ms}`;
});
export const logger = (label: string) =>
	createLogger({
		transports: [
			new transports.Console({
				level: "info",
				format: format.combine(
					format.label({ label }),
					format.colorize(),
					format.timestamp({ format: "DD-MM-YYYY HH:mm:ss" }),
					format.ms(),
					format.splat(),
					format.simple()
				),
			}),
		],
	});
