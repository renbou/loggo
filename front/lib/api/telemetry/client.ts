import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport";
import { TelemetryClient } from "./pb";

export interface LogBatch {
  messages: Array<string>;
  nextPageToken: string;
}

export class Client {
  telemetryClient: TelemetryClient;

  constructor() {
    const transport = new GrpcWebFetchTransport({
      baseUrl: "http://localhost:3000",
    });

    this.telemetryClient = new TelemetryClient(transport);
  }

  listLogMessages(
    from: Date,
    to: Date,
    filter: string,
    pageToken: string
  ): LogBatch {
    return { messages: new Array(), nextPageToken: "" };
  }
}
