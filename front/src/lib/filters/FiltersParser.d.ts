import type antlr4 from "antlr4";

type filterKind = "text" | "scoped" | "and" | "or" | "not";

declare class FiltersParser extends antlr4.Parser {
  constructor(input: any);
  filter(_p: any): FilterContext;
}

declare namespace FiltersParser {
  export { FilterContext, ExpressionContext };
}

export default FiltersParser;

declare class FilterContext extends antlr4.ParserRuleContext {}

declare class ExpressionContext extends antlr4.ParserRuleContext {
  filterKind: filterKind;
}
