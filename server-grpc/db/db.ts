import mongoose from "mongoose";

export const connectDB = async () => {
	const uri = `mongodb+srv://${process.env.MONGO_USER}:${process.env.MONGO_PW}@cluster0-cs2gr.mongodb.net/${process.env.MONGO_DB}?retryWrites=true&w=majority`;

	try {
		mongoose.set("strictQuery", false);
		await mongoose.connect(uri);
		return mongoose.connection;
	} catch (error) {
		console.error(`Error: ${error}`);
		process.exit(1);
	}
};
