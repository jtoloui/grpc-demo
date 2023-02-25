import winston, { createLogger, transports, format } from "winston";

export const logger = (label: string) =>
	createLogger({
		transports: [
			new transports.Console({
				level: "info",
				format: format.combine(
					format.label({ label }),
					format.timestamp(),
					format.ms(),
					format.json()
				),
			}),
		],
	});
