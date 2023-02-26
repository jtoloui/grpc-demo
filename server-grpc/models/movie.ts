import { Document, Model, model, Schema } from "mongoose";

interface IMovie {
	title: string;
	year: number;
	director: string;
}

interface IMovieDoc extends Document {
	title: string;
	year: number;
	director: string;
}

interface movieModelInterface extends Model<IMovieDoc> {
	build(attr: IMovie): IMovieDoc;
}

const movieSchema = new Schema({
	title: {
		type: String,
		required: true,
	},
	year: {
		type: Number,
		required: true,
	},
	director: {
		type: String,
		required: true,
	},
});

const Movie = model<IMovieDoc, movieModelInterface>("Movie", movieSchema);

export { Movie };
