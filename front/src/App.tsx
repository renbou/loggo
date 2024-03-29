import { useMemo, useState } from "react";
import LogList from "./components/LogList";
import Navbar from "./components/Navbar";
import * as api from "./lib/api/telemetry";
import { parseFilter } from "./lib/filters";

type searchError = {
  title: string;
  info: string;
};

function App() {
  // apiClient is created only once
  const apiClient = useMemo(() => new api.Client(), []);

  const [searchError, setSearchError] = useState<searchError>();
  const [logMessages, setLogMessages] = useState<api.LogMessage[]>();

  // TODO: add support for streaming logs
  async function runSearch(search: string, from: Date, to: Date) {
    const result = parseFilter(search);
    if (result.errorMessage) {
      setSearchError({
        title: "Failed to parse search query",
        info: result.errorMessage,
      });
      return;
    }

    const logBatch = await apiClient.listLogMessages(
      from,
      to,
      result.filter,
      ""
    );

    setSearchError(undefined);
    setLogMessages(logBatch.messages);
  }

  return (
    <div className="flex flex-col h-screen">
      <Navbar onSearch={runSearch}></Navbar>
      <LogList messages={logMessages || []} error={searchError}></LogList>
    </div>
  );
}

export default App;
