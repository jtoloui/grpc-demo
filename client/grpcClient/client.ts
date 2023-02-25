import { ChannelCredentials } from "@grpc/grpc-js";
import { MoviesServiceClient } from "@jtoloui/proto-store";

export const client = new MoviesServiceClient(
	"0.0.0.0:50051",
	ChannelCredentials.createInsecure(),
	{},
	{}
);

const deadline = new Date();
deadline.setSeconds(deadline.getSeconds() + 5);
client.waitForReady(deadline, (err) => {
	if (err) {
		console.log("error: ", err);
	}

	console.log("client is ready");
});
