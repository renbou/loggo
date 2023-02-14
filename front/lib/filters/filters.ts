import { LogFilter } from "@/lib/api/telemetry/pb";
import antlr4 from "antlr4";
import FiltersLexer from "./FiltersLexer";
import FiltersParser from "./FiltersParser";

export interface ParseResult {
  filter: LogFilter;
  error: string;
}

export function parseFilter(text: string): ParseResult {
  const chars = new antlr4.InputStream(text);
  const lexer = new FiltersLexer(chars);
  const tokens = new antlr4.CommonTokenStream(lexer);
  const parser = new FiltersParser(tokens);
  parser.buildParseTrees = true;

  const ctx = parser.filter();
  return { filter: { filter: { oneofKind: undefined } }, error: "" };
}
