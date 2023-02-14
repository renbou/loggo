import { useState } from "react";
import "./App.css";
import { parseFilter } from "@/lib/filters/filters";

function App() {
  const [text, setText] = useState("");
  const [parsed, setParsed] = useState({});

  function changed(e: React.ChangeEvent<HTMLInputElement>) {
    setText(e.target.value);
    const parsed = parseFilter(e.target.value);
    setParsed(parsed);
  }

  return (
    <div className="App">
      <div>
        <a href="https://vitejs.dev" target="_blank"></a>
        <a href="https://reactjs.org" target="_blank"></a>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        <input type="text" value={text} onChange={changed} />
        <p>
          <code>{JSON.stringify(parsed)}</code>
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </div>
  );
}

export default App;
