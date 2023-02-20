import { LogFilter } from "@/lib/api/telemetry/pb";
import antlr4 from "antlr4";
import type Recognizer from "antlr4/Recognizer";
import FiltersLexer from "./FiltersLexer";
import FiltersParser from "./FiltersParser";

export interface ParseResult {
  filter?: LogFilter;
  errorMessage?: string;
}

export function parseFilter(text: string): ParseResult {
  // Empty filter is also valid
  if (text === "") {
    return {};
  }

  // Error listeners here are removed, since by default they output to console.
  // Instead an error listener only for syntax errors is added to check if everything was parsed properly.
  const lexer = new FiltersLexer(new antlr4.InputStream(text));
  lexer.removeErrorListeners();

  const parser = new FiltersParser(new antlr4.CommonTokenStream(lexer));
  parser.buildParseTrees = true;

  const errlistener = new errorListener();
  parser.removeErrorListeners();
  parser.addErrorListener(errlistener);

  const ctx = parser.filter(undefined);
  if (errlistener.errorMessage) {
    return { errorMessage: errlistener.errorMessage };
  }

  return {
    filter: traverse(ctx.getChild<FiltersParser.ExpressionContext>(0)!),
  };
}

class errorListener extends antlr4.error.ErrorListener {
  errorMessage?: string;

  syntaxError(
    recognizer: Recognizer,
    offendingSymbol: antlr4.Token,
    line: number,
    column: number,
    msg: string,
    e: antlr4.error.RecognitionException
  ): void {
    this.errorMessage = `Error near character #${column}: ${msg}`;
  }
}

function traverse(ctx: FiltersParser.ExpressionContext): LogFilter {
  switch (ctx.filterKind) {
    case "text":
      return {
        filter: {
          oneofKind: "text",
          text: { value: JSON.parse(ctx.getText()) },
        },
      };
    case "scoped":
      return {
        filter: {
          oneofKind: "scoped",
          scoped: {
            field: ctx.getChild<antlr4.tree.TerminalNode>(0)!.getText(),
            value: JSON.parse(
              ctx.getChild<antlr4.tree.TerminalNode>(2)!.getText()
            ),
          },
        },
      };
    case "and":
      return {
        filter: {
          oneofKind: "and",
          and: {
            a: traverse(ctx.getChild<FiltersParser.ExpressionContext>(0)!),
            b: traverse(ctx.getChild<FiltersParser.ExpressionContext>(2)!),
          },
        },
      };
    case "or":
      return {
        filter: {
          oneofKind: "or",
          or: {
            a: traverse(ctx.getChild<FiltersParser.ExpressionContext>(0)!),
            b: traverse(ctx.getChild<FiltersParser.ExpressionContext>(2)!),
          },
        },
      };
    case "not":
      return {
        filter: {
          oneofKind: "not",
          not: {
            a: traverse(ctx.getChild<FiltersParser.ExpressionContext>(1)!),
          },
        },
      };
  }
}
