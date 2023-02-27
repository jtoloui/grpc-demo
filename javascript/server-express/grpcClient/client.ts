import { ChannelCredentials } from "@grpc/grpc-js";
import { MoviesServiceClient } from "@jtoloui/proto-store";
import { logger as log } from "../middleware";

const logger = log("grpcClient");

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
		logger.error("error: ", err);
	}

	logger.info("MoviesServiceClient is ready");
});
