import mongoose from "mongoose";
import { logger as log } from "../logger";

export const connectDB = async () => {
	const uri = `mongodb+srv://${process.env.MONGO_USER}:${process.env.MONGO_PW}@cluster0-cs2gr.mongodb.net/${process.env.MONGO_DB}?retryWrites=true&w=majority`;
	const logger = log("db");
	try {
		mongoose.set("strictQuery", false);
		await mongoose.connect(uri);
		logger.info("Connected to database");
		return mongoose.connection;
	} catch (error) {
		logger.error(`Error: ${error}`);
		process.exit(1);
	}
};
