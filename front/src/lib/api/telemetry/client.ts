import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport";
import { LogFilter, TelemetryClient } from "./pb";
import { Timestamp } from "./pb/google/protobuf/timestamp";

export interface LogMessage {
  message: string;
  id: string;
}

export interface LogBatch {
  messages: Array<LogMessage>;
  nextPageToken: string;
}

const byteToHex: string[] = [];

for (let n = 0; n <= 0xff; ++n) {
  const hexOctet = n.toString(16).padStart(2, "0");
  byteToHex.push(hexOctet);
}

// https://stackoverflow.com/questions/40031688/javascript-arraybuffer-to-hex
function hex(buffer: Uint8Array) {
  const hexOctets = [];
  for (let i = 0; i < buffer.length; ++i) {
    hexOctets.push(byteToHex[buffer[i]]);
  }
  return hexOctets.join("");
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
    const messages = batch.messages.map(
      (message) =>
        ({
          message: td.decode(message.message),
          id: hex(message.id),
        } as LogMessage)
    );

    return { messages, nextPageToken: batch.nextPageToken };
  }
}
