import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport";
import { LogFilter, TelemetryClient } from "./pb";
import { Timestamp } from "./pb/google/protobuf/timestamp";

export interface LogBatch {
  messages: Array<string>;
  nextPageToken: string;
}

// TODO: add authorization headers here
export class Client {
  private telemetryClient: TelemetryClient;

  constructor() {
    const transport = new GrpcWebFetchTransport({
      baseUrl: location.origin,
      format: "text",
    });

    this.telemetryClient = new TelemetryClient(transport);
  }

  async listLogMessages(
    from: Date,
    to: Date,
    filter: LogFilter | undefined,
    pageToken: string
  ): Promise<LogBatch> {
    const call = this.telemetryClient.listLogMessages({
      from: Timestamp.fromDate(from),
      to: Timestamp.fromDate(to),
      filter: filter,
      pageSize: 0,
      pageToken: pageToken,
    });

    const batch = (await call.response).batch!;

    // Avoid failing on decode since formally we allow any log messages to be sent to the agregator
    const td = new TextDecoder("utf-8", { fatal: false });
    const messages = batch.messages.map((message) => td.decode(message));

    return { messages, nextPageToken: batch.nextPageToken };
  }
}
