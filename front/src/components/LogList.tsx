import { Card } from "flowbite-react";
import * as api from "../lib/api/telemetry";

type Props = {
  messages: api.LogMessage[];
  error?: { title: string; info: string };
};

function ErrorCard(props: Props["error"]) {
  return (
    <Card>
      <h3 className="text-xl font-bold text-red-600">{props?.title}</h3>
      <p className="font-normal text-gray-700">{props?.info}</p>
    </Card>
  );
}

function MessageCard(props: { message: string }) {
  let formatted = props.message;
  try {
    const json = JSON.parse(props.message);
    formatted = JSON.stringify(json, null, 2);
  } catch (_) {}

  return (
    <Card>
      <code className="whitespace-pre">{formatted}</code>
    </Card>
  );
}

function LogList(props: Props) {
  return (
    <div className="bg-slate-100 flex-grow">
      <div className="mx-auto max-w-6xl h-full rounded-t-2xl bg-gray-50 border border-gray-300 p-4 flex flex-col gap-2">
        {props.error ? (
          <ErrorCard
            title={props.error.title}
            info={props.error.info}
          ></ErrorCard>
        ) : (
          props.messages.map((m) => {
            // TODO: use message ID from the database once it's available in the API
            return <MessageCard message={m.message} key={m.id}></MessageCard>;
          })
        )}
      </div>
    </div>
  );
}

export default LogList;
